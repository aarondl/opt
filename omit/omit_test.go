package omit_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/aarondl/opt/omit"
)

func TestConstruction(t *testing.T) {
	t.Parallel()

	hello := "hello"

	val := omit.From("hello")
	checkState(t, val, omit.StateSet)
	if !val.IsSet() {
		t.Error("should be set")
	}

	val = omit.FromPtr(&hello)
	checkState(t, val, omit.StateSet)
	val = omit.FromPtr[string](nil)
	checkState(t, val, omit.StateUnset)
	if !val.IsUnset() {
		t.Error("should be unset")
	}

	val = omit.New[string]()
	checkState(t, val, omit.StateUnset)
	if !val.IsUnset() {
		t.Error("should be unset")
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	val := omit.From("hello")
	if val.MustGet() != "hello" {
		t.Error("wrong value")
	}
	if val.GetOr("hi") != "hello" {
		t.Error("wrong value")
	}
	if *val.Ptr() != "hello" {
		t.Error("wrong value")
	}

	val.Unset()
	if _, ok := val.Get(); ok {
		t.Error("should not be okay")
	}
	if val.GetOr("hi") != "hi" {
		t.Error("wrong value")
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

	val := omit.New[int]()
	if !val.Map(func(int) int { return 0 }).IsUnset() {
		t.Error("it should still be unset")
	}
	if !omit.Map(val, func(int) int { return 0 }).IsUnset() {
		t.Error("it should still be unset")
	}
	val.Set(5)
	if val.Map(func(i int) int { return i + 1 }).MustGet() != 6 {
		t.Error("wrong value")
	}
	if omit.Map(val, func(i int) int { return i + 1 }).MustGet() != 6 {
		t.Error("wrong value")
	}
}

func TestChanges(t *testing.T) {
	t.Parallel()

	val := omit.From("hello")
	checkState(t, val, omit.StateSet)
	val.Unset()
	checkState(t, val, omit.StateUnset)

	val = omit.New[string]()
	checkState(t, val, omit.StateUnset)
	val.Set("hello")
	checkState(t, val, omit.StateSet)

	val = omit.New[string]()
	checkState(t, val, omit.StateUnset)
	val.SetPtr(nil)
	checkState(t, val, omit.StateUnset)

	hello := "hello"
	val = omit.New[string]()
	checkState(t, val, omit.StateUnset)
	val.SetPtr(&hello)
	checkState(t, val, omit.StateSet)
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	val := omit.From("hello")
	checkJSON(t, val, `"hello"`)
	val.Unset()
	checkJSON(t, val, `null`)
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	hello := omit.New[string]()
	checkState(t, hello, omit.StateUnset)

	if err := json.Unmarshal([]byte("null"), &hello); err == nil {
		t.Error("cannot accept a null")
	}

	if err := json.Unmarshal([]byte(`"hello"`), &hello); err != nil {
		t.Error(err)
	}
	checkState(t, hello, omit.StateSet)

	if hello.MustGet() != "hello" {
		t.Error("expected hello")
	}

	hello.UnmarshalJSON(nil)
	checkState(t, hello, omit.StateUnset)
}

func TestMarshalText(t *testing.T) {
	t.Parallel()

	hello := omit.From("hello")
	b, err := hello.MarshalText()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "hello" {
		t.Error("expected hello")
	}

	hello.Unset()
	b, err = hello.MarshalText()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "" {
		t.Error("expected empty str")
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Parallel()

	var val omit.Val[string]
	if err := val.UnmarshalText([]byte("")); err != nil {
		t.Error(err)
	}
	checkState(t, val, omit.StateUnset)

	if err := val.UnmarshalText([]byte("hello")); err != nil {
		t.Error(err)
	}
	checkState(t, val, omit.StateSet)
	if val.MustGet() != "hello" {
		t.Error("wrong value")
	}
}

func TestScan(t *testing.T) {
	t.Parallel()

	var val omit.Val[string]
	if err := val.Scan(nil); err == nil {
		t.Error("should break trying to scan null")
	}

	if err := val.Scan("hello"); err != nil {
		t.Error(err)
	}
	checkState(t, val, omit.StateSet)
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

	var val omit.Val[string]
	if v, err := val.Value(); err != nil {
		t.Error(err)
	} else if v != nil {
		t.Error("expected v to be nil")
	}

	val = omit.From("hello")
	if v, err := val.Value(); err != nil {
		t.Error(err)
	} else if v.(string) != "hello" {
		t.Error("expected v to be nil")
	}

	date := time.Date(2000, 1, 1, 2, 30, 0, 0, time.UTC)
	nullTime := omit.From(date)
	if v, err := nullTime.Value(); err != nil {
		t.Error(err)
	} else if !v.(time.Time).Equal(date) {
		t.Error("time was wrong")
	}

	valuer := omit.From(valuerImplementation{})
	if v, err := valuer.Value(); err != nil {
		t.Error(err)
	} else if v.(int64) != 1 {
		t.Error("expect const int")
	}
}

func TestBadValuers(t *testing.T) {
	t.Parallel()

	type noValuerImplementation struct{}
	nonValuer := omit.From(noValuerImplementation{})
	if _, err := nonValuer.Value(); err == nil {
		t.Error("expect an error to be returned")
	}

	nonByteSlice := omit.From([]string{"slice"})
	if _, err := nonByteSlice.Value(); err == nil {
		t.Error("expect an error to be returned")
	}

	nonSupportedType := omit.From(make(chan string))
	if _, err := nonSupportedType.Value(); err == nil {
		t.Error("expect an error to be returned")
	}
}

func TestStateStringer(t *testing.T) {
	t.Parallel()

	if omit.StateUnset.String() != "unset" {
		t.Error("bad value")
	}
	if omit.StateSet.String() != "set" {
		t.Error("bad value")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic")
		}
	}()
	_ = omit.State(99).String()
}

func checkState[T any](t *testing.T, val omit.Val[T], state omit.State) {
	t.Helper()

	if state != val.State() {
		t.Errorf("state should be: %s but is: %s", state, val.State())
	}
}

func checkJSON[T any](t *testing.T, v omit.Val[T], s string) {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	if string(b) != s {
		t.Errorf("expect: %s, got: %s", s, b)
	}
}
