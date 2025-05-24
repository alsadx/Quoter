# Quoter

REST API сервис для хранения и управления цитатами.

## Функциональность

- Добавление новых цитат
- Получение всех цитат
- Получение случайной цитаты
- Фильтрация цитат по автору
- Удаление цитаты по ID

## Как запустить

1. Клонируй репозиторий
```bash
git clone https://github.com/alsadx/Quoter.git 
cd Quoter
```

2. Установи зависимости
```bash
go mod download
```

3. Запусти сервер
```bash
go run cmd/quoter/main.go
```

Сервис будет доступен по адресу: http://localhost:8080

## Примеры запросов через curl

#### Добавление цитаты
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius","quote":"Life is simple..."}'
  ```

💡 Совет: если в JSON есть апостроф (') или пробел — используй двойные кавычки или оберни всё в одинарные:
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Jimmy Carr", "quote":"Everyone is jealous of what you'"'"'ve got, no one is jealous of how you got it."}'
```

#### Получение всех цитат
```bash
curl http://localhost:8080/quotes
```

#### Получение случайной цитаты
```bash
curl http://localhost:8080/quotes/random
```

#### Фильтрация по автору
```bash
curl 'http://localhost:8080/quotes?author=Confucius'
```

💡 Совет: если в имени автора есть пробел , используй %20 или +
```bash
curl 'http://localhost:8080/quotes?author=Oscar%20Wilde'
```

#### Удаление цитаты по ID
curl -X DELETE http://localhost:8080/quotes/1

## Unit-тесты
Проект содержит unit-тесты для:
- Хранилища (storage/tests)
- Хэндлеров (handlers/handlers_test.go)

## Как запустить тесты
```bash
go test ./...
```

## Docker

#### Запуск через docker-compose

Из директории /docker:
```bash
docker-compose up -d
```

Из корня проекта:
```bash
docker-compose -f docker/docker-compose.yml up -d
```

#### Остановить контейнер

Из директории /docker:
```bash
docker-compose down
```

Из корня проекта:
```bash
docker-compose -f docker/docker-compose.yml down
```
