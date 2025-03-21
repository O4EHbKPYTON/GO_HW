FROM golang:1.23-alpine AS runner

WORKDIR /app

# Устанавливаем нужные пакеты
RUN apk add --no-cache curl git make

# Устанавливаем Swag и BeeGo CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/beego/bee/v2@latest

# Копируем код в контейнер
COPY . .

# Загружаем зависимости
RUN go mod download

# Генерируем Swagger-документацию
RUN swag init

# Запускаем Beego с генерацией документации
CMD ["bee", "run", "-gendoc=true", "-downdoc=true"]
