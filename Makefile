PROJECT_NAME:=multibot
GO_GET:=go get
PLUGINS:=reminder filer

all: multibot subdirs

multibot: deps
	@echo "Building ${PROJECT_NAME}"
	@go build

deps:
	@echo "Installing dep..."
	@$(GO_GET) -u github.com/golang/dep/cmd/dep
	@echo "Download dependencies..."
	@dep ensure

subdirs: bin_plugins $(PLUGINS)

$(PLUGINS):
	@$(MAKE) -C plugins/$@
	@$(MAKE) -C plugins/$@ install

bin_plugins:
	@install -m 0755 -d bin_plugins

clean:
	@rm -rf bin_plugins
	@rm -rf ${PROJECT_NAME}