package test

import (
	"encoding/json"
	"log"
	"sync"
	"testing"

	"github.com/Ringloop/mr-plow/elastic"
	"github.com/Ringloop/mr-plow/movedata"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type insertTypedIntegrationTest struct{}

// test case config scenario
func (*insertTypedIntegrationTest) ReadConfig() ([]byte, error) {

	return ([]byte(`
pollingSeconds: 1
database: "postgres://user:pwd@localhost:5432/postgres?sslmode=disable"
queries:
  - query: "select * from test.table1 where last_update > $1"
    index: "out_index"
    updateDate: "last_update"
    fields:
      - name: email
        type: String
      - name: user_id
        type: Integer
elastic:
  url: http://localhost:9200
`), nil)
}

func TestInsertTypedIntegration(t *testing.T) {
	//given (some data on sql db)
	conf := initConfigIntegrationTest(t, &insertTypedIntegrationTest{})
	db := initSqlDB(t, conf)
	defer db.Close()
	repo, err := elastic.NewDefaultClient()
	if err != nil {
		t.Error("error in creating elastic connection", err)
		t.FailNow()
	}
	repo.Delete(conf.Queries[0].Index)

	insertData(db, "mario@rossi.it", t)
	originalLastDate, err := repo.FindLastUpdateOrEpochDate(conf.Queries[0].Index, conf.Queries[0].UpdateDate)
	if err != nil {
		t.Error("error in getting last date", err)
		t.FailNow()
	}

	//when (moving data to elastic
	var allMovesDone sync.WaitGroup
	allMovesDone.Add(5)
	mover := movedata.New(db, conf, &conf.Queries[0])
	doMove := func() {
		defer allMovesDone.Done()
		errRoutine := mover.MoveData()
		if errRoutine != nil {
			t.Error("error data moving", err)
			t.FailNow()
		}
	}

	//(testing also long-running execution by executing the function as separate go routine))
	for i := 0; i < 5; i++ {
		go doMove()
	}
	allMovesDone.Wait()

	//then (last date on elastic should be updated)
	lastImportedDate, err := repo.FindLastUpdateOrEpochDate(conf.Queries[0].Index, conf.Queries[0].UpdateDate)
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
	err = mover.MoveData()
	if err != nil {
		t.Error("error data moving", err)
		t.FailNow()
	}

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
}
