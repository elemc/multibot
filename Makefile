PROJECT_NAME:=multibot
GO_GET:=go get
PLUGINS:=reminder

all: multibot subdirs

multibot: deps
	@echo "Building ${PROJECT_NAME}"
	@go build

deps:
	${GO_GET} gopkg.in/telegram-bot-api.v4
	${GO_GET} github.com/sirupsen/logrus
	${GO_GET} github.com/spf13/viper
	${GO_GET} github.com/go-pg/pg

subdirs: bin_plugins $(PLUGINS)

$(PLUGINS):
	@$(MAKE) -C plugins/$@
	@$(MAKE) -C plugins/$@ install

bin_plugins:
	@install -m 0755 -d bin_plugins

clean:
	@rm -rf bin_plugins
	@rm -rf ${PROJECT_NAME}