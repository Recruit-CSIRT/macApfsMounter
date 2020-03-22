APPNAME := macApfsMounter

.PHONY: deps
deps:
	export GO111MODULE=on
	go mod -d
	go mod tidy

.PHONY: cli
cli:
	export GO111MODULE=on
	mkdir -p build
	go build -o ./build/amtr ./cmd/cli/

.PHONY: gui
gui:
	export QT_HOMEBREW=true
	export GO111MODULE=off
	qtdeploy build desktop ./cmd/$(APPNAME)/

	mkdir -p build
	rm -rf ./build/$(APPNAME).app

	cp -R ./automator/$(APPNAME).app ./build/
	mkdir -p ./build/$(APPNAME).app/Contents/apps/
	cp -R ./cmd/$(APPNAME)/deploy/darwin/$(APPNAME).app ./build/$(APPNAME).app/Contents/apps/$(APPNAME).app

.PHONY: clean
clean:
	rm -rf ./cmd/$(APPNAME)/deploy
	rm -rf ./cmd/$(APPNAME)/darwin
	rm -rf ./build/*