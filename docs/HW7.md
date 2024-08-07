# Домашнее задание 7 "Работа с gRPC"
## Основное задание
Перевести проект из домашнего задания 6 на работу с gRPC.
Для этого:
- Создать proto контракт сервиса
- В проекте нужно добавить в Makefile команды для генерации .go файлов из proto файлов и установки нужных зависимостей (использовать protoc)
- Сгенерировать gRPC сервис
- Покрыть gRPC-хендлеры тестами

## Дополнительное задание
- Добавить HTTP-gateway и валидацию protobuf сообщений
- Добавить swagger-ui и возможность совершать запросы из сваггера к сервису (поднять swagger-ui сервер)

### Дедлайны сдачи и проверки задания:
- 13 июля 23:59 (сдача) / 16 июля, 23:59 (проверка)


## Настройка переменных окружения

Для корректной работы приложения необходимо настроить следующие переменные окружения:

- `DATABASE_URL`: URL подключения к базе данных. Пример: `postgres://user:password@localhost:5432/database_name?sslmode=disable`
- `KAFKA_BROKERS`: Список адресов брокеров Kafka, разделенных запятыми. Пример: `localhost:9092,localhost:9093`
- `KAFKA_TOPIC`: Название топика Kafka, в который приложение будет отправлять сообщения. Пример: `my_topic`
- `OUTPUT_MODE`: Режим вывода информации (`stdout` для вывода в стандартный поток вывода или `kafka` для отправки сообщений в Kafka). Пример: `stdout`
- `GRPC_PORT`: Порт, на котором будет запущен gRPC сервер. Пример: `50051`

Эти переменные можно задать при использовании файла `.env` в корне проекта для их определения. Также возможно указание переменных окружения в `docker-compose.yml` для запуска в контейнере.