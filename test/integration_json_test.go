package test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"dariobalinzo.com/elastic/v2/config"
	"dariobalinzo.com/elastic/v2/elastic"
	"dariobalinzo.com/elastic/v2/movedata"
	_ "github.com/lib/pq"
)

func initSqlDB_local(t *testing.T, conf *config.ImportConfig) *sql.DB {
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
		additional_info JSON,
		last_update TIMESTAMP NOT NULL
	)
	
	`)

	if err != nil {
		t.Error("error creating schema", err)
		t.FailNow()
	}

	return db
}

func TestIntegrationWithJSON(t *testing.T) {
	//given (some data on sql db)
	conf := initConfigIntegrationTestWithJson(t)
	db := initSqlDB_local(t, conf)
	defer db.Close()
	elastic.Delete(conf.Queries[0].Index)

	email := "mario@rossi.it"
	json := `
{
	"str_col": "String Data",
    "int_col": 4237,
    "bool_col": true,
    "float_col": 48.94065780742467
}`
	insertDataWithJSON(db, email, json, t)
	originalLastDate, err := elastic.FindLastUpdateOrEpochDate(conf.Queries[0].Index)
	if err != nil {
		t.Error("error in getting last date", err)
		t.FailNow()
	}

	//when (moving data to elastic)
	err = movedata.MoveData(db, conf.Queries[0])
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

	log.Println("JSON_TEST: original date", originalLastDate)
	log.Println("JSON_TEST: date after import", lastImportedDate)

	if !lastImportedDate.After(*originalLastDate) {
		t.Error("error date not incremented!")
		t.FailNow()
	}

}

type readerIntegrationTestWithJson struct{}

// 'readerTest' implementing the Interface
func (*readerIntegrationTestWithJson) ReadConfig() ([]byte, error) {

	configIntegrationWithJson := `
database: "postgres://user:pwd@localhost:5432/postgres?sslmode=disable"
queries:
  - index: "out_index"
    query: "select * from test.table1 where last_update > $1"
    JSONFields:
      - fieldName: additional_info
        attributes:
          - attributeName: str_col
            attributeType: string
          - attributeName: int_col
            attributeType: integer
          - attributeName: bool_col
            attributeType: boolean
          - attributeName: float_col
            attributeType: float
`

	// Prepare data you want to return without reading from the file
	return []byte(configIntegrationWithJson), nil
}

func initConfigIntegrationTestWithJson(t *testing.T) *config.ImportConfig {
	testReader := readerIntegrationTestWithJson{}
	conf, err := config.ParseConfiguration(&testReader)
	if err != nil {
		t.Error("error reading conf", err)
		t.FailNow()
	}
	return conf
}

func insertDataWithJSON(db *sql.DB, email string, json string, t *testing.T) {
	sql_statement := fmt.Sprintf(`
	INSERT INTO test.table1 (email, additional_info, last_update)
	VALUES ('%s', '%s', now());	
	`, email, json)
	_, err := db.Exec(sql_statement)
	if err != nil {
		t.Error("Error insert temp table: ", err)
		t.FailNow()
	}
}
