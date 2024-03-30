package testutils

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewTestDB creates a new connection to the integration test database and also registers a cleanup
// function to teardown the tables in the database and close the connection
func NewTestDB(t *testing.T) *pgxpool.Pool {
	// connect to the test DB
	conn, err := pgxpool.New(
		context.Background(),
		"postgresql://test_web:3c523592-852d-42be-915c-d5931792e39e@localhost:5432/test",
	)
	if err != nil {
		t.Fatal(err)
	}
	// create necessary tables and insert necessary data
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(context.Background(), string(script)); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		// remove all the tables at the end of the test and return the connection to the pool
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := conn.Exec(context.Background(), string(script)); err != nil {
			t.Fatal(err)
		}
		conn.Close()
	})

	return conn
}
