OUTPUT := ./bin
BIN := friendly-bot

build:
	go build -o ${OUTPUT}/${BIN} *.go

clean:
	rm ${OUTPUT}/* || true
