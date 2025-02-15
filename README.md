# Тестовое задание для стажера-разработчика

## Задача

Реализовать сервис, предоставляющий API для создания сокращённых ссылок.

Ссылка должна быть:

- **Уникальной**: на один оригинальный URL должна ссылаться только одна сокращённая ссылка.
- **Длиной 10 символов**.
- **Состоять из**: символов латинского алфавита в нижнем и верхнем регистре, цифр и символа `_` (подчёркивание).

## Требования к сервису

Сервис должен быть написан на **Go** и принимать следующие HTTP-запросы:

1. **POST** – сохраняет оригинальный URL в базе и возвращает сокращённый.
2. **GET** – принимает сокращённый URL и возвращает оригинальный.

### Дополнительное условие (*по желанию*)

- Реализовать работу сервиса через **gRPC**: составить `.proto`-файл и реализовать сервис с двумя соответствующими
  эндпойнтами.

## Требования к реализации

- **Docker-образ**: сервис должен быть распространён в виде Docker-образа.
- **Хранилище**: поддержка **in-memory**-решения и **PostgreSQL**. Выбор хранилища осуществляется параметром при
  запуске.
- **Тестирование**: реализованный функционал должен быть покрыт **unit-тестами**.

## Ожидаемый результат

Решение предоставить в виде **публичного репозитория** на [GitHub](https://github.com).

## Критерии оценки

1. **Алгоритм генерации ссылок**
    - Как генерируются ссылки и почему предложенный алгоритм будет работать.
    - Соответствие требованиям и простота понимания.

2. **Структура проекта**
    - Логика разбиения типов по файлам, файлов по пакетам и пакетов по приложению.

3. **Обработка ошибок**
    - Как сервис обрабатывает ошибки в различных сценариях.

4. **Удобство и логика использования**
    - Насколько сервис удобен и логичен в эксплуатации.

5. **Масштабируемость**
    - Как сервис ведёт себя при высокой нагрузке (сотни пользователей одновременно, как, например, YouTube или ya.cc).

6. **Долговременная работа**
    - Что произойдёт, если сервис будет работать длительное время.

7. **Чистота кода**
    - Общий уровень качества кода.

## Решение

### Конфигурация

Настройки сервиса хранятся в файле `config/config.yml`, где указываются порты и сопутствующие параметры. По умолчанию:

- HTTP-сервер работает на порту **8080**
- gRPC-сервер работает на порту **50051**

#### Пример конфигурации

```yaml
server:
  http_port: ":8080"
  grpc_port: ":50051"
  timeout: "10s"
  idle_timeout: "15s"

storage:
  type: "postgres" # memory, postgres
  postgres:
    host: "postgres"
    port: 5432
    user: "postgres"
    password: "password"
    dbname: "shortener"

log:
  level: "prod" # local, prod
```

### Как работает In-Memory хранилище

In-Memory хранилище реализовано в пакете `memory`. Оно использует два `map` для хранения данных:

- `storage` для хранения соответствия `shortURL -> originalURL`
- `reverse` для хранения обратного соответствия `originalURL -> shortURL`

За счёт двух мап мы можем за O(1) проверять кейс на уже записаный rl

При добавлении нового URL выполняется проверка на существование, после чего данные записываются в обе карты с блокировкой `sync.RWMutex` для потокобезопасности.

#### Пример кода

```go
func (s *StorageInMemory) Put(url, shortURL string) error {
    s.rvMu.Lock()
    defer s.rvMu.Unlock()

    if _, ok := s.storage[shortURL]; ok || _, ok := s.reverse[url]; ok {
        return errs.ErrURLIsExist
    }

    s.storage[shortURL] = url
    s.reverse[url] = shortURL
    return nil
}
```

### Как работает генератор случайных строк

Генерация случайных коротких URL выполняется в пакете `random`.
Используется криптографически безопасный генератор случайных чисел из пакета `crypto/rand`. Для генерации строки случайно выбираются символы из `alphabet`, содержащего буквы, цифры и `_`.

#### Пример кода

```go
func NewRandomString(stringLength int) (string, error) {
    if stringLength <= 0 {
        return "", ErrInvalidLength
    }

    var builder strings.Builder
    alphaLength := big.NewInt(int64(utf8.RuneCount(alphabet)))

    for range stringLength {
        n, err := rand.Int(rand.Reader, alphaLength)
        if err != nil {
            return "", fmt.Errorf("failed to generate random integer: %w", err)
        }
        builder.WriteByte(alphabet[n.Int64()])
    }
    return builder.String(), nil
}
```

### API

#### HTTP

##### Сокращение ссылки

**Запрос:**

- **Метод:** `POST`
- **Эндпоинт:** `/shorten`
- **Тело запроса (JSON):**

  ```json
  {
    "url": "https://example.com"
  }
  ```

**Ответ:**

- **200 OK**

  ```json
  {
    "short_url": "example",
    "status": "OK"
  }
  ```

- **500 Internal Server Error** (если URL уже существует)

  ```json
  {
    "error": "url already exists",
    "status": "Error"
  }
  ```

- **400 Bad Request** (если URL невалидный)

  ```json
  {
    "error": "invalid URL format",
    "status": "Error"
  }
  ```

##### Получение оригинальной ссылки

**Запрос:**

- **Метод:** `POST`
- **Эндпоинт:** `/resolve`
- **Тело запроса (JSON):**

  ```json
  {
    "short_url": "example"
  }
  ```

**Ответ:**

- **200 OK**

  ```json
  {
    "original_url": "https://example.com",
    "status": "OK"
  }
  ```

- **404 Not Found** (если ссылка не найдена)

  ```json
  {
    "error": "URL not found",
    "status": "Error"
  }
  ```

#### gRPC

Файл спецификации: `proto/urlshortener.proto`

```protobuf
syntax = "proto3";

package urlshortener;

option go_package = "../internal/grpc/urlshortener";

service URLShortener {
  rpc Shorten (ShortenRequest) returns (ShortenResponse);
  rpc Resolve (ResolveRequest) returns (ResolveResponse);
}

message ShortenRequest {
  string url = 1;
}

message ShortenResponse {
  string short_url = 1;
}

message ResolveRequest {
  string short_url = 1;
}

message ResolveResponse {
  string original_url = 1;
}
```

### Тестирование

Функционал сервиса покрыт юнит-тестами. Для запуска тестов:

```sh
go test -v ./...
```

Покрытие:

```sh
        url-shortener/cmd/url-shortener         coverage: 0.0% of statements
        url-shortener/cmd/url-shortener/server/httpserver               coverage: 0.0% of statements
        url-shortener/cmd/url-shortener/server/grpcserver               coverage: 0.0% of statements
        url-shortener/internal/config           coverage: 0.0% of statements
        url-shortener/internal/grpc/urlshortener                coverage: 0.0% of statements
ok      url-shortener/internal/grpc/server      0.120s  coverage: 92.3% of statements
?       url-shortener/internal/storage/errs     [no test files]
        url-shortener/internal/http/middleware/mvlogger         coverage: 0.0% of statements
ok      url-shortener/internal/http/handlers/resolve    0.472s  coverage: 100.0% of statements
        url-shortener/internal/logger           coverage: 0.0% of statements
        url-shortener/internal/storage          coverage: 0.0% of statements
        url-shortener/internal/storage/postgres         coverage: 0.0% of statements
ok      url-shortener/internal/http/handlers/shorten    0.498s  coverage: 84.2% of statements
ok      url-shortener/internal/service  0.378s  coverage: 83.3% of statements
ok      url-shortener/internal/storage/memory   0.426s  coverage: 94.1% of statements
ok      url-shortener/pkg/util/random   0.210s  coverage: 90.0% of statements
```

## Запуск

### Через Docker

```sh
docker-compose up --build
```

### Локально

```sh
go run cmd/url-shortener/main.go
```
