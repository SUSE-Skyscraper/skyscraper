# SUSE Skyscraper

## Development Environment

### Prerequisites

* skaffold
* minikube (with [ingress-dns](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/))
* kubectl
* make
* cfssl
* helm
* nodejs 14+
* golang 1.18+

**openSUSE Tumbleweed:**

```bash
sudo zypper in skaffold minikube kubernetes-client helm nodejs16 npm16 go1.18 make cfssl
```

### Configuration

1. Generate the development TLS CA certificates:
    ```bash
    make dev_ca
    ```
2. (optional) Add `kubernetes/ca/keys/ca.pem` to your browser.
3. Build the backend configuration file `kubernetes/config/config.yaml` with `kubernetes/config/config.yaml.example` as an example.
4. Build the frontend configuration file `kubernetes/config/environmnent.local.ts` with `kubernetes/config/environmnent.local.ts.example` as an example.
5. Generate the Kubernetes configuration files:
    ```bash
    make dev_config
    ```
6. Start minikube:
    ```bash
    minikube start
    ```
7. Start skaffold:
    ```bash
    skaffold dev
    ```
8. If everything is set up properly (including [ingress-dns](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/)), you should now be able to reach Skyscraper at https://skyscraper-web.test.
