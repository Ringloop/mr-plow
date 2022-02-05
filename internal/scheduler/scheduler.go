package scheduler

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ringloop/mr-plow/internal/config"
	"github.com/Ringloop/mr-plow/internal/movedata"
)

type Scheduler struct {
	Done chan os.Signal
}

func NewScheduler() Scheduler {
	s := Scheduler{make(chan os.Signal)}
	signal.Notify(s.Done, os.Interrupt, syscall.SIGTERM)
	return s
}

func (s *Scheduler) MoveDataUntilExit(conf *config.ImportConfig, db *sql.DB, query *config.QueryModel, finished chan bool) {
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
		case <-s.Done:
			log.Println("stopping query execution, bye...")
			finished <- true
			return
		}
	}
}
