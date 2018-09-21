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

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	agTestServer    *bool
	agTestGraphName = "ag_test_go"
)

func TestMain(m *testing.M) {
	agTestServer = flag.Bool("ag.test.server", false, "Run server tests")
	flag.Parse()

	setUpServerTest()
	c := m.Run()
	tearDownServerTest()

	os.Exit(c)
}

func setUpServerTest() {
	if !*agTestServer {
		return
	}

	// Other environment variables: PGHOST, PGPORT, PGUSER, PGPASSWORD, ...
	os.Setenv("PGDATABASE", "postgres")
	os.Setenv("PGSSLMODE", "disable")

	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := fmt.Sprintf(`CREATE GRAPH IF NOT EXISTS %q`, agTestGraphName)
	_, err = db.Exec(q)
	if err != nil {
		log.Fatal(err)
	}
}

func tearDownServerTest() {
	if !*agTestServer {
		return
	}

	db, err := sql.Open("postgres", "")
	if err == nil {
		q := fmt.Sprintf(`DROP GRAPH %q CASCADE`, agTestGraphName)
		db.Exec(q)
		db.Close()
	}
}

func skipUnlessServerTest(t *testing.T) {
	if !*agTestServer {
		t.SkipNow()
	}
}

func mustOpenAndSetGraph(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "")
	if err != nil {
		t.Fatal(err)
	}
	q := fmt.Sprintf(`SET graph_path = %q`, agTestGraphName)
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
