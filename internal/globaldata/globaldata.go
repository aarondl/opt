package globaldata

import (
	"database/sql/driver"
	"reflect"
	"time"
)

var (
	JSONNull               = []byte("null")
	DriverValuerIntf       = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	TimeTimeType           = reflect.TypeOf(time.Time{})
	Int64Val         int64 = 0
	Int64Type              = reflect.TypeOf(Int64Val)
)
