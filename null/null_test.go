package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/aarondl/opt/null"
)

func TestConstruction(t *testing.T) {
	t.Parallel()

	hello := "hello"

	val := null.From("hello")
	checkState(t, val, null.StateSet)
	if !val.IsSet() {
		t.Error("should be set")
	}

	val = null.FromPtr(&hello)
	checkState(t, val, null.StateSet)
	val = null.FromPtr[string](nil)
	checkState(t, val, null.StateNull)
	if !val.IsNull() {
		t.Error("should be null")
	}

	val = null.New[string]()
	checkState(t, val, null.StateNull)
	if !val.IsNull() {
		t.Error("should be null")
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	val := null.From("hello")
	if val.MustGet() != "hello" {
		t.Error("wrong value")
	}
	if val.GetOr("hi") != "hello" {
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

	val := null.New[int]()
	if !val.Map(func(int) int { return 0 }).IsNull() {
		t.Error("it should still be null")
	}
	if !null.Map(val, func(int) int { return 0 }).IsNull() {
		t.Error("it should still be null")
	}
	val.Set(5)
	if val.Map(func(i int) int { return i + 1 }).MustGet() != 6 {
		t.Error("wrong value")
	}
	if null.Map(val, func(i int) int { return i + 1 }).MustGet() != 6 {
		t.Error("wrong value")
	}
}

func TestChanges(t *testing.T) {
	t.Parallel()

	val := null.From("hello")
	checkState(t, val, null.StateSet)
	val.Null()
	checkState(t, val, null.StateNull)

	val = null.New[string]()
	checkState(t, val, null.StateNull)
	val.Set("hello")
	checkState(t, val, null.StateSet)

	val = null.New[string]()
	checkState(t, val, null.StateNull)
	val.SetPtr(nil)
	checkState(t, val, null.StateNull)

	hello := "hello"
	val = null.New[string]()
	checkState(t, val, null.StateNull)
	val.SetPtr(&hello)
	checkState(t, val, null.StateSet)
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	val := null.From("hello")
	checkJSON(t, val, `"hello"`)
	val.Null()
	checkJSON(t, val, `null`)
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	hello := null.New[string]()
	checkState(t, hello, null.StateNull)

	if err := json.Unmarshal([]byte("null"), &hello); err != nil {
		t.Error(err)
	}
	checkState(t, hello, null.StateNull)

	if err := json.Unmarshal([]byte(`"hello"`), &hello); err != nil {
		t.Error(err)
	}
	checkState(t, hello, null.StateSet)

	if hello.MustGet() != "hello" {
		t.Error("expected hello")
	}

	if err := hello.UnmarshalJSON(nil); err == nil {
		t.Error("expected an error")
	}
}

func TestMarshalText(t *testing.T) {
	t.Parallel()

	hello := null.From("hello")
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
}

func TestUnmarshalText(t *testing.T) {
	t.Parallel()

	var val null.Val[string]
	checkState(t, val, null.StateNull)
	if err := val.UnmarshalText([]byte("hello")); err != nil {
		t.Error(err)
	}
	checkState(t, val, null.StateSet)
	if val.MustGet() != "hello" {
		t.Error("wrong value")
	}
}

func TestScan(t *testing.T) {
	t.Parallel()

	var val null.Val[string]
	if err := val.Scan(nil); err != nil {
		t.Error(err)
	}
	checkState(t, val, null.StateNull)

	if err := val.Scan("hello"); err != nil {
		t.Error(err)
	}
	checkState(t, val, null.StateSet)
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

	var val null.Val[string]
	if v, err := val.Value(); err != nil {
		t.Error(err)
	} else if v != nil {
		t.Error("expected v to be nil")
	}

	val = null.From("hello")
	if v, err := val.Value(); err != nil {
		t.Error(err)
	} else if v.(string) != "hello" {
		t.Error("expected v to be nil")
	}

	date := time.Date(2000, 1, 1, 2, 30, 0, 0, time.UTC)
	nullTime := null.From(date)
	if v, err := nullTime.Value(); err != nil {
		t.Error(err)
	} else if !v.(time.Time).Equal(date) {
		t.Error("time was wrong")
	}

	valuer := null.From(valuerImplementation{})
	if v, err := valuer.Value(); err != nil {
		t.Error(err)
	} else if v.(int64) != 1 {
		t.Error("expect const int")
	}
}

func TestBadValuers(t *testing.T) {
	t.Parallel()

	type noValuerImplementation struct{}
	nonValuer := null.From(noValuerImplementation{})
	if _, err := nonValuer.Value(); err == nil {
		t.Error("expect an error to be returned")
	}

	nonByteSlice := null.From([]string{"slice"})
	if _, err := nonByteSlice.Value(); err == nil {
		t.Error("expect an error to be returned")
	}

	nonSupportedType := null.From(make(chan string))
	if _, err := nonSupportedType.Value(); err == nil {
		t.Error("expect an error to be returned")
	}
}

func TestStateStringer(t *testing.T) {
	t.Parallel()

	if null.StateNull.String() != "null" {
		t.Error("bad value")
	}
	if null.StateSet.String() != "set" {
		t.Error("bad value")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic")
		}
	}()
	_ = null.State(99).String()
}

func checkState[T any](t *testing.T, val null.Val[T], state null.State) {
	t.Helper()

	if state != val.State() {
		t.Errorf("state should be: %s but is: %s", state, val.State())
	}
}

func checkJSON[T any](t *testing.T, v null.Val[T], s string) {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	if string(b) != s {
		t.Errorf("expect: %s, got: %s", s, b)
	}
}
