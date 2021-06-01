BIN=./bin
LIST=$(BIN)/trm

build: $(LIST)

clean:
	rm -f $(LIST)

$(BIN)/%:
	go build -o $@