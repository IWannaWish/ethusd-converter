# ethusd-converter

`ethusd-converter` — это pet-проект на Go, предназначенный для практики разработки отказоустойчивых микросервисов в экосистеме Web3.  
Сервис получает on-chain балансы ETH, WETH и популярных ERC-20 токенов по Ethereum-адресу и переводит их в доллары США по курсам Chainlink.  
Результаты доступны через CLI и gRPC API.

Проект использует production-ориентированную архитектуру: с gRPC, кэшами (in-memory + Redis), брокером сообщений (NATS) и метриками Prometheus.

---

## 📌 Особенности

- CLI и gRPC интерфейсы
- Получение on-chain балансов ETH и токенов
- Chainlink Price Feeds (`latestRoundData()` on-chain)
- 2-уровневый TTL-кэш: `go-cache` + Redis
- Очередь задач через NATS
- Прометей-метрики: время ответа, ошибки, кол-во запросов
- Makefile, Docker-образ, Linter + Staticcheck

---

## 🗂 Структура проекта

- `cmd/cli` — CLI-приложение
- `internal/api` — gRPC-сервер
- `internal/core` — бизнес-логика
- `internal/cache` — реализация 2-level TTL cache
- `internal/broker` — взаимодействие с брокером (NATS)
- `internal/eth` — взаимодействие с Ethereum, контракты
- `proto/` — protobuf-схемы

---

## 🚀 Быстрый старт

### Требования

- Go 1.22+
- `protoc`
- gRPC + protobuf-плагины:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  ```

- Доп. инструменты (опционально):
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  go install honnef.co/go/tools/cmd/staticcheck@latest
  ```

---

## 🛠 Установка

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/yourname/ethusd-converter.git
   cd ethusd-converter
   ```

2. Установите зависимости:

   ```bash
   make tidy
   ```

---

## ⚙️ Сборка и запуск

### CLI

```bash
make run
```

### Docker

Сборка и запуск:

```bash
make docker-build
make docker-run
```

---

## 💻 Пример CLI-запуска

```bash
./ethusd-converter 0x1234567890abcdef...
```

Вывод:
```
Address: 0x1234...abcd
ETH:   1.245 ETH  ≈ $4,312.90
WETH:  0.875 WETH ≈ $3,032.00
DAI:   1500 DAI   ≈ $1,500.00

Total: ≈ $8,844.90
```

---

## ✅ Тестирование

Полный запуск всех тестов:

```bash
make test
```

Только unit-тесты:

```bash
make test-only
```

---

## 🎯 Цель проекта

Проект используется как pet-проект для собеседований и практики разработки микросервисов в стиле "production-grade Go":  
работа с Ethereum, взаимодействие с контрактами, gRPC, кэширование, брокеры сообщений, метрики и т.д.

---

## 📄 License

Проект распространяется под лицензией MIT. Подробности в [LICENSE](LICENSE).
