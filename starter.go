package main

import (
	"database/sql"
	"log"

	"dariobalinzo.com/elastic/v2/config"
	"dariobalinzo.com/elastic/v2/movedata"

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
	db, err := sql.Open("postgres", conf.SqlConfig)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	log.Println("Connected to postgres", err)

	for _, c := range conf.Queries {
		go func(c config.QueryModel) {
			moveErr := movedata.MoveData(db, c.Query, c.Index)
			if moveErr != nil {
				log.Fatal("error execurting query", err)
			}
		}(c)

	}

}
