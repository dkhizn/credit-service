# Ипотечный калькулятор

Сервис расчёта параметров ипотеки с локальным in-memory кэшем.

## Оглавление

- [Описание](#описание)  
- [Функциональность](#функциональность)  
- [Архитектура](#архитектура)  
- [Конфигурация](#конфигурация)
- [Запуск через Makefile](#запуск-через-makefile)  
  - [1. Локальный запуск](#1-локальный-запуск)  
  - [2. Тестирование](#2-тестирование)  
  - [3. Линтинг](#3-линтинг)  
  - [4. Сборка Docker-образа](#4-сборка-docker-образа)  
  - [5. Запуск контейнера](#5-запуск-контейнера)  
  - [6. Остановка и удаление контейнера](#6-остановка-и-удаление-контейнера)  
- [REST API](#rest-api)  
  - [POST /execute](#post-execute)  
  - [GET /cache](#get-cache) 
## Описание

Этот сервис реализует ипотечный калькулятор, который на вход принимает JSON-запрос с параметрами объекта недвижимости, первоначальным взносом, сроком и программой кредитования.  

Он вычисляет:
- годовую процентную ставку (8%, 9% или 10%)  
- сумму кредита  
- аннуитетный ежемесячный платёж 
- переплату за весь срок  
- дату последнего платежа  

Результат каждого расчёта сохраняется в локальном in-memory кэше.

## Функциональность

1. **POST /execute** – выполнить расчёт ипотеки.  
   - Проверки:
     - введена ипотечная программа;  
     - введена лишь одна ипотечная программа из трёх(`salary` / `military` / `base`);
     - первоначальный взнос ≥ 20% от стоимости объекта;
   - В случае ошибок возвращается HTTP 400 + JSON с ключом `"error"`.  
   - Успешный ответ (HTTP 200) – JSON вида:
     ```json
     {
       "result": {
         "params": {
           "object_cost": 5000000,
           "initial_payment": 1000000,
           "months": 240
         },
         "program": {
           "salary": true
         },
         "aggregates": {
           "rate": 8,
           "loan_sum": 4000000,
           "monthly_payment": 33458,
           "overpayment": 4029920,
           "last_payment_date": "2044-02-18"
         }
       }
     }
     ```

2. **GET /cache** – получить массив всех ранее выполненных расчётов.  
   - Если кэш пуст, возвращается HTTP 400 + JSON:
     ```json
     {
       "error": "empty cache"
     }
     ```
   - Если есть записи, возвращается HTTP 200 + JSON-массив:  
     ```json
     [
       {
         "id": 0,
         "params": {
           "object_cost": 5000000,
           "initial_payment": 1000000,
           "months": 240
         },
         "program": {
           "salary": true
         },
         "aggregates": {
           "rate": 8,
           "loan_sum": 4000000,
           "monthly_payment": 33458,
           "overpayment": 4029920,
           "last_payment_date": "2044-02-18"
         }
       },
       ...
     ]
     ```

## Архитектура

Чистая (Clean) архитектура, слои:

- **domain**  
  `internal/domain/credit.go`  
  – доменная модель `Credit` (входные поля + поля результатов).

- **application (use-case)**  
  `internal/application/credit-service/`  
  – `CredtiService`:
  - `Execute(*Credit)` – расчёты, сохранение в кэш;
  - `GetAll() []Credit` – вернуть все сохранённые кредиты.

- **adapters / primary (HTTP-adapter)**  
  `internal/adapters/primary/http-adapter/controller/`  
  – контроллеры:
  
  - `Post(w,r)` → DTO → `ToDomain` → `svc.Execute` → DTO ответа;
  - `Cache(w,r)` → `svc.GetAll()` → формирование JSON-списка.  
  `internal/adapters/primary/http-adapter/router/router.go` – регистрация маршрутов.

- **adapters / secondary (репозиторий)**  
  `internal/adapters/secondary/repositories/credit-repository/`  
  – `СreditRepository` хранит срез `[]domain.Credit` + `sync.RWMutex`.

- **cmd / main.go**  
  Точка входа: создание контроллера, роутера, middleware, запуск HTTP-сервера.

## Конфигурация

Файл `config.yml` в корне:

```yaml
port: 8080
```


## Запуск через Makefile

Проект полностью управляется через `Makefile`.

### 1. Локальный запуск

```bash
make run
```

Запускает Docker-контейнер, публикует порт `8080`.

### 2. Тестирование

```bash
make test
```

Запускает юнит-тесты. Покрытие:

- Контроллеры: ~92%
- Middleware и роутер: 100%
- Репозиторий: 100%
- Бизнес-логика (CreditService): ~94%
- Конфигурация: 100%

### 3. Линтинг

```bash
make lint
```

Анализ кода через `golangci-lint`.

### 4. Сборка Docker-образа

```bash
make build
```

Собирает Docker-образ с бинарником `credit-service`.

### 5. Запуск контейнера

```bash
make run
```

Запускает сервис в Docker-контейнере.

### 6. Остановка и удаление контейнера

```bash
make stop
```

Останавливает и удаляет контейнер.  
Для полной очистки:

```bash
make clean
```

## REST API

### POST /execute

Запрос:

```json
{
  "params": {
    "object_cost": 5000000,
    "initial_payment": 1000000,
    "months": 240
  },
  "program": {
    "salary": true
  }
}
```

Ответ:

```json
{
  "result": {
    "params": {
      "object_cost": 5000000,
      "initial_payment": 1000000,
      "months": 240
    },
    "program": {
      "salary": true
    },
    "aggregates": {
      "rate": 8,
      "loan_sum": 4000000,
      "monthly_payment": 33458,
      "overpayment": 4029920,
      "last_payment_date": "2044-02-18"
    }
  }
}
```

Ошибки:

```json
{
  "error": "initial_payment must be at least 20% of object_cost"
}
```

### GET /cache

Если кэш пуст:

```json
{
  "error": "empty cache"
}
```

Если есть записи:

```json
[
  {
    "id": 0,
    "params": { ... },
    "program": { ... },
    "aggregates": { ... }
  },
  ...
]
```

## Примеры запросов

```bash
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "object_cost": 5000000,
      "initial_payment": 1000000,
      "months": 240
    },
    "program": {
      "salary": true
    }
  }'

curl http://localhost:8080/cache
```


