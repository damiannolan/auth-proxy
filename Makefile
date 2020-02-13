APP = auth-proxy
USER := $(shell whoami)

deploy:
	docker-compose up

destroy:
	docker-compose down

install:
	go get ./...

test:
	go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

tidy:
	go mod tidy

verify:
	go mod verify

.PHONY: \
	deploy \
	destroy \
	install \
	test \
	tidy \
	verify