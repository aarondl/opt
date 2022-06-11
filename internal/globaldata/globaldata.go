package globaldata

import (
	"database/sql/driver"
	"encoding"
	"reflect"
	"time"
)

var (
	JSONNull                          = []byte("null")
	DriverValuerIntf                  = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	EncodingTextMarshalerIntf         = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	EncodingTextUnmarshalerIntf       = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	TimeTimeType                      = reflect.TypeOf(time.Time{})
	Int64Val                    int64 = 0
	Int64Type                         = reflect.TypeOf(Int64Val)
)
