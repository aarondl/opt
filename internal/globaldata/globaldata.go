package globaldata

import (
	"database/sql/driver"
	"encoding"
	"reflect"
)

var (
	JSONNull                    = []byte("null")
	DriverValuerIntf            = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	EncodingTextMarshalerIntf   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	EncodingTextUnmarshalerIntf = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	Int64Type                   = reflect.TypeOf(int64(0))
	ByteSliceType               = reflect.TypeOf([]byte(nil))
)
