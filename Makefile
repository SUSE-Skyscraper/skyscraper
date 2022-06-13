.ONESHELL:
.PHONY: lint fmt build test dev_cert

lint:
	golangci-lint run

fmt:
	go mod tidy
	go fmt ./cmd/... ./internal/...

build:
	go build -v ./cmd/server/main.go

test:
	go test -v ./cmd/... ./internal/...

dev_ca:
	cd ca/keys
	cfssl gencert -initca ../ca-csr.json | cfssljson -bare ca

dev_cert:
	cd ca/keys
	cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ../config.json -profile=www ../skyscraper-web-csr.json | cfssljson -bare skyscraper-web
	kubectl create secret tls skyscraper-web-tls --cert=skyscraper-web.pem --key=skyscraper-web-key.pem --dry-run=client --output yaml > ../../kubernetes/skyscraper-web-tls.yaml
	cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ../config.json -profile=www ../skyscraper-backend-csr.json | cfssljson -bare skyscraper-backend
	kubectl create secret tls skyscraper-backend-tls --cert=skyscraper-backend.pem --key=skyscraper-backend-key.pem --dry-run=client --output yaml > ../../kubernetes/skyscraper-backend-tls.yaml
