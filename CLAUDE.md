# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Communication Guidelines

- Write all code and code comments in English
- Communicate with the user in Russian
- Maintain professional style in both languages

### Agents for Testing

⚠️ **ЗАПРЕЩЕНО писать тесты напрямую. ВСЕГДА используй агента через Task tool.**

Когда нужно написать тесты:
1. **СТОП** — не пиши тест-код сам
2. **Используй Task tool** с `subagent_type: "gateway-unit-test-writer"` или `"gateway-test-coordinator"`
3. Агент знает все паттерны проекта (testify, GWT, connect.CodeOf и т.д.)

Available agents:
- **`gateway-test-coordinator`** — для комплексного покрытия (распределяет между unit и isolation)
- **`gateway-unit-test-writer`** — для unit тестов (validation, mocks, error handling)
- **`gateway-isolation-test-writer`** — для e2e тестов (full flows with real infrastructure)

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
