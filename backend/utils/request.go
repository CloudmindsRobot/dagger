package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

func HttpRequest(apiURL string, method string, headers map[string]string, data interface{}, contentType string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
	}

	var err error
	var req *http.Request

	if contentType == "json" {
		requestData, _ := json.Marshal(data)
		req, err = http.NewRequest(method, apiURL, bytes.NewBuffer(requestData))
		req.Header.Set("Content-Type", "application/json")
	} else if contentType == "yaml" {
		req, err = http.NewRequest(method, apiURL, bytes.NewBufferString(data.(string)))
		req.Header.Set("Content-Type", "application/yaml")
	}

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
	if response.StatusCode == 200 || response.StatusCode == 201 || response.StatusCode == 202 || response.StatusCode == 204 {
		return string(body), nil
	} else if response.StatusCode == 429 {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("触发限流"))
		return "", nil
	} else {
		return "", fmt.Errorf("%s", string(body))
	}
}
