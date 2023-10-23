# User API

Тестовое задание по написанию сервиса работы с пользователями

## Старт сервиса

1. Клонируем репозиторий проекта:
```
git clone https://github.com/LinqCod/user-service.git
```
2. Далее необходимо сбилдить и запустить сервис. Сам сервис, как и бд (в данном случае postgres) работают в докер контейнерах, их конфигурация и связь прописаны в `docker-compose` файле. Также перед запуском следует установить пакет `golang-migrate`, если еще не был установлен.

### Команды

Для удобства все команды прописаны в `Makefile`

 - Для запуска сервиса достаточно использовать команду `make service-up`, которая сбилдит и запустит приложение, доступное на порту 8080:

```
make service-up
```

 - Для последующих запусков сервиса достаточно использовать команду `make service-run`, которая будет запускать уже существующий контейнер.

```
make service-run
```

 - Для завершения работы сервиса следует воспользоваться командой:

```
make service-down
```

 - Для создания файлов миграции таблицы пользователей необходимо выполнить следующую команду (для удобства файлы миграции, которые создают и дропают таблицу пользователей, уже присутствуют в проекте):

```
make migration-create-user
```

 - Чтобы накатить миграции достаточно выполнить команду:

```
make migration-up
``` 

 - Для роллбека миграций воспользуйтесь командой:

```
make migration-down
``` 

 - Для удобства добавлена команда для входа в postgres терминал созданной базы данных:

```
make shell-postgres
```
 
## О реализации

 - Сервис реализован, используя основные концепции `Clean Architecture` (разделение на слои handlers, repository и usecase, также прописаны доменные модели для каждого кейса)
 - Для создания основного контейнера сервиса используется `multi-stage build` дабы уменьшить вес итогового контейнера и увеличить скорость его запуска
 - Конфигурация вынесена в `.env` файл. Для инициализации и взаимодействия с конфигом используется пакет `viper`
 - Для логирования используется качественная библиотека `zap`
 - Реализован механизм `graceful shutdown` для грамотного завершения сервиса

## Эндпоинты

 - POST /api/v1/users - создание пользователя. На вход принимает json-объект с обязательными полями name, surname. Остальные поля (patronymic, age, gender, nationality) являются необязательными. Поля age, gender и nationality запрашиваются со сторонних апи. 
 - DELETE /api/v1/users/:id - удаление пользователя по айди. На вход принимает айди пользователя в качестве path параметра.
 - PATCH /api/v1/users/:id - изменение пользователя по айди. На вход принимает айди пользователя в качестве path параметра, а также json-объект с новыми значениями полей. Если какие то поля не будут указаны в request body, в базе они останутся прежними.
 - GET /api/v1/users/:id - получение пользователя по айди. На вход принимает айди пользователя в качестве path параметра.
 - GET /api/v1/users - получение отфильтрованного списка пользователей. Фильтрация происходит путем указания следующих query параметров: gender, minAge, maxAge, nationality, count (пагинация). Параметры можно указать как все, так и часть, или вообще не указывать. В случае, если какие то параметры не указаны, они не будут учитываться при фильтрации. 