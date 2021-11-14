package test

import (
	"database/sql"
	"testing"

	"dariobalinzo.com/elastic/v2/config"
	_ "github.com/lib/pq"
)

type readerIntegrationTest struct{}

// 'readerTest' implementing the Interface
func (r *readerIntegrationTest) ReadConfig() ([]byte, error) {

	testComplexConfig := `
sql: "postgres://user:pwd@localhost:5432?sslmode=disable"
queries:
  - query: "select * from table1 where last_updated > $1"
    index: "index1"
`

	// Prepare data you want to return without reading from the file
	return []byte(testComplexConfig), nil
}

func TestIntegration(t *testing.T) {
	testReader := readerComplexTest{fileName: "Sample File Name"}
	conf, err := config.ParseConfiguration(&testReader)
	if err != nil {
		t.Errorf("Parsing config, got error %s", err)
	}

	db, err := sql.Open("postgres", conf.SqlConfig)
	if err != nil {
		t.Error("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP IF EXISTS table1")
	if err != nil {
		t.Error("Error dropping temp table: ", err)
	}
	db.Exec(`CREATE TABLE table1 (
		user_id serial PRIMARY KEY,
		email VARCHAR ( 255 ) UNIQUE NOT NULL,
		last_updated TIMESTAMP NOT NULL
	)`)
	if err != nil {
		t.Error("Error creating temp table: ", err)
	}

	db.Exec(`
	INSERT INTO table1 (email,last_updated) VALUES('mario@rossi.it', now())
	`)
	if err != nil {
		t.Error("Error insert temp table: ", err)
	}

}
