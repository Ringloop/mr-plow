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

func MoveDataUntilExit(conf *config.ImportConfig, db *sql.DB, query *config.QueryModel) {
	ticker := time.NewTicker(time.Duration(conf.PollingSeconds * 1000000000))
	defer ticker.Stop()
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	mover := movedata.New(db, conf, query)
	for {
		select {
		case <-ticker.C:
			moveErr := mover.MoveData()
			if moveErr != nil {
				log.Fatal("error executing query", moveErr)
			}
		case <-done:
			log.Println("stopping query execution, bye...")
			return
		}
	}
}
