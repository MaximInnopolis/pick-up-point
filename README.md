# Домашнее задание №8 "Дорога к реальным сервисам"


## Цель

Оптимизация нагрузки на сервис, наблюдение и контроль работы сервиса

## Основное задание

- Реализуйте простое кеширование (In-memory Cache) для уменьшения нагрузки на базу данных вашего приложения
- Внедрите механизм инвалидации кеша, объясните свой выбор в документации пакета кеш клиента (godoc/md)
- Необходимо реализовать сбор метрик с вашего приложения. Добавьте кастомную метрику по количеству выданных заказов

## Дополнительное задание

- Добавьте ограничение потребления ресурсов вашего кеша (LRU/LFU или их модификации). Обновите документацию пакета кеширования
- Добавьте трейсинг запросов внутри приложения
- Выделите пакет кеширования в потенциально универсальный, с Generic типами ключей и значений, и конфигурацией снаружи (стратегия кеширования для экономии ресурсов, TTL кешированных значений и др.)

### Дедлайны сдачи и проверки задания:
- 20 июля 23:59 (сдача) / 23 июля, 23:59 (проверка)


## Инвалидация кэша

Инвалидация кэша выполнена путем периодической инвалидации просроченных записей, чей TTL истек, таким образом устаревшие
данные удаляются из кэша. Это реализовано методом```InvalidateExpired```, который вызывается в горутине с периодичностью
раз в минуту


## Документация по домашним заданиям

- [Домашнее задание 1](docs/HW1.md)
- [Домашнее задание 2](docs/HW2.md)
- [Домашнее задание 3](docs/HW3.md)
- [Домашнее задание 4](docs/HW4.md)
- [Домашнее задание 5](docs/HW5.md)
- [Домашнее задание 6](docs/HW6.md)
- [Домашнее задание 7](docs/HW7.md)