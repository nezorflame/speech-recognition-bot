package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nezorflame/speech-recognition-bot/internal/file"
	"github.com/nezorflame/speech-recognition-bot/pkg/yandex"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Bot describes Telegram bot
type Bot struct {
	ctx  context.Context
	api  *tgbotapi.BotAPI
	cfg  *viper.Viper
	ySDK *yandex.SDK
}

// NewBot creates new instance of Bot
func NewBot(ctx context.Context, cfg *viper.Viper, sdk *yandex.SDK) (*Bot, error) {
	if cfg == nil {
		return nil, errors.New("empty config")
	}

	api, err := tgbotapi.NewBotAPI(cfg.GetString("telegram.token"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to Telegram")
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Enabling debug mode for bot")
		api.Debug = true
	}

	log.Debugf("Authorized on account %s", api.Self.UserName)
	return &Bot{api: api, cfg: cfg, ctx: ctx, ySDK: sdk}, nil
}

// Start starts to listen the bot updates channel
func (b *Bot) Start() {
	update := tgbotapi.NewUpdate(0)
	update.Timeout = b.cfg.GetInt("telegram.timeout")
	b.listen(b.api.GetUpdatesChan(update))
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.ySDK.Close()
	b.api.StopReceivingUpdates()
}

func (b *Bot) listen(updates tgbotapi.UpdatesChannel) {
	for u := range updates {
		if u.Message == nil { // ignore any non-Message Updates
			continue
		}

		if !isIDInWhitelist(u.Message.Chat.ID, b.cfg.GetStringSlice("telegram.whitelist")) {
			go b.reject(u.Message)
			continue
		}

		switch {
		case strings.HasPrefix(u.Message.Text, b.cfg.GetString("commands.start")):
			go b.hello(u.Message)
		case u.Message.Voice != nil:
			go b.parseVoice(u.Message)
		}
	}
}

func (b *Bot) hello(msg *tgbotapi.Message) {
	b.reply(msg.Chat.ID, msg.MessageID, b.cfg.GetString("messages.hello"))
}

func (b *Bot) reject(msg *tgbotapi.Message) {
	b.reply(msg.Chat.ID, msg.MessageID, b.cfg.GetString("errors.whitelist"))
}

func (b *Bot) parseVoice(msg *tgbotapi.Message) {
	replyID := b.reply(msg.Chat.ID, msg.MessageID, b.cfg.GetString("messages.in_progress"))
	if replyID == 0 {
		b.reply(msg.Chat.ID, msg.MessageID, b.cfg.GetString("errors.unknown"))
		return
	}

	audioFile, err := b.api.GetFile(tgbotapi.FileConfig{FileID: msg.Voice.FileID})
	if err != nil {
		log.WithError(err).Errorf("Unable to get audio file")
		b.reply(msg.Chat.ID, msg.MessageID, b.cfg.GetString("errors.download"))
	}

	fileLink := audioFile.Link(b.cfg.GetString("telegram.token"))
	filePath, err := file.Download(fileLink)
	if err != nil {
		log.WithError(err).Errorf("Unable to download audio file")
		b.reply(msg.Chat.ID, msg.MessageID, b.cfg.GetString("errors.download"))
	}

	result, err := b.getRecognition(filePath)
	if err != nil {
		log.WithError(err).Errorf("Unable to download audio file")
		b.reply(msg.Chat.ID, msg.MessageID, b.cfg.GetString("errors.download"))
	}

	if err = b.editMsg(msg.Chat.ID, replyID, result); err != nil {
		log.WithError(err).WithField("reply_id", replyID).Errorf("Unable to edit the message")
	}
}

func (b *Bot) getRecognition(filePath string) (string, error) {
	rClient, err := b.ySDK.NewRecognitionClient(b.ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to create STT client")
	}
	defer rClient.Close()
	return rClient.SimpleRecognize(filePath, b.cfg.GetString("yandex.lang"))
}

func (b *Bot) reply(chatID int64, msgID int, text string) int {
	logger := log.WithField("chat_id", chatID).WithField("msg_id", msgID)
	logger.Debug("Sending reply to the message")
	msg := tgbotapi.NewMessage(chatID, fmt.Sprint(text))
	if msgID != 0 {
		msg.ReplyToMessageID = msgID
	}
	msg.ParseMode = tgbotapi.ModeMarkdown

	replyMsg, err := b.api.Send(msg)
	if err != nil {
		logger.WithError(err).Errorf("Unable to send the message")
		return 0
	}
	return replyMsg.MessageID
}

func (b *Bot) editMsg(chatID int64, msgID int, text string) error {
	msgConfig := tgbotapi.NewEditMessageText(chatID, msgID, text)
	_, err := b.api.Send(msgConfig)
	return err
}

func isIDInWhitelist(chatID int64, whitelist []string) bool {
	cID := strconv.FormatInt(chatID, 10)
	for _, id := range whitelist {
		if id == cID {
			return true
		}
	}
	return false
}
