PROJECT_NAME:=fwego
GO_GET:="@go get"

all: multibot bin_plugins/log_messages.so bin_plugins/save_messages.so

multibot: deps
	@echo "Building ${PROJECT_NAME}"
	@go build

deps:
	${GO_GET} gopkg.in/telegram-bot-api.v4
	${GO_GET} github.com/sirupsen/logrus
	${GO_GET} github.com/spf13/viper

bin_plugins/log_messages.so:
	@pushd plugins/log_messages
	go build
	install -m 0644 log_messages.so ../../bin_plugins/
	@popd

bin_plugins/save_messages.so:
	@pushd plugins/save_messages
	go build
	install -m 0644 save_messages.so ../../bin_plugins/
	@popd