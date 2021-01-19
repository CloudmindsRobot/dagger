package utils

import (
	"dagger/backend/runtime"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func QueryRange(query string, limit int, start string, end string, direction string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	url := "/loki/api/v1/query_range"
	queryURL := fmt.Sprintf("%s%s?query=%s&start=%s&end=%s&limit=%d&direction=%s", runtime.LokiServer, url, query, start, end, limit, direction)
	Log4Zap(zap.InfoLevel).Info(fmt.Sprintf("loki api query url: %s", queryURL))
	repeat := 0
	var data string
	var err error
	for {
		if repeat < 5 {
			data, err = HttpRequest(queryURL, "GET", nil, params, "json")
			if err != nil {
				repeat++
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		} else {
			return nil, err
		}
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonRes)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("Unmarshal loki query range response error %s", err))
		return nil, err
	}

	if _, ok := jsonRes["data"]; ok {
		return jsonRes["data"].(map[string]interface{}), nil
	}

	return nil, fmt.Errorf("unknown error")
}

func Labels(start string, end string) []interface{} {
	url := "/loki/api/v1/labels"
	queryURL := fmt.Sprintf("%s%s?start=%s&end=%s", runtime.LokiServer, url, start, end)
	repeat := 0
	var data string
	var err error
	for {
		if repeat < 5 {
			data, err = HttpRequest(queryURL, "GET", nil, nil, "json")
			if err != nil {
				repeat++
				Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("get loki labels error %s", err))
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
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("Unmarshal loki labels response error %s", err))
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
	for {
		if repeat < 5 {
			data, err = HttpRequest(queryURL, "GET", nil, nil, "json")
			if err != nil {
				repeat++
				Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("get loki label values error %s", err))
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
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("Unmarshal loki label values response error %s", err))
		return nil
	}

	if _, ok := jsonRes["data"]; ok {
		values := jsonRes["data"].([]interface{})
		return values
	}

	return nil
}

func CreateOrUpdateRuleGroup(namespace string, yaml string) (bool, error) {
	url := fmt.Sprintf("/loki/api/v1/rules/%s", namespace)
	var data string
	var err error

	data, err = HttpRequest(fmt.Sprintf("%s%s", runtime.LokiServer, url), "POST", nil, yaml, "yaml")
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("create or update loki rule group error %s", err))
		return false, err
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonRes)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("Unmarshal loki rule response error %s", err))
		return false, err
	}

	if _, ok := jsonRes["status"]; ok {
		values := jsonRes["status"].(string)
		return (values == "success"), nil
	}

	return false, err
}

func DeleteRuleGroup(namespace string, groupName string) (bool, error) {
	url := fmt.Sprintf("/loki/api/v1/rules/%s/%s", namespace, groupName)
	var data string
	var err error

	data, err = HttpRequest(fmt.Sprintf("%s%s", runtime.LokiServer, url), "DELETE", nil, "", "yaml")
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("delete loki rule group error %s", err))
		return false, err
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonRes)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("Unmarshal loki rule response error %s", err))
		return false, err
	}

	if _, ok := jsonRes["status"]; ok {
		values := jsonRes["status"].(string)
		return (values == "success"), nil
	}

	return false, nil
}

func LoadRules(namespace string) string {
	url := fmt.Sprintf("/loki/api/v1/rules/%s", namespace)
	repeat := 0
	var data string
	var err error
	for {
		if repeat < 5 {
			data, err = HttpRequest(fmt.Sprintf("%s%s", runtime.LokiServer, url), "GET", nil, nil, "json")
			if err != nil {
				repeat++
				Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("get loki rule group error %s", err))
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		} else {
			return ""
		}
	}

	return data
}
