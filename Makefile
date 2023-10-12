GO_ARG_LINUX=GOOS=linux GOARCH=amd64

SRC_PATH?=$(shell pwd)
BIN_PATH?=$(SRC_PATH)/bin

BOOTSTRAP_APP_PREINSTALL_NAME=bootstrapapp-pre-install
BOOTSTRAP_APP_POST_NAME=bootstrapapp-post
BOOTSTRAP_APP_PREDELETE_NAME=bootstrapapp-pre-delete

BOOTSTRAP_CMD_PREINSTALL_PATH=./cmd/bootstrap-pre-install
BOOTSTRAP_CMD_POST_PATH=./cmd/bootstrap-post
BOOTSTRAP_CMD_PREDELETE_PATH=./cmd/bootstrap-pre-delete

build-bootstrap-pre-install:
	@echo "\n Building bootstrap-pre-install binary"
	env $(GO_ARG_LINUX) CGO_ENABLED=0 go build -o "$(BIN_PATH)/$(BOOTSTRAP_APP_PREINSTALL_NAME)" -ldflags \
		"-X main.version=$(MAIN_VERSION)" "$(BOOTSTRAP_CMD_PREINSTALL_PATH)"

build-bootstrap-post:
	@echo "\n Building bootstrapapp-post binary"
	env $(GO_ARG_LINUX) CGO_ENABLED=0 go build -o "$(BIN_PATH)/$(BOOTSTRAP_APP_POST_NAME)" -ldflags \
		"-X main.version=$(MAIN_VERSION)" "$(BOOTSTRAP_CMD_POST_PATH)"

build-bootstrap-pre-delete:
	@echo "\n Building bootstrapapp-pre-delete binary"
	env $(GO_ARG_LINUX) CGO_ENABLED=0 go build -o "$(BIN_PATH)/$(BOOTSTRAP_APP_PREDELETE_NAME)" -ldflags \
		"-X main.version=$(MAIN_VERSION)" "$(BOOTSTRAP_CMD_PREDELETE_PATH)"

local-pre-install:
	go build -o "$(BIN_PATH)/$(BOOTSTRAP_APP_PREINSTALL_NAME)" "$(BOOTSTRAP_CMD_PREINSTALL_PATH)"
	"$(BIN_PATH)/$(BOOTSTRAP_APP_PREINSTALL_NAME)"

local-post:
	go build -o "$(BIN_PATH)/$(BOOTSTRAP_APP_POST_NAME)" "$(BOOTSTRAP_CMD_POST_PATH)"
	"$(BIN_PATH)/$(BOOTSTRAP_APP_POST_NAME)"

local-pre-delete:
	go build -o "$(BIN_PATH)/$(BOOTSTRAP_APP_PREDELETE_NAME)" "$(BOOTSTRAP_CMD_PREDELETE_PATH)"
	"$(BIN_PATH)/$(BOOTSTRAP_APP_PREDELETE_NAME)"

build-all: build-bootstrap-post build-bootstrap-pre-delete build-bootstrap-pre-install
