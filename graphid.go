package ag

import (
	"bytes"
	"database/sql/driver"
	"errors"
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

var nullGraphId = GraphId{}

var graphIdRegexp = regexp.MustCompile(`^(\d+)\.(\d+)$`)

// NewGraphId returns GraphId of str if str is between "1.1" and
// "65535.281474976710656". If str is "NULL", it returns GraphId whose Valid is
// false. Otherwise, it returns error.
func NewGraphId(str string) (GraphId, error) {
	if str == "NULL" {
		return nullGraphId, nil
	}

	m := graphIdRegexp.FindStringSubmatch(str)
	if m == nil {
		return GraphId{}, fmt.Errorf("bad graphid representation: %q", str)
	}

	i, err := strconv.ParseUint(m[1], 10, 16)
	if err != nil {
		return GraphId{}, errors.New("invalid label ID: " + err.Error())
	}
	if i == 0 {
		return GraphId{}, fmt.Errorf("invalid label ID: %d", i)
	}

	i, err = strconv.ParseUint(m[2], 10, 48)
	if err != nil {
		return GraphId{}, errors.New("invalid local ID: " + err.Error())
	}
	if i == 0 {
		return GraphId{}, fmt.Errorf("invalid local ID: %d", i)
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

type graphIdArray []GraphId

// separated by comma (see graphid in pg_type.h)
const graphIdDelim = byte(054)

var graphIdNullValue = []byte("NULL")

func (a *graphIdArray) Scan(src interface{}) error {
	if src == nil {
		*a = nil
		return nil
	}

	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for _graphid: %T", src)
	}

	// remove surrounding braces
	b = b[1 : len(b)-1]

	// bytes.Split() returns [][]byte{[]byte{}} even if len(b) < 1.
	// In this case, return empty []GraphId to distinguish between NULL and
	// empty _graphid.
	if len(b) < 1 {
		*a = []GraphId{}
		return nil
	}

	sa := bytes.Split(b, []byte{graphIdDelim})

	gids := make([]GraphId, len(sa))
	for i, s := range sa {
		if bytes.Equal(s, graphIdNullValue) {
			gids[i] = nullGraphId
			continue
		}

		err := gids[i].Scan(s)
		if err != nil {
			return errors.New("bad _graphid representation: " + err.Error())
		}
	}

	*a = gids
	return nil
}

func (a graphIdArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// '{' + "d.d,"*n - ',' + '}' = 1+4*n-1+1 = 4*n+1
		b := make([]byte, 1, 4*n+1)
		b[0] = byte('{')
		for i, gid := range a {
			if i > 0 {
				b = append(b, graphIdDelim)
			}
			val, _ := gid.Value()
			if val == nil {
				b = append(b, graphIdNullValue...)
			} else {
				b = append(b, val.([]byte)...)
			}
		}
		b = append(b, '}')

		return b, nil
	}

	return []byte("{}"), nil
}
