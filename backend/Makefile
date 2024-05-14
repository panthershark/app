.PHONY: clean update gen test build install run purge migrate
all: gen update test

ifeq ($(PREFIX),)
    PREFIX := /opt/some-app-name
endif

clean:
	go clean --testcache -cache 

update:
	go mod tidy && go mod vendor
	
gen: 
	go generate ./...

test: 
	go test ./...

build:
	go build ./cmd/server/main.go ./server

install: server
	install -m 755 server $(PREFIX)/bin/
	install -m 644 etc/.env $(PREFIX)/etc/

run:
	go run ./cmd/server/main.go

migrate: 
	dbmate up

purge:
	rm -rf $(PREFIX)