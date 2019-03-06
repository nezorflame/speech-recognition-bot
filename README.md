# speech-recognition-bot [![CircleCI](https://circleci.com/gh/nezorflame/speech-recognition-bot/tree/master.svg?style=svg)](https://circleci.com/gh/nezorflame/speech-recognition-bot/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/nezorflame/speech-recognition-bot)](https://goreportcard.com/report/github.com/nezorflame/speech-recognition-bot) [![GolangCI](https://golangci.com/badges/github.com/nezorflame/speech-recognition-bot.svg)](https://golangci.com)

Speech recognition bot for Telegram using [Yandex SpeechKit API](https://github.com/yandex-cloud/docs/blob/master/en/speechkit/stt/index.md) through gRPC.

Currently only the client for the SpeechKit API is implemented (as POC). Bot will be implemented at the later stages.

## Installation

This project uses Go modules.
To install it, starting with Go 1.12 you can just use `go get`:

`go get github.com/nezorflame/speech-recognition-bot`

or

`go install github.com/nezorflame/speech-recognition-bot/cmd/speech-client`

Also you can just clone this repo and use the build/install targets from `Makefile`.

## Prerequisits

Make sure you have acquired:

- [OAuth token](https://oauth.yandex.ru/authorize?response_type=token&client_id=1a6990aa636648e9b2ef855fa7bec2fb)
- Folder ID (can be found at your [Cloud](https://console.cloud.yandex.ru/folders/) page after you've selected your project (in the form of `https://console.cloud.yandex.ru/folders/YOUR_FOLDER_ID`)

## Client usage

```text
-audio-file string
    Audio file path (for recognition)
-folder-id string
    Yandex Cloud folder ID
-token string
    Yandex Cloud OAuth token
```
