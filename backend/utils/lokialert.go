package utils

import (
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/runtime"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func PushToAlertPlatform(body interface{}) bool {
	alertPlatformHost, _ := runtime.Cfg.GetValue("dayu", "dayu_alert_engine")
	url := fmt.Sprintf("%s/api/v1/external", alertPlatformHost)
	repeat := 0
	var err error
	for {
		if repeat < 5 {
			_, err = HttpRequest(url, "POST", nil, body, "json")
			if err != nil {
				repeat++
				Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("push loki alert error %s", err))
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		} else {
			return false
		}
	}

	return true
}

func CacheRule() error {
	var rules []models.LogRule
	databases.DB.Preload("User").Preload("Groups").Preload("Groups.LogUserGroup").Preload("Groups.LogUserGroup.Users").Preload("Groups.LogUserGroup.Users.User").Model(&models.LogRule{}).Find(&rules)

	ruleMap := make(map[string]models.LogRule)
	for _, rule := range rules {
		ruleMap[rule.Key] = rule
	}

	err := databases.GC.Set("rule-map", ruleMap)
	return err
}
