# Ипотечный калькулятор

Сервис расчёта параметров ипотеки с локальным in-memory кэшем.

## Оглавление

- [Описание](#описание)  
- [Функциональность](#функциональность)  
- [Архитектура](#архитектура)  
- [Конфигурация](#конфигурация)
- [Тестовое покрытие](#тестовое-покрытие)   
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
- [Примеры запросов](#примеры-запросов)  
- [Технологии](#технологии)  

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


## Тестовое покрытие

В проекте покрытие unit-тестами составляет более 80% для всех ключевых частей. Ниже приведены результаты `go test -cover` по пакетам:

- **cmd** (точка входа)  
  - Покрытие: 0.0% (без тестов, так как `main.go` только инициализирует приложение)  
- **internal/domain**  
  - Покрытие: нет тестов (модели без логики)  
- **internal/adapters/primary/http-adapter**  
  - Пакет `middleware`: 100.0%  
  - Пакет `router`: 100.0%  
  - Пакет `controller`: 91.9%  
  - Пакет (прочие): 87.5%  
- **internal/adapters/secondary/repositories/credit-repository**  
  - Покрытие: 100.0%  
- **internal/application/credit-service**  
  - Покрытие: 93.8%  
- **internal/config**  
  - Покрытие: 100.0%  

Итоговое покрытие по проекту — более 90% для основных компонентов (контроллеры, middleware, роутер, репозиторий, сервис, конфиг).



