# cyberpb2

Protobuf-контракты и утилиты для сервисов Cybermetrica.

## Установка

```bash
go get github.com/ansanych/cyberpb2
```

## Структура модуля

```
proto/
├── cyber.proto           # Базовые protobuf-контракты (Status, HealthReply, Empty и др.)
├── cyber.pb.go           # Сгенерированный Go-код из cyber.proto
├── cybermetrica.proto    # Protobuf-контракты сервиса Cybermetrica
├── cybermetrica.pb.go    # Сгенерированный Go-код из cybermetrica.proto
├── cyberfuel.proto       # Protobuf-контракты сервиса Cyberfuel
├── cyberfuel.pb.go       # Сгенерированный Go-код из cyberfuel.proto
├── status_json.go        # Кастомная JSON-сериализация для enum Status
└── protojson.go          # Утилиты для сериализации protobuf-сообщений в JSON
```

## Использование

### Импорт

```go
import (
    "github.com/ansanych/cyberpb2/proto"
)
```

### Сериализация protobuf-ответов в JSON

Enum `Status` по умолчанию сериализуется как **integer** (0, 1, 2...).  
Для получения строковых значений ("OK", "ERROR"...), используйте `proto.MarshalProtoJSON`:

```go
// Создаём ответ
reply := &proto.HealthReply{
    Postgres: proto.Status_OK,
    Mongo:    proto.Status_ACTIVE,
    Parser:   proto.Status_ERROR,
    Disk: &proto.Disk{
        All:  500.0,
        Used: 120.5,
        Free: 379.5,
    },
}

// Сериализация в JSON с enum-строками
jsonBytes, err := proto.MarshalProtoJSON(reply)
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(jsonBytes))
// {"postgres":"OK","mongo":"ACTIVE","parser":"ERROR","disk":{"all":500,"used":120.5,"free":379.5}}
```

### Стандартный json.Marshal для отдельных enum-значений

Для отдельных значений `Status` можно использовать стандартный `json.Marshal` — он будет выводить строку благодаря кастомному `MarshalJSON`:

```go
status := proto.Status_OK
jsonBytes, _ := json.Marshal(status)
fmt.Println(string(jsonBytes)) // "OK"
```

### Десериализация JSON в protobuf

```go
data := []byte(`{"postgres":"OK","mongo":"ACTIVE","parser":"ERROR"}`)
reply := &proto.HealthReply{}
err := proto.UnmarshalProtoJSON(data, reply)
// Принимает как строки ("OK"), так и числа (3)
```

### Использование с HTTP-сервером

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    // Получаем ответ от gRPC-сервиса
    reply, err := grpcClient.Health(ctx, &proto.Empty{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Сериализуем в JSON с enum-строками
    w.Header().Set("Content-Type", "application/json")
    jsonBytes, _ := proto.MarshalProtoJSON(reply)
    w.Write(jsonBytes)
}
```

## Генерация Go-кода из proto-файлов

```bash
make gen-go
```

## Генерация Python-кода

```bash
make gen-py