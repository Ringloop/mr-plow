package scheduler

import (
	"database/sql"
	"github.com/Ringloop/Mr-Plow/config"
	"github.com/Ringloop/Mr-Plow/movedata"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Done = make(chan os.Signal)

func init() {
	signal.Notify(Done, os.Interrupt, syscall.SIGTERM)
}

func MoveDataUntilExit(conf *config.ImportConfig, db *sql.DB, query *config.QueryModel, finished chan bool) {
	ticker := time.NewTicker(time.Duration(conf.PollingSeconds * 1000000000))
	defer ticker.Stop()

	mover := movedata.New(db, conf, query)
	for {
		select {
		case <-ticker.C:
			moveErr := mover.MoveData()
			if moveErr != nil {
				log.Printf("error executing query %s", moveErr)
			}
		case <-Done:
			log.Println("stopping query execution, bye...")
			finished <- true
			return
		}
	}
}
