package opt

import (
	"database/sql/driver"
	"encoding"
	"reflect"

	"github.com/aarondl/opt/internal/globaldata"
)

// ToDriverValue generates the appropriate driver.Value
// from a given value
func ToDriverValue(val any) (driver.Value, error) {
	refVal := reflect.ValueOf(val)
	if refVal.Type().Implements(globaldata.DriverValuerIntf) {
		valuer := refVal.Interface().(driver.Valuer)
		return valuer.Value()
	}

	// If it is of time.Time type, return it as is.
	if refVal.Type() == globaldata.TimeType {
		return val, nil
	}

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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return refVal.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(refVal.Uint()), nil
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
