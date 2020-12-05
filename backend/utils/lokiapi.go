package utils

import (
	"dagger/backend/runtime"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func QueryRange(query string, limit int, start string, end string, direction string) map[string]interface{} {
	params := make(map[string]interface{})
	url := "/loki/api/v1/query_range"
	queryURL := fmt.Sprintf("%s%s?query=%s&start=%s&end=%s&limit=%d&direction=%s", runtime.LokiServer, url, query, start, end, limit, direction)
	Log4Zap(zap.InfoLevel).Info(fmt.Sprintf("loki api query url: %s", queryURL))
	repeat := 0
	var data string
	var err error
	var code int
	for {
		if repeat < 5 {
			data, code, err = HttpRequest(queryURL, "GET", nil, params)
			if err != nil {
				repeat++
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		} else {
			return nil
		}
	}

	if code != 200 && code != 201 {
		Log4Zap(zap.ErrorLevel).Error(data)
		return nil
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonRes)
	if err != nil {
		Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("Unmarshal loki query range response error %s", err))
		return nil
	}

	if _, ok := jsonRes["data"]; ok {
		return jsonRes["data"].(map[string]interface{})
	}

	return nil
}

func Labels(start string, end string) []interface{} {
	url := "/loki/api/v1/labels"
	queryURL := fmt.Sprintf("%s%s?start=%s&end=%s", runtime.LokiServer, url, start, end)
	repeat := 0
	var data string
	var err error
	var code int
	for {
		if repeat < 5 {
			data, code, err = HttpRequest(queryURL, "GET", nil, nil)
			if err != nil {
				repeat++
				Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("get loki labels error %s", err))
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		} else {
			return nil
		}
	}

	if code != 200 && code != 201 {
		Log4Zap(zap.ErrorLevel).Error(data)
		return nil
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonRes)
	if err != nil {
		Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("Unmarshal loki labels response error %s", err))
		return nil
	}

	if _, ok := jsonRes["data"]; ok {
		values := jsonRes["data"].([]interface{})
		vals := []interface{}{}
		for _, value := range values {
			if value.(string) != "__name__" {
				vals = append(vals, value)
			}
		}
		return vals
	}

	return nil
}

func LabelValues(label string, start string, end string) []interface{} {
	queryURL := fmt.Sprintf("%s/loki/api/v1/label/%s/values?start=%s&end=%s", runtime.LokiServer, label, start, end)
	repeat := 0
	var data string
	var err error
	var code int
	for {
		if repeat < 5 {
			data, code, err = HttpRequest(queryURL, "GET", nil, nil)
			if err != nil {
				repeat++
				Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("get loki label values error %s", err))
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		} else {
			return nil
		}
	}

	if code != 200 && code != 201 {
		Log4Zap(zap.ErrorLevel).Error(data)
		return nil
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonRes)
	if err != nil {
		Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("Unmarshal loki label values response error %s", err))
		return nil
	}

	if _, ok := jsonRes["data"]; ok {
		values := jsonRes["data"].([]interface{})
		return values
	}

	return nil
}
