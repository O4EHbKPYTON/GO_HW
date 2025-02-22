** Разоварачивание нового проекта

Установка зависимостей

go get -u github.com/beego/beego/v2
go get -u github.com/beego/bee
go get -u github.com/lib/pq

bee api  название проекта //создаёт стркутуру проекта

Генерация роутинга
bee generate routers

запуск с генерацией swagger
bee run -gendoc=true -downdoc=true

Описание директорий
conf                    дирекктории с конфигурационном файлом
routers                 содержит структуру HTTP путей
controllers              контроллер, содержит код для автоматической генерации документации в формате swagger и описание точек входа
swagger                 директория  содержащая swagger
tests                   директория с тестами, пока не трогаем
models                  директория с моделями, описание структуры полей БД и обработка входящих запросов 