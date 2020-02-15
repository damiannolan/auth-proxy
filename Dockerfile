# Intermediate Builder
FROM golang:1.13 AS build-env

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /opt/app/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -i -o app ./cmd/app/

# Application Image
# https://github.com/GoogleContainerTools/distroless/tree/master/base
FROM gcr.io/distroless/base:latest

COPY --from=build-env /opt/app/app /usr/local/bin/app

CMD ["/usr/local/bin/app"]

EXPOSE 8079