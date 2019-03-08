package main

import (
	"context"
	"os"

	"github.com/nezorflame/speech-recognition-bot/internal/app"
	"github.com/nezorflame/speech-recognition-bot/pkg/yandex"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	ycToken       string
	folderID      string
	lang          string
	audioFilePath string
)

func init() {
	pflag.StringVar(&ycToken, "token", "", "Yandex Cloud OAuth token")
	pflag.StringVar(&folderID, "folder-id", "", "Yandex Cloud folder ID")
	pflag.StringVar(&lang, "lang", "en-US", "Language to detect")
	pflag.StringVar(&audioFilePath, "audio-file", "", "Audio file path (for recognition)")
	level := pflag.String("log-level", "INFO", "Logrus log level (DEBUG, WARN, etc.)")
	pflag.Parse()

	logLevel, err := log.ParseLevel(*level)
	if err != nil {
		log.Fatalf("Unknown log level: %s", *level)
	}
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)

	if ycToken == "" || folderID == "" || audioFilePath == "" {
		pflag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	log.Info("Launching speech recognition client")
	defer log.Println("Client finished")
	app.PrintInfo()

	ctx := context.Background()
	sdk, err := yandex.NewSDK(ctx, ycToken, folderID)
	if err != nil {
		log.Fatalf("Unable to create: %v", err)
	}
	defer sdk.Close()

	rClient, err := sdk.NewRecognitionClient(ctx)
	if err != nil {
		log.Fatalf("Unable to create STT client: %v", err)
	}
	defer rClient.Close()

	result, err := rClient.SimpleRecognize(audioFilePath, lang)
	if err != nil {
		log.Fatalf("Unable to recognize audio file: %v", err)
	}

	log.Println("Recognition result:", result)
}
