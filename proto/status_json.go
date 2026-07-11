// Package proto — кастомная JSON-сериализация для enum Status
//
// Проблема: protobuf генерирует тип Status как int32, и при сериализации
// через стандартный encoding/json enum-значения выводятся как числа (0, 1, 2...).
// Это неудобно для HTTP API, где ожидаются строковые значения ("OK", "ERROR"...).
//
// Решение: добавляем методы MarshalJSON/UnmarshalJSON для типа Status,
// чтобы при JSON-сериализации использовались строковые имена enum.
package proto

import (
	"fmt"
	"strings"
)

// MarshalJSON реализует json.Marshaler для Status.
// Вместо числа выводит строковое имя enum, например "OK" вместо 3.
func (x Status) MarshalJSON() ([]byte, error) {
	name, ok := Status_name[int32(x)]
	if !ok {
		return nil, fmt.Errorf("unknown Status value: %d", x)
	}
	return []byte(`"` + name + `"`), nil
}

// UnmarshalJSON реализует json.Unmarshaler для Status.
// Принимает как строковое имя ("OK"), так и число (3).
func (x *Status) UnmarshalJSON(data []byte) error {
	// Пробуем распарсить как строку
	s := strings.Trim(string(data), `"`)
	if s != "" {
		val, ok := Status_value[s]
		if ok {
			*x = Status(val)
			return nil
		}
	}

	// Пробуем распарсить как число
	var n int32
	_, err := fmt.Sscanf(string(data), "%d", &n)
	if err != nil {
		return fmt.Errorf("cannot unmarshal Status from %s", string(data))
	}
	*x = Status(n)
	return nil
}
