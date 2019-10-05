.PHONY:all
all: build_debug run

.PHONY:run
run:
	./nogo45

.PHONY:build_debug
build_debug: clean
	go build -tags=debug .

.PHONY:clean
clean:
	rm -rf nogo45
