package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func Post(url string, body interface{}, result interface{}) bool {
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

	return err == nil && resp.StatusCode < 300
}
