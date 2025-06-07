package globaldata

import (
	"database/sql/driver"
	"encoding"
	"reflect"
)

var (
	JSONNull                      = []byte("null")
	DriverValuerIntf              = reflect.TypeFor[driver.Valuer]()
	EncodingTextMarshalerIntf     = reflect.TypeFor[encoding.TextMarshaler]()
	EncodingTextUnmarshalerIntf   = reflect.TypeFor[encoding.TextUnmarshaler]()
	EncodingBinaryMarshalerIntf   = reflect.TypeFor[encoding.BinaryMarshaler]()
	EncodingBinaryUnmarshalerIntf = reflect.TypeFor[encoding.BinaryUnmarshaler]()
)
