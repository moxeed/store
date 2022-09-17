package common

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type CallState struct {
	IsOk        bool
	IsAmbiguous bool
}

func Post(url string, body interface{}, result interface{}) CallState {
	data, err := json.Marshal(body)
	Log(err)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	Log(err)

	respBody, err := io.ReadAll(resp.Body)
	Log(err)
	err = json.Unmarshal(respBody, result)
	Log(err)
	err = resp.Body.Close()
	Log(err)

	return CallState{
		IsOk:        err == nil && resp.StatusCode < 300,
		IsAmbiguous: err != nil || resp.StatusCode >= 500,
	}
}

func WriteIfNoError(c *echo.Context, err error, body interface{}) error {
	if err == nil {
		err = (*c).JSON(200, body)
	} else {
		err = (*c).JSON(400, Error{
			Status:  400,
			Message: err.Error(),
		})
	}
	return err
}
