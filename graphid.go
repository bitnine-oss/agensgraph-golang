package ag

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
)

// GraphId is a unique ID for a vertex and an edge.
type GraphId struct {
	// Valid is true if GraphId is not NULL
	Valid bool

	b []byte
}

var graphIdRegexp = regexp.MustCompile(`^(\d+)\.(\d+)$`)

// NewGraphId returns GraphId of str if str is between "0.0" and
// "65535.281474976710656". Otherwise, it returns error.
func NewGraphId(str string) (GraphId, error) {
	m := graphIdRegexp.FindStringSubmatch(str)
	if m == nil {
		return GraphId{}, fmt.Errorf("bad graphid representation: %q", str)
	}

	var err error

	_, err = strconv.ParseUint(m[1], 10, 16)
	if err != nil {
		return GraphId{}, fmt.Errorf("invalid label ID: %s", err)
	}

	_, err = strconv.ParseUint(m[2], 10, 48)
	if err != nil {
		return GraphId{}, fmt.Errorf("invalid local ID: %s", err)
	}

	return GraphId{true, []byte(str)}, nil
}

// Equal reports whether gid and x are the same GraphId.
func (gid GraphId) Equal(x GraphId) bool {
	if !gid.Valid || !x.Valid {
		return false
	}
	return bytes.Equal(gid.b, x.b)
}

func (gid GraphId) String() string {
	if gid.Valid {
		return string(gid.b)
	} else {
		return "NULL"
	}
}

// Scan implements the database/sql Scanner interface.
func (gid *GraphId) Scan(src interface{}) error {
	if src == nil {
		gid.Valid, gid.b = false, nil
		return nil
	}

	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for graphid: %T", src)
	}

	gid.b = append([]byte(nil), b...)
	gid.Valid = true
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (gid GraphId) Value() (driver.Value, error) {
	if gid.Valid {
		return gid.b, nil
	} else {
		return nil, nil
	}
}
