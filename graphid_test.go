package ag

import (
	"bytes"
	"database/sql/driver"
	"testing"
)

func TestNewGraphId(t *testing.T) {
	tests := []struct {
		str string
		ok  bool
	}{
		{"NULL", true},
		{"-1.0", false},
		{"0.-1", false},
		{"0.0", false},
		{"1.1", true},
		{"65535.281474976710655", true},
		{"65536.281474976710655", false},
		{"65535.281474976710656", false},
	}
	for _, c := range tests {
		_, err := NewGraphId(c.str)
		if err == nil {
			if !c.ok {
				t.Errorf("error expected for %q", c.str)
			}
		} else {
			if c.ok {
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
		{mustNewGraphId("NULL"), mustNewGraphId("NULL"), false},
		{mustNewGraphId("NULL"), mustNewGraphId("1.1"), false},
		{mustNewGraphId("1.1"), mustNewGraphId("NULL"), false},
		{mustNewGraphId("1.1"), mustNewGraphId("1.1"), true},
		{mustNewGraphId("65535.281474976710655"), mustNewGraphId("65535.281474976710655"), true},
		{mustNewGraphId("1.1"), mustNewGraphId("65535.281474976710655"), false},
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
	src := 0
	var gid GraphId
	err := gid.Scan(src)
	if err == nil {
		t.Errorf("error expected for %T", src)
	}
}

func TestGraphIdScanDuplicate(t *testing.T) {
	src := []byte("1.1")

	var gid GraphId
	_ = gid.Scan(src)

	if &gid.b[0] == &src[0] {
		t.Error("GraphId references underlying array")
	}
}

func TestGraphIdArrayScan(t *testing.T) {
	tests := []struct {
		src  interface{}
		gids []GraphId
	}{
		{
			[]byte("{NULL,1.1,65535.281474976710655}"),
			[]GraphId{
				mustNewGraphId("NULL"),
				mustNewGraphId("1.1"),
				mustNewGraphId("65535.281474976710655"),
			},
		},
		{
			[]byte("{}"),
			[]GraphId{},
		},
		{
			nil,
			nil,
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
			if !gid.Valid && !c.gids[i].Valid {
				continue
			}
			if !gid.Equal(c.gids[i]) {
				t.Errorf("got %q, want %q", gid, c.gids[i])
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
				mustNewGraphId("NULL"),
				mustNewGraphId("1.1"),
				mustNewGraphId("65535.281474976710655"),
			},
			[]byte("{NULL,1.1,65535.281474976710655}"),
		},
		{
			[]GraphId{},
			[]byte("{}"),
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
			b, ok := val.([]byte)
			if !ok {
				t.Errorf("got %T, want []byte", val)
				continue
			}
			if !bytes.Equal(b, c.val.([]byte)) {
				t.Errorf("got %q, want %q", string(b), string(c.val.([]byte)))
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
		t.Error("got NULL, want Valid GraphId")
	}

	var cnt int64
	q := `MATCH (n:gid) WHERE id(n) = $1 RETURN count(*)`
	err = db.QueryRow(q, gid).Scan(&cnt)
	if err != nil {
		t.Fatal(err)
	}
	if cnt != 1 {
		t.Errorf("got %d, want %d", cnt, 1)
	}

	gids := []GraphId{mustNewGraphId("1.1"), mustNewGraphId("NULL"), mustNewGraphId("65535.281474976710655")}
	var gidsOut []GraphId
	err = db.QueryRow(`SELECT $1::_graphid`, Array(gids)).Scan(Array(&gidsOut))
	if err != nil {
		t.Fatal(err)
	}
	for i, gid := range gidsOut {
		if !gid.Valid && !gids[i].Valid {
			continue
		}
		if !gid.Equal(gids[i]) {
			t.Errorf("got %q, want %q", gid, gids[i])
			break
		}
	}

	err = db.QueryRow(`SELECT NULL::_graphid`).Scan(Array(&gidsOut))
	if err != nil {
		t.Fatal(err)
	}
	if gidsOut != nil {
		t.Errorf("got %v, want nil", gidsOut)
	}
}
