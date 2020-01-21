# speech-recognition-bot [![CircleCI](https://circleci.com/gh/nezorflame/speech-recognition-bot/tree/master.svg?style=svg)](https://circleci.com/gh/nezorflame/speech-recognition-bot/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/nezorflame/speech-recognition-bot)](https://goreportcard.com/report/github.com/nezorflame/speech-recognition-bot) [![GolangCI](https://golangci.com/badges/github.com/nezorflame/speech-recognition-bot.svg)](https://golangci.com/r/github.com/nezorflame/speech-recognition-bot) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnezorflame%2Fspeech-recognition-bot.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnezorflame%2Fspeech-recognition-bot?ref=badge_shield)

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

## Bot

Located at `cmd/speech-recognizer-bot`.
Uses config in the TOML format. Example can be found at `config.example.toml`.

Usage:

```text
--config string      Config file name (default "config")
--log-level string   Logrus log level (DEBUG, INFO, WARN, etc.) (default "INFO")
```

## Client test app

Located at `cmd/speech-client`.

Usage:

```text
--audio-file string   Audio file path (for recognition)
--folder-id string    Yandex Cloud folder ID
--lang string         Language to detect (default "en-US")
--log-level string    Logrus log level (DEBUG, WARN, etc.) (default "INFO")
--token string        Yandex Cloud OAuth token
```

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnezorflame%2Fspeech-recognition-bot.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnezorflame%2Fspeech-recognition-bot?ref=badge_large)
