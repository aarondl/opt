package omitnull

import (
	"bytes"
	"database/sql/driver"
	"net"
	"testing"
	"time"

	"github.com/aarondl/opt"
	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
)

func TestConstruction(t *testing.T) {
	t.Parallel()

	hello := "hello"

	val := From("hello")

	checkState(t, val, StateSet)
	if !val.IsSet() {
		t.Error("should be set")
	}

	val = FromPtr(&hello)
	checkState(t, val, StateSet)
	val = FromPtr[string](nil)
	checkState(t, val, StateNull)
	if !val.IsNull() {
		t.Error("should be null")
	}

	val = Val[string]{}
	checkState(t, val, StateUnset)
	if !val.IsUnset() {
		t.Error("should be unset")
	}
}

func TestConversions(t *testing.T) {
	val := FromOmit(omit.Val[int]{})
	checkState(t, val, StateUnset)
	val = FromOmit(omit.From(5))
	checkState(t, val, StateSet)
	if val.MustGet() != 5 {
		t.Error("wrong value")
	}
	o := val.MustGetOmit()
	if !o.IsSet() {
		t.Error("should be set")
	}
	if o.MustGet() != 5 {
		t.Error("wrong value")
	}
	o, ok := val.GetOmit()
	if !ok || !o.IsSet() {
		t.Error("should be set")
	}
	val.Unset()
	o = val.MustGetOmit()
	if !o.IsUnset() {
		t.Error("should be unset")
	}
	o, ok = val.GetOmit()
	if !ok || !o.IsUnset() {
		t.Error("should be unset")
	}
	val.SetPtr(nil)
	_, ok = val.GetOmit()
	if ok {
		t.Error("should be nil")
	}
	val = FromNull(null.Val[int]{})
	checkState(t, val, StateNull)
	val = FromNull(null.From(5))
	checkState(t, val, StateSet)
	if val.MustGet() != 5 {
		t.Error("wrong value")
	}
	n := val.MustGetNull()
	if !n.IsSet() {
		t.Error("should be set")
	}
	if n.MustGet() != 5 {
		t.Error("wrong value")
	}
	n, ok = val.GetNull()
	if !ok || !n.IsSet() || n.MustGet() != 5 {
		t.Error("should be set")
	}
	val.Null()
	n = val.MustGetNull()
	if !n.IsNull() {
		t.Error("should be null")
	}
	n, ok = val.GetNull()
	if !ok || !n.IsNull() {
		t.Error("should be null")
	}
	val.Unset()
	o = val.MustGetOmit()
	if !o.IsUnset() {
		t.Error("should be unset")
	}
	o, ok = val.GetOmit()
	if !ok || !o.IsUnset() {
		t.Error("should be unset")
	}
	val.SetPtr(nil)
	_, ok = val.GetOmit()
	if ok {
		t.Error("should be nil")
	}
}

func TestConversionOmitFail(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("should have panic'd")
		}
	}()
	val := FromPtr[int](nil)
	_ = val.MustGetOmit() // nil is not an allowed value
}

func TestConversionNullFail(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("should have panic'd")
		}
	}()
	var val Val[int]
	_ = val.MustGetNull() // omitted is not an allowed value
}

func TestGet(t *testing.T) {
	t.Parallel()

	val := From("hello")
	if val.MustGet() != "hello" {
		t.Error("wrong value")
	}
	if val.GetOr("hi") != "hello" {
		t.Error("wrong value")
	}
	if val.GetOrZero() != "hello" {
		t.Error("wrong value")
	}
	if *val.MustPtr() != "hello" {
		t.Error("wrong value")
	}

	val.Null()
	if _, ok := val.Get(); ok {
		t.Error("should not be okay")
	}
	if val.GetOr("hi") != "hi" {
		t.Error("wrong value")
	}
	if val.GetOrZero() != "" {
		t.Error("wrong value")
	}
	if val.MustPtr() != nil {
		t.Error("should be nil")
	}

	val.Unset()
	if _, ok := val.Get(); ok {
		t.Error("should not be okay")
	}
	if val.GetOr("hi") != "hi" {
		t.Error("wrong value")
	}

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("should have panic'd")
			}
		}()
		_ = val.MustGet()
		t.Error("should not have made it here")
	}()

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("should have panic'd")
			}
		}()
		_ = val.MustPtr()
		t.Error("should not have made it here")
	}()
}

