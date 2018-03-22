GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_SERVER=perfchat_server
SERVER_SRC=./cmd/server
CLIENT_SRC=./cmd/client
BINARY_CLIENT=perfchat_client
LOCAL_SERVER_CONF=./deployment/local/server.conf
LOCAL_CLIENT_CONF=./deployment/local/client.conf
SERVER_CONF=server.conf
CLIENT_CONF=client.conf

.PHONY: all test build server client clean

all: test build

build: server client

local: test build
	@cp $(LOCAL_CLIENT_CONF) $(CLIENT_CONF)
	@cp $(LOCAL_SERVER_CONF) $(SERVER_CONF)

server:
	@echo "-> Building server"
	@$(GOBUILD) -o $(BINARY_SERVER) $(SERVER_SRC)

client:
	@echo "-> Building client"
	@$(GOBUILD) -o $(BINARY_CLIENT) $(CLIENT_SRC)

test:
	$(GOTEST) -v ./...

clean:
	@echo "-> Cleaning"
	@$(GOCLEAN)
	@rm -f $(BINARY_SERVER)
	@rm -f $(BINARY_CLIENT)
	@rm -f $(CLIENT_CONF)
	@rm -f $(SERVER_CONF)