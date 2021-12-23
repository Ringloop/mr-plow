package main

import (
	"database/sql"
	"github.com/Ringloop/Mr-Plow/scheduler"
	"log"

	"github.com/Ringloop/Mr-Plow/config"
	_ "github.com/lib/pq"
)

func main() {
	ymlConfReader := config.Reader{FileName: "config.yml"} //TODO parse commandline yml path, now assuming is in current dir
	conf, err := config.ParseConfiguration(&ymlConfReader)
	if err != nil {
		log.Fatal("Cannot parse config file", err)
	}
	ConnectAndStart(conf)
}

func ConnectAndStart(conf *config.ImportConfig) {
	db, err := sql.Open("postgres", conf.Database)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("error in closing postgres connection", err)
		}
	}(db)
	log.Println("Connected to postgres")
	for _, c := range conf.Queries {
		go scheduler.MoveDataUntilExit(conf, db, &c)
	}
}
