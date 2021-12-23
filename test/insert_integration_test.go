package test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Ringloop/Mr-Plow/elastic"
	"github.com/Ringloop/Mr-Plow/movedata"
	"github.com/Ringloop/Mr-Plow/test_util"
	_ "github.com/lib/pq"
)

type insertIntegrationTest struct{}

// test case config scenario
func (*insertIntegrationTest) ReadConfig() ([]byte, error) {

	testComplexConfig := `
database: "postgres://user:pwd@localhost:5432/postgres?sslmode=disable"
queries:
  - query: "select * from test.table1 where last_update > $1"
    index: "out_index"
    updateDate: "last_update"
elastic:
  url: http://localhost:9200
`

	// Prepare data you want to return without reading from the file
	return []byte(testComplexConfig), nil
}

func TestInsertIntegration(t *testing.T) {
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
	originalLastDate, err := repo.FindLastUpdateOrEpochDate(conf.Queries[0].Index, conf.Queries[0].UpdateDate)
	if err != nil {
		t.Error("error in getting last date", err)
		t.FailNow()
	}

	//when (moving data to elastic)
	mover := movedata.New(db, conf, &conf.Queries[0])
	err = mover.MoveData()
	if err != nil {
		t.Error("error data moving", err)
		t.FailNow()
	}

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

	test_util.AssertEqual(t, len(response1.Hits.Hits), 1)
	test_util.AssertEqual(t, response1.Hits.Hits[0].Source.Email, "mario@rossi.it")
	test_util.AssertNotNull(t, response1.Hits.Hits[0].Source.LastUpdate)
	test_util.AssertNotNull(t, response1.Hits.Hits[0].Source.UserID)

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

	test_util.AssertEqual(t, len(response2.Hits.Hits), 2)
}
