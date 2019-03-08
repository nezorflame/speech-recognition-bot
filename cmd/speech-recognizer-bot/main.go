package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nezorflame/speech-recognition-bot/internal/app"
	"github.com/nezorflame/speech-recognition-bot/internal/config"
	"github.com/nezorflame/speech-recognition-bot/pkg/telegram"
	"github.com/nezorflame/speech-recognition-bot/pkg/yandex"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var configName string

func init() {
	pflag.StringVar(&configName, "config", "config", "Config file name")
	level := pflag.String("log-level", "INFO", "Logrus log level (DEBUG, WARN, etc.)")
	pflag.Parse()

	logLevel, err := log.ParseLevel(*level)
	if err != nil {
		log.Fatalf("Unknown log level: %s", *level)
	}
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)

	if configName == "" {
		pflag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Info("Starting speech recognition bot")
	app.PrintInfo()

	// Init config
	cfg, err := config.New(configName)
	if err != nil {
		log.WithError(err).Fatal("Unable to create config")
	}
	log.Info("Config created")

	// Init Yandex Cloud SDK
	sdk, err := yandex.NewSDK(ctx, cfg.GetString("yandex.token"), cfg.GetString("yandex.folder_id"))
	if err != nil {
		log.WithError(err).Fatal("Unable to init Yandex SDK")
	}
	defer sdk.Close()
	log.Info("Yandex SDK initiated")

	// Create bot
	bot, err := telegram.NewBot(ctx, cfg, sdk)
	if err != nil {
		log.WithError(err).Fatal("Unable to create bot")
	}
	log.Info("Bot created")

	// Init graceful stop chan
	log.Debug("Initiating system signal watcher")
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Warnf("Caught sig %+v, stopping the app in a second", sig)
		bot.Stop()
		cancel()
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	// Start the bot
	log.Info("Starting the bot")
	bot.Start()
}
