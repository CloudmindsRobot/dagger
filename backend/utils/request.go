package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpRequest(apiURL string, method string, headers map[string]string, data interface{}) (string, int, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
	}

	requestData, _ := json.Marshal(data)
	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(requestData))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if err != nil {
		return "", 500, err
	}

	response, err := client.Do(req)
	if err != nil {
		return "", 500, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", 500, err
	}

	defer response.Body.Close()
	return string(body), response.StatusCode, nil
}
