package test

import (
	"database/sql"
	"testing"

	"github.com/Ringloop/mr-plow/internal/config"
)

func initConfigIntegrationTest(t *testing.T, testReader config.IReader) *config.ImportConfig {
	conf, err := config.ParseConfiguration(testReader)
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
		email VARCHAR ( 255 ) NOT NULL,
		last_update TIMESTAMP NOT NULL
	)
	
	`)

	if err != nil {
		t.Error("error creating schema", err)
		t.FailNow()
	}

	return db
}

func insertData(db *sql.DB, email string, t *testing.T) {
	_, err := db.Exec(`
		INSERT INTO test.table1 (email,last_update) 
		VALUES($1, now())
	`, email)
	if err != nil {
		t.Error("Error insert temp table: ", err)
		t.FailNow()
	}
}
