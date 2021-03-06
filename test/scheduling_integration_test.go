package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Ringloop/mr-plow/internal/scheduler"
	"github.com/stretchr/testify/require"

	"github.com/Ringloop/mr-plow/internal/elastic"
	_ "github.com/lib/pq"
)

func TestSchedulingIntegration(t *testing.T) {
	//given (some data on sql db)
	conf := initConfigIntegrationTest(t, &insertIntegrationTest{})
	db := initSqlDB(t, conf)
	defer db.Close()
	repo, err := elastic.NewDefaultClient()
	if err != nil {
		t.Error("error in creating elastic connection", err)
		t.FailNow()
	}
	repo.Delete(conf.Queries[0].Index)

	insertData(db, "mario@rossi.it", t)

	//when (starting the scheduler)
	finished := make(chan bool)
	s := scheduler.NewScheduler()
	go s.MoveDataUntilExit(conf, db, &conf.Queries[0], finished)
	time.Sleep(2 * time.Second)

	indexContent1, err := repo.FindIndexContent("out_index", "last_update")
	defer (*indexContent1).Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var response1 ElasticTestResponse
	if err := json.NewDecoder(*indexContent1).Decode(&response1); err != nil {
		t.Error(err)
		t.FailNow()
	}

	require.Equal(t, len(response1.Hits.Hits), 1)
	require.Equal(t, response1.Hits.Hits[0].Source.Email, "mario@rossi.it")
	require.NotNil(t, response1.Hits.Hits[0].Source.LastUpdate)
	require.NotNil(t, response1.Hits.Hits[0].Source.UserID)

	//and when (inserting new data)
	insertData(db, "mario@rossi.it", t)

	// and then (the data is moved)
	time.Sleep(2 * time.Second)
	indexContent2, err := repo.FindIndexContent("out_index", "last_update")
	defer (*indexContent2).Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var response2 ElasticTestResponse
	if err := json.NewDecoder(*indexContent2).Decode(&response2); err != nil {
		t.Error(err)
		t.FailNow()
	}

	require.Equal(t, len(response2.Hits.Hits), 2)

	s.Done <- FakeExitSignal{}

	//and when (inserting new data again)
	time.Sleep(2 * time.Second)
	insertData(db, "mario@rossi.it", t)

	//and then, nothing new is extracted
	indexContent3, err := repo.FindIndexContent("out_index", "last_update")
	defer (*indexContent3).Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var response3 ElasticTestResponse
	if err := json.NewDecoder(*indexContent3).Decode(&response3); err != nil {
		t.Error(err)
		t.FailNow()
	}

	require.Equal(t, <-finished, true)
	require.Equal(t, len(response3.Hits.Hits), len(response2.Hits.Hits))
}
