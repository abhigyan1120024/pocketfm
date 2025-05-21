# PocketFM Assignment

This project is a simple Go web service using the Gin framework, containerized with Docker, and secured with JWT authentication via Envoy proxy. It includes utilities for generating RSA keys and JWT tokens.

## Project Structure

```
deployment.yaml
Dockerfile
envoy-deployment.yaml
envoy.yaml
go.mod
go.sum
kind-config.yaml
main.go
jwt_generator/
  jwt_generator.go
key_generator/
  key_generator.go
```

## Features

- REST API with Gin (`/hello`, `/health` endpoints)
- JWT authentication (RS256) enforced by Envoy
- Utilities for generating RSA keys and JWT tokens
- Kubernetes manifests for deployment
- Kind config for local Kubernetes cluster

## Getting Started

### 1. Generate RSA Keys and JWKS

To generate a new RSA key pair and JWKS:

```sh
cd key_generator
go run key_generator.go
```

This will create `private_key.pem` and `public_key.pem` in the `key_generator` directory and print a JWKS JSON string to the console.

**Important:**  
Copy the printed JWKS JSON string and update the `envoy.yaml` file under the `inline_string` field for `local_jwks`.  
For example, replace:

```yaml
inline_string: |
  {
    "keys": [
      {
        ...
      }
    ]
  }
```

with the new JWKS output.

### 2. Generate a JWT Token

```sh
cd jwt_generator
go run jwt_generator.go
```

This will output a JWT token signed with your private key.

### 3. Kubernetes Deployment

> **Note:** The server Docker image is already built and pushed to Docker Hub. (abhigyan04017/pocket-fm:v1.1) 
> You do **not** need to rebuild the image or modify `deployment.yaml`.

#### Install Kind

If you don't have [Kind](https://kind.sigs.k8s.io/) installed, you can install it using Homebrew:

```sh
brew install kind
```

Or download it manually:

```sh
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.23.0/kind-darwin-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

- Use [kind-config.yaml](kind-config.yaml) to create a local Kind cluster:

  ```sh
  kind create cluster --config kind-config.yaml
  ```

- Create the namespaces for Envoy and the Gin service:

  ```sh
  kubectl create namespace envoy-ns
  kubectl create namespace pocket-fm
  ```

- Deploy the Gin app and Envoy proxy:

  ```sh
  kubectl create configmap envoy-config --from-file=envoy.yaml -n envoy-ns
  kubectl apply -f deployment.yaml
  kubectl apply -f envoy-deployment.yaml
  ```

### 4. Accessing the Service

- The Gin app is exposed internally on port 8000.
- Envoy is exposed on NodePort 30010. Access the API via Envoy:

  ```
  curl -H "Authorization: Bearer <your-jwt-token>" http://localhost:30010/hello
  ```

## Endpoints

- `GET /hello` — Returns `{"message": "Hello, World!"}`
- `GET /health` — Returns `{"status": "ok"}`