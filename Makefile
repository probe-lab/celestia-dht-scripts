# Compiler Variables 
GOCC=go
TARGET_PATH=./cmd/cnames
BIN_PATH=./build
BIN=./build/cnames


# Make Operations
.PHONY: install uninstall build clean tidy audit test docker 

install:
	$(GOCC) install $(GIT_PACKAGE)

uninstall:
	$(GOCC) clean $(GIT_PACKAGE)

build:
	$(GOCC) get $(TARGET_PATH)
	$(GOCC) build -o $(BIN) $(TARGET_PATH)

clean:
	rm -r $(BIN_PATH)

tidy:
	gofumpt -w -l .
	$(GOCC) mod tidy -v

audit:
	$(GOCC) mod verify
	$(GOCC) vet ./...
	$(GOCC) run honnef.co/go/tools/cmd/staticcheck@latest ./...
	$(GOCC) test -race -buildvcs -vet=off $(TARGET_PATH)