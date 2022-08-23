GO_ARG_LINUX=GOOS=linux GOARCH=amd64

SRC_PATH?=$(shell pwd)
BIN_PATH?=$(SRC_PATH)/bin
BOOTSTRAP_APP_POST_NAME=bootstrapapp-post
BOOTSTRAP_APP_PREDELETE_NAME=bootstrapapp-pre-delete

BOOTSTRAP_CMD_POST_PATH=./cmd/bootstrap-post
BOOTSTRAP_CMD_PREDELETE_PATH=./cmd/bootstrap-pre-delete


build-bootstrap-post:
	@echo "\n Building bootstrap binary"
	env $(GO_ARG_LINUX) CGO_ENABLED=0 go build -v -o "$(BIN_PATH)/$(BOOTSTRAP_APP_PATH)$(BOOTSTRAP_APP_POST_NAME)" -ldflags \
		"-X main.version=$(MAIN_VERSION)" "$(BOOTSTRAP_CMD_POST_PATH)"

build-bootstrap-pre-delete:
	@echo "\n Building bootstrap binary"
	env $(GO_ARG_LINUX) CGO_ENABLED=0 go build -v -o "$(BIN_PATH)/$(BOOTSTRAP_APP_PATH)$(BOOTSTRAP_APP_PREDELETE_NAME)" -ldflags \
		"-X main.version=$(MAIN_VERSION)" "$(BOOTSTRAP_CMD_PREDELETE_PATH)"

build-all: build-bootstrap-post build-bootstrap-pre-delete