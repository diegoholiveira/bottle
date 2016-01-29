BIN_PATH = $(GOPATH)/bin/bottle

bottle: clean
	go install

clean:
	rm -f $(BIN_PATH)
