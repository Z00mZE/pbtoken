# pbtoken
Шифрованный токен чтоб не потерять (как всегда). Пакет создавался для использования сериализации protobuf

#   Установка зависимостей

##  Установка protobuf
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Генерация

```shell
protoc --go_out=. --go-grpc_out=. proto/*.proto
```