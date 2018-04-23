package ag

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/lib/pq"
)

// Array returns Scanner and Valuer for arr that is an array or a slice of
// the following types; GraphId.
func Array(arr interface{}) interface {
	sql.Scanner
	driver.Valuer
} {
	switch arr := arr.(type) {
	case []GraphId:
		return pq.GenericArray{arr}
	case *[]GraphId:
		return pq.GenericArray{arr}
	default:
		panic(fmt.Errorf("unexpected type %T for Array", arr))
	}
}
