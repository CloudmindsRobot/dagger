package utils

import (
	"dagger/backend/gin/runtime"
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
	for {
		if repeat < 5 {

			data, err = HttpRequest(queryURL, "GET", nil, params)
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

func Labels() []interface{} {
	url := "/loki/api/v1/labels"
	repeat := 0
	var data string
	var err error
	for {
		if repeat < 5 {
			data, err = HttpRequest(fmt.Sprintf("%s%s", runtime.LokiServer, url), "GET", nil, nil)
			if err != nil {
				repeat++
				Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("get loki labels error %s", err))
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		} else {
			return []interface{}{}
		}
	}
	var jsonRes map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonRes)
	if err != nil {
		Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("Unmarshal loki labels response error %s", err))
		return []interface{}{}
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

	return []interface{}{}
}

func LabelValues(label string) []interface{} {
	params := make(map[string]interface{})
	h, _ := time.ParseDuration("-1h")
	t := time.Now().Add(24 * h * 7).Unix()
	url := fmt.Sprintf("/loki/api/v1/label/%s/values?start=%s", label, fmt.Sprintf("%d000000000", t))
	repeat := 0
	var data string
	var err error
	for {
		if repeat < 5 {
			data, err = HttpRequest(fmt.Sprintf("%s%s", runtime.LokiServer, url), "GET", nil, params)
			if err != nil {
				repeat++
				Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("get loki label values error %s", err))
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		} else {
			return []interface{}{}
		}
	}
	var jsonRes map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonRes)
	if err != nil {
		Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("Unmarshal loki label values response error %s", err))
		return []interface{}{}
	}

	if _, ok := jsonRes["data"]; ok {
		values := jsonRes["data"].([]interface{})
		return values
	}

	return []interface{}{}
}
