package null

import (
	"bytes"
	"database/sql/driver"
	"net"
	"testing"
	"time"

	"github.com/aarondl/opt"
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

	val = FromCond("hello", true)
	checkState(t, val, StateSet)
	val = FromCond("hello", false)
	checkState(t, val, StateNull)
	if !val.IsNull() {
		t.Error("should be null")
	}

	val = Val[string]{}
	checkState(t, val, StateNull)
	if !val.IsNull() {
		t.Error("should be null")
	}
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
	if *val.Ptr() != "hello" {
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
	if val.Ptr() != nil {
		t.Error("should be nil")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("should have panic'd")
		}
	}()
	_ = val.MustGet()
}

func TestMap(t *testing.T) {
	t.Parallel()

	val := Val[int]{}
	if !val.Map(func(int) int { return 0 }).IsNull() {
		t.Error("it should still be null")
	}
	if !Map(val, func(int) int { return 0 }).IsNull() {
		t.Error("it should still be null")
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

	val = Val[string]{}
	checkState(t, val, StateNull)
	val.Set("hello")
	checkState(t, val, StateSet)

	val = Val[string]{}
	checkState(t, val, StateNull)
	val.SetPtr(nil)
	checkState(t, val, StateNull)

	hello := "hello"
	val = Val[string]{}
	checkState(t, val, StateNull)
	val.SetPtr(&hello)
	checkState(t, val, StateSet)
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	val := From("hello")
	checkJSON(t, val, `"hello"`)
	val.Null()
	checkJSON(t, val, `null`)
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	hello := Val[string]{}
	checkState(t, hello, StateNull)

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

	if err := hello.UnmarshalJSON(nil); err == nil {
		t.Error("expected an error")
	}
}

func TestMarshalText(t *testing.T) {
	t.Parallel()

	hello := From("hello")
	b, err := hello.MarshalText()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "hello" {
		t.Error("expected hello")
	}

	hello.Null()
	b, err = hello.MarshalText()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "" {
		t.Error("expected hello")
	}

	marshaller := From(net.IPv4(1, 1, 1, 1))
	if b, err := marshaller.MarshalText(); err != nil {
		t.Error(err)
	} else if !bytes.Equal(b, []byte("1.1.1.1")) {
		t.Error("wrong value")
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Parallel()

	var val Val[string]
	checkState(t, val, StateNull)
	if err := val.UnmarshalText([]byte("hello")); err != nil {
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
	checkState(t, unmarshaller, StateNull)

	if err := unmarshaller.UnmarshalText([]byte("1.1.1.1")); err != nil {
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

func checkState[T any](t *testing.T, val Val[T], state state) {
	t.Helper()

	if state != val.State() {
		t.Errorf("state should be: %s but is: %s", state, val.State())
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
