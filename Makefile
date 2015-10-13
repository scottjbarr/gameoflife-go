# See http://peter.bourgon.org/go-in-production/
GO ?= go

BUILD_DIR = build

APP = gameoflife-go
APP_BUILD = $(BUILD_DIR)/$(APP)

BUILD = $(GO) build

GO_FILES = `ls *.go | grep -v test`

all: clean build

build:
	$(BUILD) -o $(APP_BUILD)

run:
	$(GO) run $(GO_FILES) -file ./docs/glider.L -iterations 31

test:
	$(GO) test

clean:
	rm -rf $(BUILD_DIR)
