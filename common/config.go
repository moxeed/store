package common

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DataBaseDsn string
	ListenPort  int
}

var Configuration Config

func init() {
	configFile, err := os.Open("./config.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&Configuration)
}
