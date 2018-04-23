package ag_test

import (
	"database/sql"

	"github.com/bitnine-oss/agensgraph-golang"
)

var (
	db *sql.DB
)

func ExampleGraphId_Scan() {
	var gid ag.GraphId
	db.QueryRow(`MATCH (n) RETURN id(n) LIMIT 1`).Scan(&gid)
}

func ExampleGraphId_Value() {
	gid, _ := ag.NewGraphId("1.1")
	db.QueryRow(`MATCH (n) WHERE id(n) = $1 RETURN n`, gid)
}
