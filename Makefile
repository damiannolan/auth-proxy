VERSION := 0.1.0

DOCKER_REG = damiannolan
DOCKER_IMAGE = auth-proxy
DOCKER_IMAGE_TAG = $(VERSION)
USER := $(shell whoami)
	
deploy:
	docker-compose up

destroy:
	docker-compose down

docker-build:
	docker build -t $(DOCKER_REG)/$(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG) .

docker-build-dev:
	docker build -t $(DOCKER_REG)/$(DOCKER_IMAGE):$(USER) .

docker-push:
	docker push $(DOCKER_REG)/$(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG)

docker-push-dev:
	docker push $(DOCKER_REG)/$(DOCKER_IMAGE):$(USER)

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
	docker-build \
	docker-build-dev \
	docker-push \
	docker-push-dev \
	install \
	test \
	tidy \
	verify