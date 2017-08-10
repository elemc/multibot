PROJECT_NAME:=multibot
GO_GET:=go get

all: multibot bin_plugins/log_messages.so bin_plugins/save_messages.so

multibot: deps
	@echo "Building ${PROJECT_NAME}"
	@go build

deps:
	${GO_GET} gopkg.in/telegram-bot-api.v4
	${GO_GET} github.com/sirupsen/logrus
	${GO_GET} github.com/spf13/viper
	${GO_GET} github.com/go-pg/pg

bin_plugins/log_messages.so: bin_plugins
	@echo "Building plugin: log_messages..."
	@go build -buildmode=plugin -o bin_plugins/log_messages.so plugins/log_messages/main.go

bin_plugins/save_messages.so: bin_plugins
	@echo "Building plugin: save_messages..."
	@go build -buildmode=plugin -o bin_plugins/save_messages.so plugins/save_messages/main.go

bin_plugins:
	@install -m 0755 -d bin_plugins
