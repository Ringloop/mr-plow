package config

import (
	"io/ioutil"
	"log"
)

// Defining an interface so that functionality of 'readConfig()' can be mocked
type IReader interface {
	ReadConfig() ([]byte, error)
}

type Reader struct {
	FileName string
}

// 'reader' implementing the Interface
// Function to read from actual file
func (r *Reader) ReadConfig() ([]byte, error) {
	configFile, err := ioutil.ReadFile(r.FileName)

	if err != nil {
		log.Fatal(err)
	}
	return configFile, err
}
