package domain

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/samber/oops"
)

type Object[T any] struct {
	data T
}

func NewObject[T any](data T) Object[T] {
	return Object[T]{data: data}
}

func (op *Object[T]) Data() T {
	return op.data
}

func (op Object[T]) Value() (driver.Value, error) {
	bytes, err := json.Marshal(op.data)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to marshal Object data")
	}
	return bytes, nil
}

func (op *Object[T]) Scan(value any) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return oops.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	if err := json.Unmarshal(bytes, &op.data); err != nil {
		return oops.Wrapf(err, "failed to unmarshal Object data")
	}
	return nil
}

func (op Object[T]) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(op.data)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to marshal Object data")
	}
	return bytes, nil
}

func (op *Object[T]) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &op.data); err != nil {
		return oops.Wrapf(err, "failed to unmarshal Object data")
	}
	return nil
}
