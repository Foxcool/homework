# 3 homework

Сделать простейший RESTful CRUD по созданию, удалению, просмотру и обновлению пользователей.
Пример API - https://app.swaggerhub.com/apis/otus55/users/1.0.0

* Добавить базу данных с persistent volume для приложения.
* Docker-Image базы данных должен использоваться из официального докер-репозитория.
* Под с БД должен запускаться StatefulSet-ом с количеством реплик - 1.
* Конфигурация приложения должна храниться в Configmaps.
* Доступы к БД должны храниться в Secrets.
* Первоначальные миграции должны быть оформлены в качестве Pod-ы (или Job-ы).
* Ingress-ы должны также вести на url arch.homework/otusapp/* (как и в прошлом занятии)

На выходе должны быть предоставлена
1) ссылка на директорию в github, где находится директория с манифестами кубернетеса
2) команда kubectl apply -f, которая запускает в правильном порядке манифесты кубернетеса.
3) Postman коллекция, в которой будут представлены примеры запросов к сервису на создание, получение, изменение и удаление пользователя. Запросы из коллекции должны работать сразу после применения манифестов, без каких-то дополнительных подготовительных действий.

Задание со звездочкой (необязательное, но дает дополнительные баллы):
* +3 балла за шаблонизацию приложения в helm 3 чартах
* +2 балла за использование официального helm чарта для БД и подключение его в чарт приложения в качестве зависимости. 


## OpenAPI

### Generate API server code

```shell script
oapi-codegen --package api --generate types,server,spec,skip-prune ./docs/swagger.yaml > ./api/swagger.gen.go
```


## Docker
    sudo docker-compose up -d
    
## Helm

Run
```shell script
cd homework-chart
helm install mywork .
```

Install deps
```shell script
helm dependency update
```

## Postman collection

    tests/postman/collection.json