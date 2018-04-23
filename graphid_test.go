package ag

import (
	"database/sql/driver"
	"testing"
)

func TestNewGraphId(t *testing.T) {
	tests := []struct {
		str   string
		valid bool
	}{
		{"-1.0", false},
		{"0.-1", false},
		{"0.0", true},
		{"65535.281474976710655", true},
		{"65536.281474976710655", false},
		{"65535.281474976710656", false},
	}
	for _, c := range tests {
		_, err := NewGraphId(c.str)
		if err == nil {
			if !c.valid {
				t.Errorf("error expected for %q", c.str)
			}
		} else {
			if c.valid {
				t.Error(err)
			}
		}
	}
}

func mustNewGraphId(str string) GraphId {
	gid, err := NewGraphId(str)
	if err != nil {
		panic(err)
	}
	return gid
}

func TestGraphIdEqual(t *testing.T) {
	tests := []struct {
		x     GraphId
		y     GraphId
		equal bool
	}{
		{GraphId{}, GraphId{}, false},
		{GraphId{}, mustNewGraphId("0.0"), false},
		{mustNewGraphId("0.0"), GraphId{}, false},
		{mustNewGraphId("0.0"), mustNewGraphId("0.0"), true},
		{mustNewGraphId("65535.281474976710655"), mustNewGraphId("65535.281474976710655"), true},
		{mustNewGraphId("0.0"), mustNewGraphId("65535.281474976710655"), false},
	}
	for _, c := range tests {
		if c.x.Equal(c.y) != c.equal {
			t.Errorf("got %q.Equal(%q) == %t, want %t", c.x, c.y, c.x.Equal(c.y), c.equal)
		}
	}
}

func TestGraphIdScanNil(t *testing.T) {
	var gid GraphId
	err := gid.Scan(nil)
	if err != nil {
		t.Fatal(err)
	}
	if gid.Valid {
		t.Errorf("got %q, want NULL", gid)
	}
}

func TestGraphIdScanType(t *testing.T) {
	var src interface{}
	var gid GraphId
	var err error

	src = 0
	err = gid.Scan(src)
	if err == nil {
		t.Errorf("error expected for %T", src)
	}

	src = []byte{}
	err = gid.Scan(src)
	if err != nil {
		t.Fatal(err)
	}
	if !gid.Valid {
		t.Errorf("got %s, want %q", gid, string(src.([]byte)))
	}
}

func TestGraphIdScanDuplicate(t *testing.T) {
	src := []byte("0.0")

	var gid GraphId
	_ = gid.Scan(src)

	if &gid.b[0] == &src[0] {
		t.Error("GraphId references underlying array")
	}
}

func TestGraphIdArrayScan(t *testing.T) {
	tests := []struct {
		src  []byte
		gids []GraphId
	}{
		{
			[]byte("{0.0,65535.281474976710655}"),
			[]GraphId{
				mustNewGraphId("0.0"),
				mustNewGraphId("65535.281474976710655"),
			},
		},
	}
	for _, c := range tests {
		var gids []GraphId
		a := Array(&gids)
		err := a.Scan(c.src)
		if err != nil {
			t.Error(err)
			continue
		}

		if len(gids) != len(c.gids) {
			t.Errorf("got len(gids) == %d, want %d", len(gids), len(c.gids))
			continue
		}

	EqualLoop:
		for i, gid := range gids {
			if !gid.Equal(c.gids[i]) {
				t.Errorf("got %s, want %s", gid, c.gids[i])
				break EqualLoop
			}
		}
	}
}

func TestGraphIdArrayValue(t *testing.T) {
	tests := []struct {
		gids []GraphId
		val  driver.Value
	}{
		{
			[]GraphId{
				mustNewGraphId("0.0"),
				mustNewGraphId("65535.281474976710655"),
			},
			`{"0.0","65535.281474976710655"}`,
		},
		{
			nil,
			nil,
		},
	}
	for _, c := range tests {
		val, err := Array(c.gids).Value()
		if err != nil {
			t.Error(err)
			continue
		}

		if val == nil {
			if c.val != nil {
				t.Errorf("error expected for %v", c.gids)
			}
		} else {
			s, ok := val.(string)
			if !ok {
				t.Errorf("got %T, want string", val)
				continue
			}
			if s != c.val.(string) {
				t.Errorf("got %s, want %s", s, c.val.(string))
			}
		}
	}
}

func TestServerGraphId(t *testing.T) {
	skipUnlessServerTest(t)

	db := mustOpenAndSetGraph(t)
	defer db.Close()

	_, err := db.Exec(`CREATE (:gid)`)
	if err != nil {
		t.Fatal(err)
	}

	var gid GraphId
	err = db.QueryRow(`MATCH (n:gid) RETURN id(n)`).Scan(&gid)
	if err != nil {
		t.Fatal(err)
	}
	if !gid.Valid {
		t.Fatal("got NULL, want Valid GraphId")
	}

	var cnt int64
	q := `MATCH (n:gid) WHERE id(n) = $1 RETURN count(*)`
	err = db.QueryRow(q, gid).Scan(&cnt)
	if err != nil {
		t.Fatal(err)
	}
	if cnt != 1 {
		t.Errorf("got %d, want 1", cnt)
	}
}
