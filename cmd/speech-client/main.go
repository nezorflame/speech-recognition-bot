package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/nezorflame/speech-recognition-bot/internal/app"
	"github.com/nezorflame/speech-recognition-bot/pkg/yandex"
)

var (
	ycToken       string
	folderID      string
	audioFilePath string
	debug         bool
)

func init() {
	flag.StringVar(&ycToken, "token", "", "Yandex Cloud OAuth token")
	flag.StringVar(&folderID, "folder-id", "", "Yandex Cloud folder ID")
	flag.StringVar(&audioFilePath, "audio-file", "", "Audio file path (for recognition)")
	flag.BoolVar(&debug, "debug", false, "Show application and debug info")
	flag.Parse()

	if ycToken == "" || folderID == "" || audioFilePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	log.Println("Launching speech recognition client")
	app.PrintInfo(debug)
	defer log.Println("Client finished")

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

	result, err := rClient.SimpleRecognize(audioFilePath)
	if err != nil {
		log.Fatalf("Unable to recognize audio file: %v", err)
	}

	log.Println("Recognition result:", result)
}
