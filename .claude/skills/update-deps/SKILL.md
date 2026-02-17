# Skill: Update Go Dependencies

Обновление Go-зависимостей в gateway-проекте.

## Важно

- `github.com/demo/contracts` подключен через `replace ../contracts` — он **не обновляется** через `go get`
- Критические пакеты (требуют особого внимания при обновлении): `connectrpc.com/connect`, `google.golang.org/protobuf`
- Приватных зависимостей нет, `GOPRIVATE` не нужен
- Кодогенерации нет, `task gen` не нужен

## Step 0: Pre-flight Checks

Проверить что:
1. Текущая ветка **не** `main` — если `main`, предупредить пользователя и спросить как продолжить
2. Нет uncommitted changes — если есть, предупредить и спросить пользователя

```bash
git branch --show-current
git status --porcelain
```

## Step 1: Show Available Updates

Показать доступные обновления:

```bash
cd /Users/fdorov/projects/demo/gateway && go list -m -u all 2>/dev/null | grep '\[.*\]'
```

Пояснить пользователю:
- Какие обновления доступны (major/minor/patch)
- Что `github.com/demo/contracts` управляется через `replace` директиву и не обновляется через `go get`
- Отметить обновления критических пакетов (`connectrpc.com/connect`, `google.golang.org/protobuf`)

**Спросить пользователя** какую стратегию обновления выбрать (Step 2).

## Step 2: Update Dependencies

Три стратегии на выбор пользователя:

### All (minor + patch)
```bash
cd /Users/fdorov/projects/demo/gateway && go get -u ./...
go mod tidy
```

### Patch only (безопасное обновление)
```bash
cd /Users/fdorov/projects/demo/gateway && go get -u=patch ./...
go mod tidy
```

### Specific packages
```bash
cd /Users/fdorov/projects/demo/gateway && go get package@version
go mod tidy
```

## Step 3: Show Changes

Показать что изменилось:

```bash
cd /Users/fdorov/projects/demo/gateway && git diff go.mod go.sum
```

## Step 4: Verification Pipeline

Запустить проверки последовательно. При ошибке на любом шаге — остановиться и перейти к Step 5.

```bash
cd /Users/fdorov/projects/demo/gateway && task fmt
cd /Users/fdorov/projects/demo/gateway && task lint
cd /Users/fdorov/projects/demo/gateway && task build
cd /Users/fdorov/projects/demo/gateway && task test
```

Если все проверки прошли — перейти к Step 6.

## Step 5: Handle Failures

Если верификация провалилась, предложить пользователю три варианта:

1. **Fix** — попытаться исправить проблему (обновить код, откатить конкретную зависимость)
2. **Rollback** — откатить все изменения: `git checkout go.mod go.sum`
3. **Abort** — оставить как есть, пользователь разберётся сам

## Step 6: Commit

**Только по явному запросу пользователя.** Не коммитить автоматически.

Попытаться извлечь номер тикета из имени ветки (например `FEAT-123-update-deps` → `FEAT-123`).

Формат коммита:
```
[TICKET-ID] update go dependencies

Updated packages:
- package1: v1.0.0 → v1.1.0
- package2: v2.3.0 → v2.3.1
```

Если тикет не найден — коммит без префикса.