func TestOr(t *testing.T) {
	t.Parallel()

	var set, null, unset Val[int]
	set.Set(5)
	null.Null()
	unset.Unset()

	if set.Or(From(6)).MustGet() != 5 {
		t.Error("it should have returned 5")
	}
	if set.Or(null).MustGet() != 5 {
		t.Error("it should have returned 5")
	}
	if set.Or(unset).MustGet() != 5 {
		t.Error("it should have returned 5")
	}

	if null.Or(set).MustGet() != 5 {
		t.Error("it should have returned 5")
	}
	if !null.Or(null).IsNull() {
		t.Error("it should have returned null")
	}
	if !null.Or(unset).IsNull() {
		t.Error("it should have returned null")
	}

	if unset.Or(set).MustGet() != 5 {
		t.Error("it should have returned 5")
	}
	if !unset.Or(null).IsNull() {
		t.Error("it should have returned null")
	}
	if !unset.Or(unset).IsUnset() {
		t.Error("it should have returned unset")
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	val := Val[int]{}
	if !val.Map(func(int) int { return 0 }).IsUnset() {
		t.Error("it should still be unset")
	}
	if !Map(val, func(int) int { return 0 }).IsUnset() {
		t.Error("it should still be unset")
	}
	val.Set(5)
	if val.Map(func(i int) int { return i + 1 }).MustGet() != 6 {
		t.Error("wrong value")
	}
	if Map(val, func(i int) int { return i + 1 }).MustGet() != 6 {
		t.Error("wrong value")
	}
}

func TestChanges(t *testing.T) {
	t.Parallel()

	val := From("hello")
	checkState(t, val, StateSet)
	val.Null()
	checkState(t, val, StateNull)

	val = From("hello")
	checkState(t, val, StateSet)
	val.Unset()
	checkState(t, val, StateUnset)

	val = Val[string]{}
	checkState(t, val, StateUnset)
	val.Set("hello")
	checkState(t, val, StateSet)

	val = Val[string]{}
	checkState(t, val, StateUnset)
	val.SetPtr(nil)
	checkState(t, val, StateNull)

	hello := "hello"
	val = Val[string]{}
	checkState(t, val, StateUnset)
	val.SetPtr(&hello)
	checkState(t, val, StateSet)
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	val := From("hello")
	checkJSON(t, val, `"hello"`)
	val.Null()
	checkJSON(t, val, `null`)
	val.Unset()
	checkJSON(t, val, `null`)
}

func TestMarshalJSONIsZero(t *testing.T) {
	type testStruct struct {
		ID Val[int] `json:"id,omitzero"`
	}

	b, err := opt.JSONMarshal(testStruct{})
	if err != nil {
		t.Error(err)
	}
	if string(b) != `{}` {
		t.Errorf("expected empty json object, got: %s", b)
	}

	b, err = opt.JSONMarshal(testStruct{ID: From(0)})
	if err != nil {
		t.Error(err)
	}
	if string(b) != `{"id":0}` {
		t.Errorf("expected non-empty json object, got: %s", b)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	hello := Val[string]{}
	checkState(t, hello, StateUnset)

	if err := opt.JSONUnmarshal([]byte("null"), &hello); err != nil {
		t.Error(err)
	}
	checkState(t, hello, StateNull)

	if err := opt.JSONUnmarshal([]byte(`"hello"`), &hello); err != nil {
		t.Error(err)
	}
	checkState(t, hello, StateSet)

	if hello.MustGet() != "hello" {
		t.Error("expected hello")
	}

	hello.UnmarshalJSON(nil)
	checkState(t, hello, StateUnset)
}

func TestMarshalText(t *testing.T) {
	t.Parallel()

	hello := From("hello")
	b, err := hello.MarshalText()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "1hello" {
		t.Error("expected 1hello")
	}

	hello.Null()
	b, err = hello.MarshalText()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "0" {
		t.Error("expected a single 0")
	}

	marshaller := From(net.IPv4(1, 1, 1, 1))
	if b, err := marshaller.MarshalText(); err != nil {
		t.Error(err)
	} else if !bytes.Equal(b, []byte("11.1.1.1")) {
		t.Error("wrong value")
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Parallel()

	var val Val[string]
	if err := val.UnmarshalText([]byte("")); err != nil {
		t.Error(err)
	}
	checkState(t, val, StateUnset)

	if err := val.UnmarshalText([]byte("0")); err != nil {
		t.Error(err)
	}
	checkState(t, val, StateNull)

	if err := val.UnmarshalText([]byte("1hello")); err != nil {
		t.Error(err)
	}
	checkState(t, val, StateSet)
	if val.MustGet() != "hello" {
		t.Error("wrong value")
	}

	var unmarshaller Val[net.IP]
	if err := unmarshaller.UnmarshalText([]byte("")); err != nil {
		t.Error(err)
	}
	checkState(t, unmarshaller, StateUnset)

	if err := unmarshaller.UnmarshalText([]byte("11.1.1.1")); err != nil {
		t.Error(err)
	}
	checkState(t, unmarshaller, StateSet)
	if !unmarshaller.MustGet().Equal(net.IPv4(1, 1, 1, 1)) {
		t.Error("wrong value")
	}
}

func TestMarshalBinary(t *testing.T) {
	t.Parallel()

	hello := From("hello")
	b, err := hello.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "\x01hello" {
		t.Error("expected hello with a set byte")
	}

	hello.Null()
	b, err = hello.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(b, []byte{0}) {
		t.Error("expected a null byte")
	}

	hello.Unset()
	b, err = hello.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if len(b) != 0 {
		t.Error("expected empty byte slice")
	}

	marshaller := From(net.IPv4(1, 1, 1, 1))
	if b, err := marshaller.MarshalBinary(); err != nil {
		t.Error(err)
	} else if !bytes.Equal(b, []byte("\x011.1.1.1")) {
		t.Error("wrong value")
	}
}

func TestUnmarshalBinary(t *testing.T) {
	t.Parallel()

	var val Val[string]
	if err := val.UnmarshalBinary([]byte{}); err != nil {
		t.Error(err)
	}
	checkState(t, val, StateUnset)

	if err := val.UnmarshalBinary([]byte{0}); err != nil {
		t.Error(err)
	}
	checkState(t, val, StateNull)

	if err := val.UnmarshalBinary([]byte("\x01hello")); err != nil {
		t.Error(err)
	}
	checkState(t, val, StateSet)
	if val.MustGet() != "hello" {
		t.Error("wrong value")
	}

	var unmarshaller Val[net.IP]
	if err := unmarshaller.UnmarshalBinary([]byte{}); err != nil {
		t.Error(err)
	}
	checkState(t, unmarshaller, StateUnset)

	if err := unmarshaller.UnmarshalBinary([]byte("\x011.1.1.1")); err != nil {
		t.Error(err)
	}
	checkState(t, unmarshaller, StateSet)
	if !unmarshaller.MustGet().Equal(net.IPv4(1, 1, 1, 1)) {
		t.Error("wrong value")
	}
}

func TestScan(t *testing.T) {
	t.Parallel()

	var val Val[string]
	if err := val.Scan(nil); err != nil {
		t.Error(err)
	}
	checkState(t, val, StateNull)

	if err := val.Scan("hello"); err != nil {
		t.Error(err)
	}
	checkState(t, val, StateSet)
	if val.MustGet() != "hello" {
		t.Error("wrong value")
	}
}

type valuerImplementation struct{}

func (valuerImplementation) Value() (driver.Value, error) {
	return int64(1), nil
}

func TestValue(t *testing.T) {
	t.Parallel()

	var val Val[string]
	if v, err := val.Value(); err != nil {
		t.Error(err)
	} else if v != nil {
		t.Error("expected v to be nil")
	}

	val = From("hello")
	if v, err := val.Value(); err != nil {
		t.Error(err)
	} else if v.(string) != "hello" {
		t.Error("expected v to be nil")
	}

	date := time.Date(2000, 1, 1, 2, 30, 0, 0, time.UTC)
	nullTime := From(date)
	if v, err := nullTime.Value(); err != nil {
		t.Error(err)
	} else if !v.(time.Time).Equal(date) {
		t.Error("time was wrong")
	}

	valuer := From(valuerImplementation{})
	if v, err := valuer.Value(); err != nil {
		t.Error(err)
	} else if v.(int64) != 1 {
		t.Error("expect const int")
	}
}

func TestStateStringer(t *testing.T) {
	t.Parallel()

	if StateUnset.String() != "unset" {
		t.Error("bad value")
	}
	if StateNull.String() != "null" {
		t.Error("bad value")
	}
	if StateSet.String() != "set" {
		t.Error("bad value")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic")
		}
	}()
	_ = state(99).String()
}

func TestEqual(t *testing.T) {
	t.Parallel()

	a := Val[string]{}
	b := Val[string]{}
	if !Equal(a, b) {
		t.Error("should be equal")
	}

	a.Set("hello")
	if Equal(a, b) {
		t.Error("should not be equal")
	}

	b.Set("hello")
	if !Equal(a, b) {
		t.Error("should be equal")
	}

	b.Set("hi")
	if Equal(a, b) {
		t.Error("should not be equal")
	}

	a.Null()
	if Equal(a, b) {
		t.Error("should not be equal")
	}

	b.Null()
	if !Equal(a, b) {
		t.Error("should be equal")
	}

	a.Unset()
	if Equal(a, b) {
		t.Error("should not be equal")
	}
}

func checkState[T any](t *testing.T, val Val[T], want state) {
	t.Helper()

	if want != val.State() {
		t.Errorf("state should be: %s but is: %s", want, val.State())
	}
}

func checkJSON[T any](t *testing.T, v Val[T], s string) {
	t.Helper()

	b, err := opt.JSONMarshal(v)
	if err != nil {
		t.Error(err)
	}

	if string(b) != s {
		t.Errorf("expect: %s, got: %s", s, b)
	}
}
