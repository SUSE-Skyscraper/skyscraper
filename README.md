# SUSE Skyscraper

## Development Environment

### Prerequisites

**Local Environment:**

* PostgresSQL
* nats
* make
* nodejs 14+
* golang 1.18+
* sqlc
   * `go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`

**Minikube Environment:**

* skaffold
* minikube (with [ingress-dns](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/))
* kubectl
* cfssl
  * `go install github.com/cloudflare/cfssl/cmd/...@latest`
* helm

**openSUSE Tumbleweed:**

```bash
sudo zypper in skaffold minikube kubernetes-client helm nodejs16 npm16 go1.18 make
```

### API Specification

The api specification is located at [api/skyscraper.yaml](api/skyscraper.yaml)

### Database Migrations

The database migration files are at [cmd/app/migrate/migrations](cli/cmd/app/migrate/migrations). They're embedded into the binary, and we read them in the `migrate` command.

**Migrate Up:**
```bash
go run ./cmd/main.go migrate up
```

**Migrate Down:**
```bash
go run ./cmd/main.go migrate down
```

### Generate Database files

**Notes:**

* The queries are located in the `queries.sql` file.
* The config file for `sqlc` is located at `sqlc.yaml`.
* The database files are generated at `internal/db`.

Run the following command to generate the database files:

```bash
sqlc generate
```

### Deploy Locally

1. Ensure that you have a PostgresSQL server that you can connect to locally.
2. Ensure that you have a NATS server that you can connect to locally:
   ```bash
   docker run -p 4222:4222 --name nats nats:latest -js
   ```
3. Copy `config.yaml.example` to `config.yaml` and fill in the values.
4. Copy `web/src/environments/environment.ts` to `web/src/environments/environment.local.ts` and fill in its values.
5. Build the golang backend:
   ```bash
   go build ./cmd/main.go
   ```
6. Run database migrations:
   ```bash
   go run ./cmd/main.go migrate up
   ```
7. Run the sync job:
   ```bash
   go run ./cmd/main.go cloud-sync
   ```
8. Run the server:
   ```bash
   go run ./cmd/main.go server
   ```
9. While the server is running, start the web frontend:
   ```bash
   cd web
   npm start
   ```
10. The frontend should be live at [http://localhost:4200](http://localhost:4200).

### Deployment on Minikube

1. Generate the development TLS CA certificates:
    ```bash
    make dev_ca
    ```
2. Add helm repos:
    ```
    helm repo add nats https://nats-io.github.io/k8s/helm/charts/
    helm repo add skyscraper https://suse-skyscraper.github.io/skyscraper-helm-charts
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo update
    ```
3. Copy helm configuration values from `kubernetes/helm/*-values.yaml.example` to respective `*-values.yaml` and make changes to customize your environment.
4. (optional) Add `kubernetes/ca/keys/ca.pem` to your browser.
5. Build the backend configuration file `kubernetes/config/config.yaml` with `kubernetes/config/config.yaml.example` as an example.
6. Build the frontend configuration file `kubernetes/config/environmnent.local.ts` with `kubernetes/config/environmnent.local.ts.example` as an example.
7. Generate the Kubernetes configuration files:
    ```bash
    make dev_config
    ```
8. Start minikube:
    ```bash
    minikube start
    ```
9. Start skaffold:
    ```bash
    skaffold dev
    ```
10. If everything is set up properly (including [ingress-dns](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/)), you should now be able to reach Skyscraper at https://skyscraper-web.test.
