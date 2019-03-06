CMD_LIST?=$(shell ls ./cmd)
PKG_LIST:=$(shell go list ./...)
GIT_HASH?=$(shell git log --format="%h" -n 1 2> /dev/null)
GIT_BRANCH?=$(shell git branch 2> /dev/null | grep '*' | cut -f2 -d' ')
GIT_TAG:=$(shell git describe --exact-match --abbrev=0 --tags 2> /dev/null)
APP_VERSION?=$(if $(GIT_TAG),$(GIT_TAG),$(shell git describe --all --long HEAD 2> /dev/null))
GO_VERSION:=$(shell go version)
GO_VERSION_SHORT:=$(shell echo $(GO_VERSION)|sed -E 's/.* go(.*) .*/\1/g')
BUILD_TS:=$(shell date +%FT%T%z)
LDFLAGS:=-X 'github.com/nezorflame/speech-recognition-bot/internal/app.Version=$(APP_VERSION)'\
         -X 'github.com/nezorflame/speech-recognition-bot/internal/app.BuildTS=$(BUILD_TS)'\
         -X 'github.com/nezorflame/speech-recognition-bot/internal/app.GoVersion=$(GO_VERSION_SHORT)'\
         -X 'github.com/nezorflame/speech-recognition-bot/internal/app.GitHash=$(GIT_HASH)'\
         -X 'github.com/nezorflame/speech-recognition-bot/internal/app.GitBranch=$(GIT_BRANCH)'\

# install project dependencies
.PHONY: deps
deps:
	$(info #Install dependencies and clean up...)
	go mod tidy

# run all tests
.PHONY: test
test:
	$(info #Running tests...)
	go test -v -cover -race ./...

# run all tests with coverage
.PHONY: test-cover
test-cover:
	$(info #Running tests with coverage...)
	go test -v -coverprofile=coverage.out -race $(PKG_LIST)
	go tool cover -func=coverage.out | grep total
	rm -f coverage.out
	
.PHONY: fast-build
fast-build: deps
	$(info #Building binaries...)
	$(foreach CMD, $(CMD_LIST), $(shell $(BUILD_ENVPARMS) go build -ldflags "-s -w $(LDFLAGS) -X 'github.com/nezorflame/speech-recognition-bot/internal/app.Name=$(CMD)'" -o ./bin/$(CMD) ./cmd/$(CMD)))
	@echo

.PHONY: build
build: deps fast-build test

.PHONY: install
install:
	$(info #Installing binaries...)
	$(foreach CMD, $(CMD_LIST), $(shell $(BUILD_ENVPARMS) go install -ldflags "-s -w $(LDFLAGS) -X 'github.com/nezorflame/speech-recognition-bot/internal/app.Name=$(CMD)'" ./cmd/$(CMD)))
	@echo
