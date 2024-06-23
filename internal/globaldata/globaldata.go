package globaldata

import (
	"database/sql/driver"
	"encoding"
	"reflect"
	"time"
)

var (
	JSONNull                      = []byte("null")
	DriverValuerIntf              = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	EncodingTextMarshalerIntf     = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	EncodingTextUnmarshalerIntf   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	EncodingBinaryMarshalerIntf   = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	EncodingBinaryUnmarshalerIntf = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()
	TimeType                      = reflect.TypeOf(time.Time{})
	Int64Type                     = reflect.TypeOf(int64(0))
	ByteSliceType                 = reflect.TypeOf([]byte(nil))
)
