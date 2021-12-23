package test

import "time"

//generated thanks to https://mholt.github.io/json-to-go/

type ElasticTestResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Index  string      `json:"_index"`
			Type   string      `json:"_type"`
			ID     string      `json:"_id"`
			Score  interface{} `json:"_score"`
			Source struct {
				Email      string    `json:"email"`
				LastUpdate time.Time `json:"last_update"`
				UserID     int       `json:"user_id"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type FakeExitSignal struct{}

func (FakeExitSignal) String() string {
	return "Fake exit"
}

func (FakeExitSignal) Signal() {}
