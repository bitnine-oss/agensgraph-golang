package ag

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// Array returns Scanner and Valuer for arr that is an array or a slice of
// the following types; GraphId.
func Array(arr interface{}) interface {
	sql.Scanner
	driver.Valuer
} {
	switch arr := arr.(type) {
	case *[]GraphId:
		return (*graphIdArray)(arr)
	case []GraphId:
		return (*graphIdArray)(&arr)

	default:
		panic(fmt.Errorf("%T not supported", arr))
	}
}
