package omitnull

import (
	"bytes"
	"testing"
)

func TestPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("int8", func(t *testing.T) {
		val := From[int8](5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("int16", func(t *testing.T) {
		val := From[int16](5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("int32", func(t *testing.T) {
		val := From[int32](5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("int64", func(t *testing.T) {
		val := From[int64](5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("int", func(t *testing.T) {
		val := From(5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("uint8", func(t *testing.T) {
		val := From[uint8](5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("uint16", func(t *testing.T) {
		val := From[uint16](5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("uint32", func(t *testing.T) {
		val := From[uint32](5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("uint64", func(t *testing.T) {
		val := From[uint64](5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("uint", func(t *testing.T) {
		val := From(5)
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `5` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(int64) != 5 {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`6`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 6 {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`7`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 7 {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(int64(8)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != 8 {
			t.Error("value wrong:", v)
		}
	})
	t.Run("string", func(t *testing.T) {
		val := From("hello")
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `"hello"` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `hello` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(string) != "hello" {
			t.Error("output was wrong:", b)
		}

		if err := val.UnmarshalJSON([]byte(`"a"`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != "a" {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`b`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != "b" {
			t.Error("value wrong:", v)
		}
		if err := val.Scan("c"); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != "c" {
			t.Error("value wrong:", v)
		}
	})
	t.Run("[]byte", func(t *testing.T) {
		val := From([]byte("hello"))
		// The JSON package has the following behavior for bytes:
		// > Array and slice values encode as JSON arrays, except that []byte
		// encodes as a base64-encoded string, and a nil slice encodes as the
		// null JSON value.
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `"aGVsbG8="` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `hello` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if !bytes.Equal(b.([]byte), []byte("hello")) {
			t.Errorf("output was wrong: %#v", b)
		}
		if err := val.UnmarshalJSON([]byte(`"YQ=="`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if !bytes.Equal(v, []byte{'a'}) {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`b`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if !bytes.Equal(v, []byte{'b'}) {
			t.Error("value wrong:", v)
		}
		if err := val.Scan("c"); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if !bytes.Equal(v, []byte{'c'}) {
			t.Error("value wrong:", v)
		}
	})
	t.Run("bool", func(t *testing.T) {
		val := From(true)
		// The JSON package has the following behavior for bytes:
		// > Array and slice values encode as JSON arrays, except that []byte
		// encodes as a base64-encoded string, and a nil slice encodes as the
		// null JSON value.
		if b, err := val.MarshalJSON(); err != nil {
			t.Error(err)
		} else if string(b) != `true` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.MarshalText(); err != nil {
			t.Error(err)
		} else if string(b) != `true` {
			t.Error("output was wrong:", string(b))
		}
		if b, err := val.Value(); err != nil {
			t.Error(err)
		} else if b.(bool) != true {
			t.Error("output was wrong", b)
		}
		if err := val.UnmarshalJSON([]byte(`true`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != true {
			t.Error("value wrong:", v)
		}
		if err := val.UnmarshalText([]byte(`false`)); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != false {
			t.Error("value wrong:", v)
		}
		if err := val.Scan(true); err != nil {
			t.Error(err)
		} else if v, ok := val.Get(); !ok {
			t.Error("should not be null!")
		} else if v != true {
			t.Error("value wrong:", v)
		}
	})
}
