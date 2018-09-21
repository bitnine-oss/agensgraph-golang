# ag - Go driver support for AgensGraph

## Documentation

Please see the package documentation at https://godoc.org/github.com/bitnine-oss/agensgraph-golang for the detailed documentation and basic usage examples.

## Tests

You may run `go test` with the optional `-ag.test.server` flag for server test.

For the server test, the environment variables listed at [here](https://www.postgresql.org/docs/10/static/libpq-envars.html) can be used to set connection parameter values. There are two environment variables set by the test code; `PGDATABASE=postgres` and `PGSSLMODE=disable`.

## License

Go driver support for AgensGraph is licensed under the GNU Affero General Public License.
