# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Communication Guidelines

- Write all code and code comments in English
- Communicate with the user in Russian
- Maintain professional style in both languages

### Agents for Testing

**Always use `gateway-test-coordinator`** to write tests — it distributes work between unit and isolation agents in parallel, avoiding duplication.

Available agents:
- **`gateway-test-coordinator`** — **USE THIS** for writing tests (delegates to agents below)
- **`gateway-unit-test-writer`** - for unit tests (validation, business logic, error handling with mocks)
- **`gateway-isolation-test-writer`** - for end-to-end integration tests (full flows with real infrastructure)

## Build and Run

Проект использует [Task](https://taskfile.dev/) для автоматизации команд.

```bash
# Запуск (порт 8080)
task run

# Сборка в ./bin/gateway
task build

# Форматирование + линтинг + сборка
task

# Только форматирование
task fmt

# Только линтинг (golangci-lint)
task lint

# Тесты
task test

# Установка инструментов (golangci-lint)
task install-tools
```

**Prerequisite:** Order Service must be running on localhost:8081.

## Architecture

This is a Connect RPC gateway service that proxies requests to an Order Service backend.

### Key Components

- **cmd/app/main.go** - Entry point. Creates the Connect RPC client for the backend Order Service and registers the gateway server with HTTP/2 (h2c) support.

- **internal/domain/orders/server.go** - Implements `orderv1connect.OrderServiceHandler` interface. Delegates each RPC method to a dedicated handler.

- **internal/domain/orders/grpc_*_handler.go** - Individual handlers for each RPC method. Each handler:
  - Has a `Handle(ctx, req)` method
  - Contains request validation in a `validate()` method
  - Proxies to the backend client after validation

### Handler Pattern

Each Connect RPC method follows this pattern:
1. Create a handler struct with the backend client
2. Implement `Handle()` that validates then proxies
3. Implement `validate()` for request validation
4. Return `connect.NewError(connect.CodeInvalidArgument, nil)` for validation failures

### Dependencies

Uses `github.com/demo/contracts` (local replace directive to `../contracts`) for:
- Proto-generated types: `github.com/demo/contracts/gen/go/order/v1`
- Connect client/handler: `github.com/demo/contracts/gen/go/order/v1/orderv1connect`
