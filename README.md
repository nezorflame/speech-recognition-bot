# speech-recognition-bot

Speech recognition bot for Telegram using [Yandex SpeechKit API](https://github.com/yandex-cloud/docs/blob/master/en/speechkit/stt/index.md) through gRPC.

Currently only the client for the SpeechKit API is implemented (as POC). Bot will be implemented at the later stages.

## Installation

This project uses Go modules.
To install it, starting with Go 1.12 you can just use `go get`:

`go get github.com/nezorflame/speech-recognition-bot`

or

`go install github.com/nezorflame/speech-recognition-bot/cmd/speech-client`

Also you can just clone this repo and use the build/install targets from `Makefile`.

## Client usage

```text
-audio-file string
    Audio file path (for recognition)
-folder-id string
    Yandex Cloud folder ID
-token string
    Yandex Cloud OAuth token
```
