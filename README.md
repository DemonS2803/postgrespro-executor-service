# Тестовое задание в PostgresPro
- стек: Golang (chi), postgresql, redis, docker (compose)

## Запуск
- git clone https://github.com/DemonS2803/postgrespro-executor-service
- docker-compose up - в корне проекта

Тесты
- go test - в корне проекта
- 

## Решения
- для хранения текущих процессов используется redis (задел на расширение системы)
- так как скрипт, запущенный в разные отрезки времени может выдать разный результат (пинг при разрыве кабеля не будет проходить), принято решение вынести результат работы скриптов в отдельную таблицу
- не было обговорено, когда и как возвращать результат работы скрипта, поэтому решение такое: пользователь запускает скрипт и получает (условно) id процесса. Далее он может отправить запрос на получение результата
- добавил базовую аутентификацию (base_token и super_token). Но нигде ее не использовал
- изначально хотелось поставить тайм-аут выполнения команд, но существует условие на долгие команды, поэтому они могут крутиться в системе вечно :-)
- результат работы "долгой" команды хранится в редисе, и периодически записывается в БД (уменьшаем нагрузку)


## Endpoints

- [POST] localhost:8080/command - создание команды
```
curl --location 'http://localhost:8080/command' \
    --header 'token: admin_token' \
    --header 'Content-Type: application/json' \
    --data '{
    "code": "echo \"new command\"",
    "description": "maybe"
    }'
```

- [GET] localhost:8080/command/exec/{commandId} - запуск команды
```
curl --location --request GET 'http://localhost:8080/command/exec/1' \
    --header 'token: admin_token' \
    --header 'Content-Type: application/json' \
    --data '{
    "code": "echo \"new command\"",
    "description": "maybe"
    }'
```       

- [GET] localhost:8080/command/result/{completedCommandId} - получение результата (completedCommandId - id, приходящий в теле запроса, запускающего скрипт)
```
curl --location --request GET 'http://localhost:8080/command/result/1' \
    --header 'token: admin_token' \
    --header 'Content-Type: application/json' \
    --data '{
    "code": "echo \"new command\"",
    "description": "maybe"
    }'
```

- [GET] localhost:8080/command/stop/{completedCommandId} - остановка долгой команды
```
curl --location --request GET 'http://localhost:8080/command/stop/1' \
    --header 'token: admin_token' \
    --header 'Content-Type: application/json' \
    --data '{
    "code": "echo \"new command\"",
    "description": "maybe"
    }'
```
- [GET] localhost:8080/command/{id} - получение данных о команде
```
curl --location --request GET 'http://localhost:8080/command/1' \
    --header 'token: admin_token' \
    --header 'Content-Type: application/json' \
    --data '{
        "code": "echo \"new command\"",
        "description": "maybe"
    }'
```
- [GET] localhost:8080/commands - получение списка команд
```
curl --location --request GET 'http://localhost:8080/commands?limit=3&offset=0' \
    --header 'token: admin_token' \
    --header 'Content-Type: application/json' \
    --data '{
        "code": "echo \"new command\"",
        "description": "maybe"
    }'
```