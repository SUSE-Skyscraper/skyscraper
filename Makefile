.ONESHELL:
.PHONY: lint fmt build dev_cert test test_unit test_integration

lint:
	golangci-lint run
	cd web
	ng lint

fmt:
	go mod tidy
	go fmt ./cmd/... ./internal/...
	cd web
	ng lint --fix

build:
	go build -v ./cmd/main.go

test: test_unit test_integration

test_unit:
	go test -v ./cmd/... ./internal/...

test_integration:
	go test -v ./test/...

dev_ca:
	cd kubernetes/ca/keys
	cfssl gencert -initca ../ca-csr.json | cfssljson -bare ca

dev_config: dev_cert
	cd kubernetes/config
	kubectl create configmap skyscraper-server-config --dry-run=client --from-file=config=config.yaml -o yaml > gen/skyscraper-server-config.yaml
	kubectl create configmap skyscraper-web-config --dry-run=client --from-file=environment=environment.local.ts -o yaml > gen/skyscraper-web-config.yaml

dev_cert:
	cd kubernetes/ca/keys
	cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ../config.json -profile=www ../skyscraper-web-csr.json | cfssljson -bare skyscraper-web
	kubectl create secret tls skyscraper-web-tls --cert=skyscraper-web.pem --key=skyscraper-web-key.pem --dry-run=client --output yaml > ../../config/gen/skyscraper-web-tls.yaml
	cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ../config.json -profile=www ../skyscraper-backend-csr.json | cfssljson -bare skyscraper-backend
	kubectl create secret tls skyscraper-backend-tls --cert=skyscraper-backend.pem --key=skyscraper-backend-key.pem --dry-run=client --output yaml > ../../config/gen/skyscraper-server-tls.yaml
