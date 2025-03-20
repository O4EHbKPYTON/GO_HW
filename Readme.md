# Разворачивание проекта

## Установка зависимостей

```
go get -u github.com/beego/beego/v2
go get -u github.com/beego/bee
go get -u github.com/lib/pq
go get -u github.com/golang-jwt/jwt/v5
```

### Cоздание нового проекта

```
bee api название проекта
```

### Генерация роутинга

```
bee generate routers
```

### запуск с генерацией swagger

```
bee run -gendoc=true -downdoc=true
```

url swagger после запуска
http://localhost:8080/swagger/

### Описание диррикторий

* conf 				дирриктория с конфигурационным файлом
* routers				содержит структуру HTTP путей
* controllers			контроллер, содержит код для автоматической генерации документации в формате swagger и описание точек входа
* swagger				дерриктория содержащая swagger
* tests               дирриктория с тестами, пока не трогаем
* models              дирриктория с моделами, описание структуры полей БД и обработка входящих запросов


### Подключение сессий
```
SessionOn = true
sessionprovider = "postgresql"
sessionproviderconfig = "указать путь до базы"
```
Создать таблицу
```
CREATE TABLE session ( session_key char(64) NOT NULL, session_data bytea, session_expiry timestamp NOT NULL, CONSTRAINT session_key PRIMARY KEY(session_key) );
```