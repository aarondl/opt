// Package null exposes a Val(ue) type that wraps a regular value with the
// ability to be 'omitted' or 'unset'.
package omit

import (
	"bytes"
	"database/sql/driver"
	"encoding"
	"errors"
	"reflect"

	"github.com/aarondl/opt"
	"github.com/aarondl/opt/internal/globaldata"
)

// state is the state of the omittable object
type state int

const (
	StateUnset state = 0
	StateSet   state = 1
)

// String -er interface implementation
func (s state) String() string {
	switch s {
	case StateUnset:
		return "unset"
	case StateSet:
		return "set"
	default:
		panic("unknown")
	}
}

// Val allows representing a value with a state of "unset" or "set".
// Its zero value is usfel and initially "unset".
type Val[T any] struct {
	value T
	state state
}

// From a value which is considered 'set'
func From[T any](val T) Val[T] {
	return Val[T]{
		value: val,
		state: StateSet,
	}
}

// FromPtr creates a value from a pointer, if the pointer is null it will be
// 'unset', if it has a value the deferenced value is stored.
func FromPtr[T any](val *T) Val[T] {
	if val == nil {
		return Val[T]{state: StateUnset}
	}
	return Val[T]{
		value: *val,
		state: StateSet,
	}
}

// FromCond conditionally creates a 'set' value if the bool is true, else
// it will return an omitted value.
func FromCond[T any](val T, ok bool) Val[T] {
	if !ok {
		return Val[T]{}
	}
	return Val[T]{
		value: val,
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

// GetOrZero returns the zero value for T if the value was omitted.
func (v Val[T]) GetOrZero() T {
	if v.state != StateSet {
		var t T
		return t
	}
	return v.value
}

// MustGet retrieves the value or panics if it's null
func (v Val[T]) MustGet() T {
	val, ok := v.Get()
	if !ok {
		panic("no value present")
	}

	return val
}

// Or returns v or other depending on their states. In general
// set > unset and therefore the one with the state highest in that
// area will win out.
//
//	v     | other | result
//	------------- | -------
//	set   | _     | v
//	unset | set   | other
//	unset | unset | v
func (v Val[T]) Or(other Val[T]) Val[T] {
	switch {
	case v.state == StateUnset && other.state == StateSet:
		return other
	default:
		return v
	}
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

// Map transforms the value inside if it is set, else it returns value of the
// same state.
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

// Unset the value (state is set to 'unset')
func (v *Val[T]) Unset() {
	var empty T
	v.value = empty
	v.state = StateUnset
}

// IsValue returns true if v contains a value (ie. not omitted/unset)
func (v Val[T]) IsValue() bool {
	return v.state == StateSet
}

// IsUnset returns true if v contains no value
func (v Val[T]) IsUnset() bool {
	return v.state == StateUnset
}

// State retrieves the internal state, mostly useful for testing.
func (v Val[T]) State() state {
	return v.state
}

// UnmarshalJSON implements json.Unmarshaler. Notably will fail to unmarshal
// if given a null.
func (v *Val[T]) UnmarshalJSON(data []byte) error {
	switch {
	case len(data) == 0:
		var zero T
		v.value = zero
		v.state = StateUnset
		return nil
	case bytes.Equal(data, globaldata.JSONNull):
		return errors.New("cannot unmarshal 'null' value into omit value")
	default:
		err := opt.JSONUnmarshal(data, &v.value)
		if err != nil {
			return err
		}
		v.state = StateSet
		return nil
	}
}

// MarshalJSON implements json.Marshaler.
//
// Note that this type cannot possibly work with the stdlib json package due
// to there being no way for the json package to omit a value based on its
// internals.
//
// That's to say even if you have an `omitempty` tag with this type, it will
// still show up in outputs as {"val": null} because this functionality is not
// supported.
//
// For a package that works well with this package see github.com/aarondl/json.
func (v Val[T]) MarshalJSON() ([]byte, error) {
	switch v.state {
	case StateSet:

		return opt.JSONMarshal(v.value)
	default:
		return globaldata.JSONNull, nil
	}
}

// MarshalJSONIsZero returns true if this value should be omitted by the json
// marshaler.
//
// Deprecated: This method was necessary to support true omitting of values
// using a json fork, since then the std library has added support for this
// and so we no longer need this method.
func (v Val[T]) MarshalJSONIsZero() bool {
	return v.IsZero()
}

// IsZero returns true if this value should be omitted by the json
// marshaler.
//
// There is a special case in which we omit the value even if the value is `set`
// which is when the value is going to write out `nil` (pointers, maps
// and slices that are nil) when marshaled.
//
// The reason this is important is if we marshal(From[[]int](nil)) with the
// special json fork, it will emit `null` without this override. This is bad
// because this same package even with the json fork cannot consume a null.
//
// In order to achieve symmetry in encoding/decoding we'll quietly omit nil
// maps, slices, and ptrs as it was likely a mistake to try to .From(nil)
// for this type of value anyway.
func (v Val[T]) IsZero() bool {
	if v.state == StateUnset {
		return true
	}

	switch reflect.TypeOf(v.value).Kind() {
	case reflect.Map, reflect.Slice, reflect.Ptr:
		if reflect.ValueOf(v.value).IsNil() {
			return true
		}
	}

	return false
}

// MarshalText implements encoding.TextMarshaler.
func (v Val[T]) MarshalText() ([]byte, error) {
	if v.state != StateSet {
		return nil, nil
	}

	refVal := reflect.ValueOf(v.value)
	if refVal.Type().Implements(globaldata.EncodingTextMarshalerIntf) {
		valuer := refVal.Interface().(encoding.TextMarshaler)
		return valuer.MarshalText()
	}

	var text string
	if err := opt.ConvertAssign(&text, v.value); err != nil {
		return nil, err
	}
	return []byte(text), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (v *Val[T]) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		var zero T
		v.value = zero
		v.state = StateUnset
		return nil
	}

	refVal := reflect.ValueOf(&v.value)
	if refVal.Type().Implements(globaldata.EncodingTextUnmarshalerIntf) {
		valuer := refVal.Interface().(encoding.TextUnmarshaler)
		if err := valuer.UnmarshalText(text); err != nil {
			return err
		}
		v.state = StateSet
		return nil
	}

	if err := opt.ConvertAssign(&v.value, string(text)); err != nil {
		return err
	}

	v.state = StateSet
	return nil
}

// MarshalBinary tries to encode the value in binary. If it finds
// type that implements encoding.BinaryMarshaler it will use that,
// it will fallback to encoding.TextMarshaler if that is implemented,
// and failing that it will attempt to do some reflect to convert between
// the types to hit common cases like Go primitives.
func (v Val[T]) MarshalBinary() ([]byte, error) {
	if v.state != StateSet {
		return nil, nil
	}

	refVal := reflect.ValueOf(v.value)
	if refVal.Type().Implements(globaldata.EncodingBinaryMarshalerIntf) {
		valuer := refVal.Interface().(encoding.BinaryMarshaler)
		return valuer.MarshalBinary()
	}

	if refVal.Type().Implements(globaldata.EncodingTextMarshalerIntf) {
		valuer := refVal.Interface().(encoding.TextMarshaler)
		return valuer.MarshalText()
	}

	var buf []byte
	if err := opt.ConvertAssign(&buf, v.value); err != nil {
		return nil, err
	}

	return buf, nil
}

// UnmarshalBinary tries to reverse the value MarshalBinary operation.
// See documentation there for details about supported types.
func (v *Val[T]) UnmarshalBinary(b []byte) error {
	if len(b) == 0 {
		var zero T
		v.value = zero
		v.state = StateUnset
		return nil
	}

	refVal := reflect.ValueOf(&v.value)
	if refVal.Type().Implements(globaldata.EncodingBinaryUnmarshalerIntf) {
		valuer := refVal.Interface().(encoding.BinaryUnmarshaler)
		if err := valuer.UnmarshalBinary(b); err != nil {
			return err
		}
		v.state = StateSet
		return nil
	}

	if refVal.Type().Implements(globaldata.EncodingTextUnmarshalerIntf) {
		valuer := refVal.Interface().(encoding.TextUnmarshaler)
		if err := valuer.UnmarshalText(b); err != nil {
			return err
		}
		v.state = StateSet
		return nil
	}

	if err := opt.ConvertAssign(&v.value, b); err != nil {
		return err
	}

	v.state = StateSet
	return nil
}

// Scan implements the sql.Scanner interface. If the wrapped type implements
// sql.Scanner then it will call that.
func (v *Val[T]) Scan(value any) error {
	if value == nil {
		return errors.New("cannot store 'null' value in omit value")
	}
	v.state = StateSet
	return opt.ConvertAssign(&v.value, value)
}

// Value implements the driver.Valuer interface. If the underlying type
// implements the driver.Valuer it will call that (when not unset).
// Go primitive types will be converted where possible.
//
// Because sql doesn't have an analog to unset it will marshal as null in these
// cases.
//
//	int64
//	float64
//	bool
//	[]byte
//	string
//	time.Time
func (v Val[T]) Value() (driver.Value, error) {
	if v.state != StateSet {
		return nil, nil
	}

	return opt.ToDriverValue(v.value)
}

// Equal compares two nullable values and returns true if they are equal.
func Equal[T comparable](a, b Val[T]) bool {
	if a.state != b.state {
		return false
	}

	// states are equal, thus if set, they could have different values
	if a.state != StateSet {
		return true
	}

	return a.value == b.value
}
