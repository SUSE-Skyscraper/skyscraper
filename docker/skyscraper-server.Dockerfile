# syntax=docker/dockerfile:1.2

FROM registry.suse.com/bci/golang:1.18 as builder

RUN zypper --non-interactive up

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o skyscraper-server ./cmd/server/main.go

FROM registry.suse.com/bci/bci-base:latest

RUN zypper --non-interactive up

COPY --from=builder /app/skyscraper-server /usr/local/bin/skyscraper-server

RUN mkdir /app
WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/skyscraper-server"]
