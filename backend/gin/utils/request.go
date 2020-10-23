package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpRequest(apiURL string, method string, headers map[string]string, data interface{}) (string, error) {
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
		return "", err
	}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode == 200 || response.StatusCode == 201 {
		return string(body), nil
	} else {
		return "", fmt.Errorf("%s", string(body))
	}
}
