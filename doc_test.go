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

func ExampleArray_graphId() {
	var gids []ag.GraphId

	// Scan()
	db.QueryRow(`MATCH ()-[es*]->() RETURN [e IN es | id(e)] LIMIT 1`).Scan(ag.Array(&gids))

	// Value()
	gid0, _ := ag.NewGraphId("1.1")
	gid1, _ := ag.NewGraphId("65535.281474976710655")
	gids = []ag.GraphId{gid0, gid1}
	db.QueryRow(`MATCH (n) WHERE id(n) IN $1 RETURN n LIMIT 1`, ag.Array(gids))
}
