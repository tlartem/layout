# GRPC

## Установка protoc

Если у вас Linux, то вам ничего скачивать не нужно, все зависимости уже есть в проекте в `bin`.

Если у вас другая ОС, то нужно сделать следующее:

1. Скачать последнюю версию protoc для вашей ОС из https://github.com/protocolbuffers/protobuf/releases (osx — MacOS, win — Windows).
2. Скопировать (заменить) protoc в папке `bin`.
3. Рядом положить папку из того же архива `include`, если ее нет.
4. Дополнительно установить зависимости:

```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

5. После чего из папки `bin` GOPATH скопировать (заменить) файлы protoc-gen-go и protoc-gen-go-grpc в папке `bin` проекта.

## Генерация protobuf файлов и grpc сервера

```bash
  make grpc_gen
```

# Redis
Продовый Redis находится по адресу:
89.111.142.19:6379

UI продового Redis:
http://89.111.142.19:8081

Чтобы на проде редис работал надо добавить переменную окружения в раздел [Variables]:
```bash
REDIS_ADDR=89.111.142.19:6379
```

# Kafka
Переменные для Gitlab
```shell
# Kafka
KAFKA_WRITER_ADDR=89.111.142.19:9092,89.111.142.19:9093,89.111.142.19:9094
KAFKA_CONSUMER_ADDR=89.111.142.19:9092,89.111.142.19:9093,89.111.142.19:9094
```
