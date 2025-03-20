FROM golang:1.23-alpine AS build
WORKDIR /build

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init

RUN apk add --no-cache make
RUN go build -o main main.go 

FROM alpine AS runnder 

WORKDIR /app

RUN apk add --no-cache make
RUN go build -o main main.go  

FROM alpine AS runner

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=build /build/main ./main
COPY --from=build /build/docs ./docs
COPY --from=build /build/swagger.yaml ./swagger.yaml

CMD ["/app/main"]
