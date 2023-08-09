.PHONY: mocks *server *client test

SERVER_DIR = "./server"
SERVER_BIN = "bin/server"

build-server:
	cd $(SERVER_DIR) && go build -o ../$(SERVER_BIN) ./cmd 
	
test-server:
	cd $(SERVER_DIR) && go test ./...

server: | build-server
	./$(SERVER_BIN)

CLIENT_DIR = "./client"

test-client:
	cd $(CLIENT_DIR) && npm test

client: 
	cd $(CLIENT_DIR) && npm start

test: | test-server test-client
