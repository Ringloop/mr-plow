package test

import (
	"database/sql"
	"log"
	"testing"

	"dariobalinzo.com/elastic/v2/config"
	"dariobalinzo.com/elastic/v2/elastic"
	"dariobalinzo.com/elastic/v2/movedata"
	_ "github.com/lib/pq"
)

func TestIntegration(t *testing.T) {
	//given (some data on sql db)
	conf := initConfigIntegrationTest(t)
	db := initSqlDB(t, conf)
	defer db.Close()
	elastic.Delete(conf.Queries[0].Index)

	insertData(db, "mario@rossi.it", t)
	originalLastDate, err := elastic.FindLastUpdateOrEpochDate(conf.Queries[0].Index)
	if err != nil {
		t.Error("error in getting last date", err)
		t.FailNow()
	}

	//when (moving data to elastic)
	err = movedata.MoveData(db, conf.Queries[0].Query, conf.Queries[0].Index)
	if err != nil {
		t.Error("error data moving", err)
		t.FailNow()
	}

	//then (last date on elastic should be updated)
	lastImportedDate, err := elastic.FindLastUpdateOrEpochDate(conf.Queries[0].Index)
	if err != nil {
		t.Error("error in getting last date", err)
		t.FailNow()
	}

	log.Println("original date", originalLastDate)
	log.Println("date after import", lastImportedDate)

	if !lastImportedDate.After(*originalLastDate) {
		t.Error("error date not incremented!")
		t.FailNow()
	}

}

type readerIntegrationTest struct{}

// 'readerTest' implementing the Interface
func (r *readerIntegrationTest) ReadConfig() ([]byte, error) {

	testComplexConfig := `
sql: "postgres://user:pwd@localhost:5432/postgres?sslmode=disable"
queries:
  - query: "select * from test.table1 where last_update > $1"
    index: "out_index"
`

	// Prepare data you want to return without reading from the file
	return []byte(testComplexConfig), nil
}

func initConfigIntegrationTest(t *testing.T) *config.ImportConfig {
	testReader := readerIntegrationTest{}
	conf, err := config.ParseConfiguration(&testReader)
	if err != nil {
		t.Error("error reading conf", err)
		t.FailNow()
	}
	return conf
}

func initSqlDB(t *testing.T, conf *config.ImportConfig) *sql.DB {
	db, err := sql.Open("postgres", conf.Database)
	if err != nil {
		t.Error("error connecting to sql db", err)
		t.FailNow()
	}

	_, err = db.Exec(`

	DROP SCHEMA IF EXISTS test CASCADE;
	CREATE SCHEMA test;

	DROP TABLE IF EXISTS test.table1;
	CREATE TABLE test.table1 (
		user_id SERIAL PRIMARY KEY,
		email VARCHAR ( 255 ) UNIQUE NOT NULL,
		last_update TIMESTAMP NOT NULL
	)
	
	`)

	if err != nil {
		t.Error("error creating schema", err)
		t.FailNow()
	}

	return db
}

func insertData(db *sql.DB, value string, t *testing.T) {
	_, err := db.Exec(`
		INSERT INTO test.table1 (email,last_update) 
		VALUES('mario@rossi.it', now())
	`)
	if err != nil {
		t.Error("Error insert temp table: ", err)
		t.FailNow()
	}
}
