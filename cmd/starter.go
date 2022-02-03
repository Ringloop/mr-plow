package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/Ringloop/mr-plow/pkg/scheduler"

	"github.com/Ringloop/mr-plow/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	configPath := flag.String("config", "./config.yml", "path of the configuration file")
	flag.Parse()

	ymlConfReader := config.Reader{FileName: *configPath}
	conf, err := config.ParseConfiguration(&ymlConfReader)
	if err != nil {
		log.Fatal("Cannot parse config file", err)
	}
	ConnectAndStart(conf)
}

func ConnectAndStart(conf *config.ImportConfig) {
	db, err := sql.Open("postgres", conf.Database)
	if err != nil {
		log.Printf("Failed to open a DB connection: %s", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("error in closing postgres connection: %s", err)
		}
	}(db)
	log.Println("Connected to postgres")

	finished := make(chan bool)
	for _, c := range conf.Queries {
		s := scheduler.NewScheduler()
		go s.MoveDataUntilExit(conf, db, &c, finished)
	}

	for i := 0; i < len(conf.Queries); i++ {
		<-finished
	}
}
