APP = auth-proxy
USER := $(shell whoami)

install:
	go get ./...

test:
	go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

tidy:
	go mod tidy

verify:
	go mod verify

.PHONY: \
	install \
	test \
	tidy \
	verify