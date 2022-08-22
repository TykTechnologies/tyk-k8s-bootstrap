GO_ARG_LINUX=GOOS=linux GOARCH=amd64

SRC_PATH?=$(shell pwd)
BIN_PATH?=$(SRC_PATH)/bin
BOOTSTRAP_APP_NAME=bootstrapapp
BOOTSTRAP_CMD_PATH=./cmd/bootstrap

build-bootstrap:
	@echo "\n Building bootstrap binary"
	env $(GO_ARG_LINUX) CGO_ENABLED=0 go build -v -o "$(BIN_PATH)/$(BOOTSTRAP_APP_PATH)$(BOOTSTRAP_APP_NAME)" -ldflags \
		"-X main.version=$(MAIN_VERSION)" "$(BOOTSTRAP_CMD_PATH)"