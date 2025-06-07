package opt

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"reflect"

	"github.com/aarondl/opt/internal/globaldata"
)

// callValuerValue returns vr.Value(), with one exception:
// If vr.Value is an auto-generated method on a pointer type and the
// pointer is nil, it would panic at runtime in the panicwrap
// method. Treat it like nil instead.
// Issue 8415.
//
// This is so people can implement driver.Value on value types and
// still use nil pointers to those types to mean nil/NULL, just like
// string/*string.
//
// This function is mirrored in the database/sql package.
func callValuerValue(vr driver.Valuer) (v driver.Value, err error) {
	if rv := reflect.ValueOf(vr); rv.Kind() == reflect.Pointer &&
		rv.IsNil() &&
		rv.Type().Elem().Implements(globaldata.DriverValuerIntf) {
		return nil, nil
	}
	return vr.Value()
}

// ToDriverValue generates the appropriate driver.Value
// from a given value
func ToDriverValue(val any) (driver.Value, error) {
	switch vr := val.(type) {
	case driver.Valuer:
		sv, err := callValuerValue(vr)
		if err != nil {
			return nil, err
		}
		if !driver.IsValue(sv) {
			return nil, fmt.Errorf("non-Value type %T returned from Value", sv)
		}
		return sv, nil

	// For now, continue to prefer the Valuer interface
	// over the decimal decompose interface.
	case decimalDecompose:
		return vr, nil
	}

	if driver.IsValue(val) {
		return val, nil
	}

	refVal := reflect.ValueOf(val)

	// If it implements encoding.TextMarshaler, use that.
	if refVal.Type().Implements(globaldata.EncodingTextMarshalerIntf) {
		marshaler := refVal.Interface().(encoding.TextMarshaler)
		return marshaler.MarshalText()
	}

	// If it implements encoding.BinaryMarshaler, use that.
	if refVal.Type().Implements(globaldata.EncodingBinaryMarshalerIntf) {
		marshaler := refVal.Interface().(encoding.BinaryMarshaler)
		return marshaler.MarshalBinary()
	}

	switch refVal.Kind() {
	case reflect.Pointer:
		// indirect pointers
		if refVal.IsNil() {
			return nil, nil
		} else {
			return ToDriverValue(refVal.Elem().Interface())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return refVal.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(refVal.Uint()), nil
	case reflect.Uint64:
		u64 := refVal.Uint()
		if u64 >= 1<<63 {
			// driver.Value only supports int64, so if the high bit is set,
			// the value cannot be represented as an int64.
			// And converting it to int64 would result in a negative value,
			return nil, fmt.Errorf("uint64 values with high bit set are not supported")
		}
		return int64(u64), nil
	case reflect.Float32, reflect.Float64:
		return refVal.Float(), nil
	case reflect.Bool:
		return refVal.Bool(), nil
	case reflect.Slice:
		if refVal.Type().Elem().Kind() == reflect.Uint8 {
			return refVal.Bytes(), nil
		}
	case reflect.String:
		return refVal.String(), nil
	}

	return val, nil
}
