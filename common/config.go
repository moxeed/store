package common

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	DataBaseDsn string
	ListenPort  int
	ZarinPal    struct {
		MerchantId          string
		BaseUrl             string
		CallBackUrl         string
		RequestRelativePath string
		VerifyRelativePath  string
		Authorization       string
		RedirectUrl         string
	}
	Job struct {
		FailedOrderRetryInMinutes  time.Duration
		OpenOrderRetryInMinutes    time.Duration
		OpenTerminalRetryInMinutes time.Duration
	}
	Store struct {
		BaseUrl       string
		CreateProduct string
		AddItem       string
		FlashBuy      string
	}
	Front struct {
		PaymentRedirect string
	}
}

var Configuration Config

func init() {
	configFile, err := os.Open("./config.json")

	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Configuration)
	if err != nil {
		panic(err)
	}
	_ = configFile.Close()
}
