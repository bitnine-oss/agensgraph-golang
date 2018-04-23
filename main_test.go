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
