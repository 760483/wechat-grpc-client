GO_CMD=go
API_BIN=wxapi-server
API_CLEAR_BIN=clear-server
API_BIN_PATH=./output

.PHONY: clear-server
clear-server:
	@echo "install clear-server start >>>"
	mkdir -p $(API_BIN_PATH)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build -o $(API_CLEAR_BIN) ./clearUser.go
	cp ./$(API_CLEAR_BIN) $(API_BIN_PATH)
	rm -rf ./$(API_CLEAR_BIN)
	@echo ">>> install clear-server complete"

.PHONY: linux
linux:
	@echo "install [linux] webchat-server start >>>"
	mkdir -p $(API_BIN_PATH)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build -o $(API_BIN) ./main.go
	cp ./$(API_BIN) $(API_BIN_PATH)
	rm -rf ./$(API_BIN)
	@echo ">>> install [linux] webchat-server complete"

.PHONY: mac
mac:
	@echo "install [mac] webchat-server start >>>"
	mkdir -p $(API_BIN_PATH)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO_CMD) build -o $(API_BIN) ./main.go
	cp ./$(API_BIN) $(API_BIN_PATH)
	rm -rf ./$(API_BIN)
	@echo ">>> install [mac] webchat-server complete"

.PHONY: win
win:
	@echo "install [windows] webchat-server start >>>"
	mkdir -p $(API_BIN_PATH)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_CMD) build -o $(API_BIN).exe ./main.go
	cp ./$(API_BIN).exe $(API_BIN_PATH)
	rm -rf ./$(API_BIN).exe
	@echo ">>> install [windows] webchat-server complete"

.PHONY: clean
clean:
	@echo "clean start >>>"
	rm -fr $(API_BIN_PATH)
	@echo ">>> clean complete"