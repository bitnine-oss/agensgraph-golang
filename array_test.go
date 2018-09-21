/*
 * Copyright (c) 2014-2018 Bitnine, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package ag

import "testing"

type testElement struct{}

func (_ testElement) readElements(b []byte) ([]interface{}, error) {
	return []interface{}{}, nil
}

func TestArrayScanError(t *testing.T) {
	tests := []interface{}{
		nil,
		0,
		(*byte)(nil),
		&testElement{},
		&[]byte{},
		&[]testElement{},
	}
	for _, i := range tests {
		err := Array(i).Scan(nil)
		if err == nil {
			t.Errorf("error expected for %T", i)
		}
	}
}

type testElementEx struct {
	testElement
}

func (_ testElementEx) Scan(src interface{}) error {
	return nil
}

func TestArrayScanNullSlice(t *testing.T) {
	var es []testElementEx
	err := Array(&es).Scan(nil)
	if err != nil {
		t.Fatalf("got %v, want NULL", es)
	}

	if es != nil {
		t.Fatalf("got %v, want nil", es)
	}
}

func TestArrayScanNullArray(t *testing.T) {
	es := [1]testElementEx{}
	err := Array(&es).Scan(nil)
	if err == nil {
		t.Fatalf("got %v, want error", es)
	}

	if _, ok := err.(NullArrayError); !ok {
		t.Errorf("NullArrayError expected")
	}
}

func TestArrayScanType(t *testing.T) {
	src := 0
	var es []testElementEx
	err := Array(&es).Scan(src)
	if err == nil {
		t.Errorf("error expected for %T", src)
	}
}

func TestArrayScanZero(t *testing.T) {
	var src interface{} = []byte(nil)
	var es []testElementEx
	err := Array(&es).Scan(src)
	if err == nil {
		t.Errorf("error expected for %T", src)
	}
}

func TestArrayScanArrayLen(t *testing.T) {
	es := [1]testElementEx{}
	err := Array(&es).Scan([]byte("dummy"))
	if err == nil {
		t.Errorf("error expected for %T", es)
	}
}

func TestArrayValue(t *testing.T) {
	var es []testElementEx
	_, err := Array(&es).Value()
	if err == nil {
		t.Errorf("error expected for Value() on Array")
	}
}
