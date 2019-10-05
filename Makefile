.PHONY:all
all: build_debug run

.PHONY:run
run:
	./nogo45

.PHONY:build_debug
build_debug: *.go assets/*
	go build -tags=debug .

.PHONY:clean
clean:
	rm -rf dist/
	rm -rf nogo45
	rm -rf nogo45-linux-amd64
	rm -rf nogo45-windows-4.0-amd64.exe
	rm -rf nogo45-darwin-10.6-amd64

.PHONY:build
build: *.go assets/* clean nogo45-linux-amd64 nogo45-windows-4.0-amd64.exe nogo45-darwin-10.6-amd64
	mkdir -p dist

.PHONY:dist
dist: build
	zip dist/windows.zip nogo45-windows-4.0-amd64.exe -r assets
	zip dist/mac.zip nogo45-darwin-10.6-amd64 -r assets
	zip dist/linux.zip nogo45-linux-amd64 -r assets

nogo45-linux-amd64: *.go assets/*
	GOOS=linux GOARCH=amd64 go build -o nogo45-linux-amd64 .

nogo45-windows-4.0-amd64.exe: *.go assets/*
	xgo -go 1.12 --targets=windows/amd64 -ldflags='-H=windowsgui' github.com/bcvery1/nogo45

nogo45-darwin-10.6-amd64: *.go assets/*
	xgo -go 1.12 --targets=darwin/amd64 github.com/bcvery1/nogo45
