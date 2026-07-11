// Package proto — утилиты для сериализации protobuf-сообщений в JSON
//
// Предоставляет функции для конвертации protobuf-ответов в JSON
// с корректным отображением enum Status в виде строк.
package proto

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// MarshalOptions — опции protojson для сериализации с enum-строками.
// Используется для формирования HTTP JSON-ответов из protobuf-сообщений.
var MarshalOptions = protojson.MarshalOptions{
	// UseEnumNumbers: false — выводить enum как строки ("OK", "ERROR"...)
	UseEnumNumbers: false,

	// EmitUnpopulated: true — включать поля с нулевыми значениями
	EmitUnpopulated: true,

	// UseProtoNames: false — использовать camelCase имена (как в json tag)
	UseProtoNames: false,
}

// MarshalProtoJSON сериализует protobuf-сообщение в JSON-строку.
// Enum-поля выводятся как строки, а не числа.
func MarshalProtoJSON(msg proto.Message) ([]byte, error) {
	return MarshalOptions.Marshal(msg)
}

// MarshalProtoJSONString сериализует protobuf-сообщение в JSON-строку (string).
func MarshalProtoJSONString(msg proto.Message) (string, error) {
	data, err := MarshalOptions.Marshal(msg)
	return string(data), err
}

// UnmarshalOptions — опции protojson для десериализации.
// Принимает enum как в виде строк, так и чисел.
var UnmarshalOptions = protojson.UnmarshalOptions{
	// DiscardUnknown: true — игнорировать неизвестные поля
	DiscardUnknown: true,
}

// UnmarshalProtoJSON десериализует JSON в protobuf-сообщение.
// Принимает enum как в виде строк ("OK"), так и чисел (3).
func UnmarshalProtoJSON(data []byte, msg proto.Message) error {
	return UnmarshalOptions.Unmarshal(data, msg)
}
