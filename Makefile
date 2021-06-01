BIN=./bin
PROG=trm
LIST=$(addprefix $(BIN)/, $(PROG))

build: $(LIST)

clean:
	rm -f $(LIST)

$(BIN)/%:
	go build -o $@