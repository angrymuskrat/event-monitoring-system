package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	DB       string `json:"DB"`
	GRPCPort string `json:"GRPCPort"`
}

func readConfig(storage string) Config {
	jsonFile, err := os.Open(storage)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var conf Config

	json.Unmarshal(byteValue, &conf)

	return conf
}
