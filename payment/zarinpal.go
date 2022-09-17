package payment

import (
	"fmt"
	"github.com/moxeed/store/common"
	"log"
	"strconv"
)

const (
	OK = iota + 100
	Repeated
)

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type TerminalRequest struct {
	MerchantId  string `json:"merchant_id"`
	Amount      int    `json:"amount"`
	CallbackUrl string `json:"callback_url"`
	Description string `json:"description"`
	Metadata    *struct {
		Mobile string `json:"mobile"`
		Email  string `json:"email"`
	} `json:"metadata,omitempty"`
}

type TerminalResponse struct {
	Data struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		Authority string `json:"authority"`
		FeeType   string `json:"fee_type"`
		Fee       int    `json:"fee"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

type VerifyRequest struct {
	MerchantId string `json:"merchant_id"`
	Amount     int    `json:"amount"`
	Authority  string `json:"authority"`
}

type VerifyResponse struct {
	Data struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		CardHash string `json:"card_hash"`
		CardPan  string `json:"card_pan"`
		RefId    int    `json:"ref_id"`
		FeeType  string `json:"fee_type"`
		Fee      int    `json:"fee"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

type VerifyResult struct {
	isRepeated    bool
	isOk          bool
	CardPan       string
	ReferenceCode string
	Fee           int
}

func openTerminal(amount int, description string) (string, error) {
	config := &common.Configuration.ZarinPal
	request := TerminalRequest{
		MerchantId:  config.MerchantId,
		Amount:      common.Abs(amount),
		CallbackUrl: config.CallBackUrl,
		Description: description,
	}
	result := TerminalResponse{}
	state := common.Post(config.BaseUrl+config.RequestRelativePath, request, &result)

	if state.IsOk {
		println(result.Data.Authority)
		return result.Data.Authority, nil
	}

	log.Println(result.Errors)
	return "", fmt.Errorf("خطا از سمت ارایه دهنده سرویس")
}

func verify(authority string, amount int) VerifyResult {
	config := &common.Configuration.ZarinPal
	request := VerifyRequest{
		MerchantId: config.MerchantId,
		Amount:     amount,
		Authority:  authority,
	}
	result := VerifyResponse{}
	state := common.Post(config.BaseUrl+config.RequestRelativePath, request, &result)

	if state.IsOk {
		return VerifyResult{
			isRepeated:    result.Data.Code == Repeated,
			isOk:          result.Data.Code == OK,
			CardPan:       result.Data.CardPan,
			ReferenceCode: strconv.Itoa(result.Data.RefId),
			Fee:           result.Data.Fee,
		}
	}

	return VerifyResult{
		isOk:       false,
		isRepeated: false,
	}
}
