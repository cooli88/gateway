# Gateway Service

REST API шлюз, который проксирует запросы к Order Service через Connect RPC.

## Запуск

```bash
go run .
```

Сервис запустится на порту 8080.

## API

### Health Check

```
GET /health
```

### Orders

```
POST /orders - создать заказ
GET /orders - получить список заказов
GET /orders/:id - получить заказ по ID
```

## Пример использования

```bash
# Создать заказ
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user1", "item": "MacBook Pro", "amount": 2499.99}'

# Список заказов
curl http://localhost:8080/orders

# Заказ по ID
curl http://localhost:8080/orders/{id}
```

## Зависимости

- `github.com/demo/contracts` - proto-контракты и Connect клиент
- Order Service должен быть запущен на localhost:8081
