package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func DoRequest(method, url string, body map[string]interface{}, response any, setReq func(req *http.Request) error) (*http.Response, error) {
	var requestBody io.Reader
	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer([]byte(bodyJSON))
	}
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, err
	}

	if setReq != nil {
		err = setReq(req)
		if err != nil {
			return nil, err
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return res, nil
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resData, &response)
	return res, err
}
