# 🌀 TON Node gRPC Wrapper

A production-ready gRPC microservice that wraps **TON API v2**,
providing secure and efficient access to all essential blockchain data and operations.

This service is tailored for **Wallet-as-a-Service platforms** and similar infrastructures.

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-ready-blue)
![Architecture](https://img.shields.io/badge/architecture-clean--arch-informational)

---

## 🚀 Features

* ✅ gRPC interface for integration with internal services
* 🔒 Built-in **authorization** middleware
* 📈 Native **Prometheus** metrics (default port: `:9090`)
* 🧾 Structured logging with contextual request tracing
* 💥 Supports:
    * Account and balance info
    * Jetton token accounts
    * Transaction traces
    * Token and native transfers

---

## 🧱 Architecture

The service follows **Clean Architecture** principles:
- Domain and use case logic is isolated from infrastructure
- Interfaces define all dependencies for inversion of control
- Handlers are organized by gRPC controllers with middleware support
- Infra contains adapters for TON API, logger, metrics

## 📂 Project Structure

```
.
├── cmd/                 # gRPC entrypoint
├── config/              # YAML config loading
├── infra/               # External deps: logger, TON API, Prometheus
├── internal/
│   ├── controller/      # gRPC handlers + middleware
│   ├── domain/          # DDD interfaces & business models
│   └── usecase/         # Application logic
├── proto/               # gRPC protobuf definitions
├── tests/               # Unit & integration tests
├── tools/               # gRPC generation configs
├── Dockerfile           # Build instructions
├── docker-compose.yaml  # Local development environment
└── Makefile             # Common development commands
```

---

## 📦 Getting Started

### Prerequisites

* Go 1.24+
* `protoc` + Go plugins
* Docker + Docker Compose

### Clone & Run

```bash
git clone https://github.com/paxyside/ton-node.git
cd ton-node

make setup && make proto

make all # for local run
make docker_run # for docker compose run
```

---

## 🛠️ gRPC Usage

### Proto

Find the proto definitions in [`proto/tonnode/ton_node.proto`](./proto/tonnode/ton_node.proto).

To generate:

```bash
make proto
```

### Example Methods

* `GetAccount` – retrieve TON account state and balance
* `GetJAccount` – get Jetton account information
* `GetSeqno` – fetch seqno of wallet address
* `GetTxTrace` – retrieve trace for a transaction
* `EmulateTxTrace` – emulate a message and get trace result
* `SendMsg` – send raw message to the TON network

---

## 🔐 Authentication

All incoming gRPC requests must include a valid token in the metadata:

```text
authorization: Bearer <your-token>
```

Configure your tokens in `config.yaml` under `app.server.auth_token`.

---

## 📊 Observability

Prometheus metrics available at:

```bash
http://localhost:9090/metrics
```

---

## 📚 Configuration

See `config.example.yaml`:

```yaml
app:

  server:
    host: "0.0.0.0"
    port: "50051"
    prometheus_host: "0.0.0.0"
    prometheus_port: "9090"
    read_header_timeout: 5s
    shutdown_timeout: 5s
    auth_token: "secret"

  node:
    url: "https://testnet.tonapi.io" # https://tonapi.io - mainnet node url, replace if needed
    api_key: "YOUR-API-KEY"
    timeout: 10s
    rate_limit: 1
    rate_burst: 5
```

---

## 🧪 Testing

```bash
go test ./... # or make tests
```

Includes:

* ❖ Unit tests for domain models and validators
* ❖ Integration tests for TON API
* ❖ Mocks and test utils for isolated execution

---

## 🧰 Makefile Commands

The project provides a Makefile with handy shortcuts for development and operations:

| Command        | Description                                                                 |
|----------------|-----------------------------------------------------------------------------|
| `make all`     | Cleans up the binary, runs linters, builds the app, and launches it locally |
| `make prepare` | Applies lint fixes, runs tests, and generates protobuf files                |
| `make lint`    | Runs static analysis using `golangci-lint`                                  |
| `make lint-fix`| Fixes formatting and minor issues automatically                             |
| `make pack`    | Compiles the binary to `./server` from `./cmd/server/main.go`               |
| `make tests`   | Runs all Go tests with verbose output                                       |
| `make local_run` | Launches the compiled binary locally (`./server`)                        |
| `make clear`   | Deletes the compiled binary                                                 |
| `make docker_run` | Builds and runs the service using Docker Compose                        |
| `make proto`   | Compiles protobuf definitions to Go using `protoc`                          |
| `make setup`   | Copies the example config and creates a Docker network `ton-node-network`   |
