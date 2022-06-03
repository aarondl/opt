// Package null exposes a Val(ue) type that wraps a regular value with the
// ability to be 'null'.
package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/aarondl/opt"
	"github.com/aarondl/opt/internal/globaldata"
)

// State is the state of the nullable object
type State int

const (
	StateNull State = 0
	StateSet  State = 1
)

// String -er interface implementation
func (s State) String() string {
	switch s {
	case StateNull:
		return "null"
	case StateSet:
		return "set"
	default:
		panic("unknown")
	}
}

// Val allows representing a value with a state of "null" or "set".
type Val[T any] struct {
	value T
	state State
}

// New creates a new value that is null, shorthand for Val[T]{}
func New[T any]() Val[T] {
	return Val[T]{}
}

// From a value which is considered 'set'
func From[T any](val T) Val[T] {
	return Val[T]{
		value: val,
		state: StateSet,
	}
}

// FromPtr creates a value from a pointer, if the pointer is null it will be
// 'null', if it has a value the deferenced value is stored.
func FromPtr[T any](val *T) Val[T] {
	if val == nil {
		return Val[T]{state: StateNull}
	}
	return Val[T]{
		value: *val,
		state: StateSet,
	}
}

// Get the underlying value, if one exists.
func (v Val[T]) Get() (T, bool) {
	if v.state == StateSet {
		return v.value, true
	}

	var empty T
	return empty, false
}

// GetOr gets the value or returns a fallback if the value does not exist.
func (v Val[T]) GetOr(fallback T) T {
	if v.state == StateSet {
		return v.value
	}
	return fallback
}

// MustGet retrieves the value or panics if it's null
func (v Val[T]) MustGet() T {
	val, ok := v.Get()
	if !ok {
		panic("no value present")
	}

	return val
}

// Map transforms the value inside if it is set, else it returns a value of the
// same state.
//
// Until a later Go version adds type parameters to methods, it is not possible
// to map to a different type. See the non-method function Map if you need
// another type.
func (v Val[T]) Map(fn func(T) T) Val[T] {
	if v.state == StateSet {
		return From(fn(v.value))
	}
	return Val[T]{state: v.state}
}

// Map transforms the value inside if it is set, else it returns a value
// of the same state.
func Map[A any, B any](v Val[A], fn func(A) B) Val[B] {
	if v.state == StateSet {
		return From(fn(v.value))
	}
	return Val[B]{state: v.state}
}

// Set the value (and the state to 'set')
func (v *Val[T]) Set(val T) {
	v.value = val
	v.state = StateSet
}

// Null sets the value to null (state is set to 'null')
func (v *Val[T]) Null() {
	var empty T
	v.value = empty
	v.state = StateNull
}

// SetPtr sets the value to (value, set) if val is non-nil or (??, null) if not.
// The value is dereferenced before stored.
func (v *Val[T]) SetPtr(val *T) {
	if val == nil {
		v.state = StateNull
		return
	}
	v.value = *val
	v.state = StateSet
}

// IsSet returns true if v contains a non-null value
func (v Val[T]) IsSet() bool {
	return v.state == StateSet
}

// IsNull returns true if v contains a null value
func (v Val[T]) IsNull() bool {
	return v.state == StateNull
}

// State retrieves the internal state, mostly useful for testing.
func (v Val[T]) State() State {
	return v.state
}

// Ptr returns a pointer to the value, or nil if null.
func (v Val[T]) Ptr() *T {
	if v.state == StateSet {
		return &v.value
	}
	return nil
}

// UnmarshalJSON implements json.Unmarshaler
func (v *Val[T]) UnmarshalJSON(data []byte) error {
	switch {
	case len(data) == 0:
		return errors.New("cannot unmarshal empty bytes into null value")
	case bytes.Equal(data, globaldata.JSONNull):
		v.state = StateNull
		return nil
	default:
		err := json.Unmarshal(data, &v.value)
		if err != nil {
			return err
		}
		v.state = StateSet
		return nil
	}
}

// MarshalJSON implements json.Marshaler.
func (v Val[T]) MarshalJSON() ([]byte, error) {
	switch v.state {
	case StateSet:
		return json.Marshal(v.value)
	default:
		return globaldata.JSONNull, nil
	}
}

// MarshalText implements encoding.TextMarshaler.
func (v Val[T]) MarshalText() ([]byte, error) {
	if v.state != StateSet {
		return nil, nil
	}

	var text string
	if err := opt.ConvertAssign(&text, v.value); err != nil {
		return nil, err
	}
	return []byte(text), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (v *Val[T]) UnmarshalText(text []byte) error {
	if err := opt.ConvertAssign(&v.value, string(text)); err != nil {
		return err
	}

	v.state = StateSet
	return nil
}

// Scan implements the sql.Scanner interface. If the wrapped type implements
// sql.Scanner then it will call that.
func (v *Val[T]) Scan(value any) error {
	if value == nil {
		v.state = StateNull
		return nil
	}
	v.state = StateSet
	return opt.ConvertAssign(&v.value, value)
}

// Value implements the driver.Valuer interface. If the underlying type
// implements the driver.Valuer it will call that (when not null).
// Go primitive types will be converted where possible.
//
//   int64
//   float64
//   bool
//   []byte
//   string
//   time.Time
func (v Val[T]) Value() (driver.Value, error) {
	if v.state != StateSet {
		return nil, nil
	}

	refVal := reflect.ValueOf(v.value)
	if refVal.Type().Implements(globaldata.DriverValuerIntf) {
		valuer := refVal.Interface().(driver.Valuer)
		return valuer.Value()
	}

	switch refVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return refVal.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		newVal := refVal.Convert(globaldata.Int64Type)
		return newVal.Int(), nil
	case reflect.Bool:
		return v.value, nil
	case reflect.String:
		return v.value, nil
	case reflect.Slice:
		itemKind := refVal.Type().Elem().Kind()
		if itemKind == reflect.Uint8 {
			return v.value, nil
		}
		return nil, errors.New("underlying slice type has no implementation for driver.Valuer")
	case reflect.Struct:
		if refVal.Type() == globaldata.TimeTimeType {
			return v.value, nil
		}
		return nil, errors.New("underlying struct type has no implementation for driver.Valuer")
	default:
		return nil, errors.New("underlying type has no implementation for driver.Valuer")
	}
}
