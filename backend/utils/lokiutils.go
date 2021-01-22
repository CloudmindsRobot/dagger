package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/runtime"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

func Exist(eles []interface{}, ele interface{}) bool {

	exist := false
	for _, element := range eles {
		if element.(map[string]interface{})["value"].(string) == ele.(string) {
			exist = true
			return exist
		}
	}
	return exist
}

func GetLogLevelExpr(level string) string {
	switch level {
	case "info":
		return `(([[]I[]])|([[【= \t][iI](?i)nfo[]】 \t]))|(INFO)`
	case "debug":
		return `(([[]D[]])|([[【= \t][dD](?i)ebug[]】 \t]))|(DEBUG)`
	case "warn":
		return `(([[]W[]])|([[【= \t][wW](?i)arn((?i)(ing))?[]】 \t]))|(WARN)|(WARNING)`
	case "error":
		return `(([[]E[]])|([[【= \t][eE](?i)rror[]】 \t]))|(ERROR)`
	default:
		return `(([[]I[]])|([[【= \t][iI](?i)nfo[]】 \t]))|(INFO)|(([[]D[]])|([[【= \t][dD](?i)ebug[]】 \t]))|(DEBUG)|(([[]W[]])|([[【= \t][wW](?i)arn((?i)(ing))?[]】 \t]))|(WARN)|(WARNING)|(([[]E[]])|([[【= \t][eE](?i)rror[]】 \t]))|(ERROR)`
	}
}

func LogLevel(message string) string {
	reg := regexp.MustCompile(GetLogLevelExpr("info"))
	if reg.MatchString(message) {
		return "info"
	}
	reg = regexp.MustCompile(GetLogLevelExpr("debug"))
	if reg.MatchString(message) {
		return "debug"
	}
	reg = regexp.MustCompile(GetLogLevelExpr("warn"))
	if reg.MatchString(message) {
		return "warn"
	}
	reg = regexp.MustCompile(GetLogLevelExpr("error"))
	if reg.MatchString(message) {
		return "error"
	}
	return "unknown"
}

func GenerateLevelRegex(level string) string {

	all := GetLogLevelExpr("all")
	if strings.ToLower(level) == "unknown" {
		return fmt.Sprintf("!~ `%s`", all)
	}
	levelArray := []string{}
	levelUnknownExist := false
	levels := strings.Split(level, ",")
	levelExpr := ""
	for _, l := range levels {
		switch strings.ToLower(l) {
		case "info":
			levelArray = append(levelArray, GetLogLevelExpr("info"))
			all = strings.ReplaceAll(all, GetLogLevelExpr("info"), "")
		case "debug":
			levelArray = append(levelArray, GetLogLevelExpr("debug"))
			all = strings.ReplaceAll(all, GetLogLevelExpr("debug"), "")
		case "warn":
			levelArray = append(levelArray, GetLogLevelExpr("warn"))
			all = strings.ReplaceAll(all, GetLogLevelExpr("warn"), "")
		case "error":
			levelArray = append(levelArray, GetLogLevelExpr("error"))
			all = strings.ReplaceAll(all, GetLogLevelExpr("error"), "")
		case "unknown":
			levelUnknownExist = true
		}
	}
	if len(levelArray) > 0 {
		levelExpr = fmt.Sprintf("|~ `%s`", strings.Join(levelArray, "|"))
	}
	if levelUnknownExist {
		all = strings.Trim(all, "|")
		if all == "" {
			return ""
		}
		levelExpr = fmt.Sprintf("!~ `%s`", all)
	}

	return levelExpr
}

func SplitDateTime(start string, end string, limit int) ([]int, int) {
	// 切20份  1586330540000 000000
	startIndex, _ := strconv.Atoi(start[0:13])
	endIndex, _ := strconv.Atoi(end[0:13])

	step := (endIndex - startIndex) / limit
	index := 0
	splitDateTimeArray := []int{}
	for {
		if index < limit {
			splitDateTimeArray = append(splitDateTimeArray, startIndex+step*index)
			index++
			continue
		}
		break
	}

	return splitDateTimeArray, step
}

func InitSplitDateTime(limit int) map[string][]int {
	chartData := make(map[string][]int)
	chartData["info"] = make([]int, limit, limit)
	chartData["debug"] = make([]int, limit, limit)
	chartData["warn"] = make([]int, limit, limit)
	chartData["error"] = make([]int, limit, limit)
	chartData["unknown"] = make([]int, limit, limit)
	return chartData
}

func TimeInPart(splitDateTime []int, timestamp string, step int) int {
	timestampIndex, _ := strconv.Atoi(timestamp[0:13])
	stepSum := (timestampIndex - splitDateTime[0]) / step

	return stepSum
}

func ShellHighlightShow(message string) string {
	highlight := message
	regExceptFilter, _ := regexp.Compile(`[[]\d*;?\d*m.*?[[]\d*;?\d*m`)
	if regExceptFilter.MatchString(message) {
		regInnerFilter, _ := regexp.Compile(`[[](\d*);?(\d*)m(.*)[[]\d*;?\d*m`)
		highlight = regExceptFilter.ReplaceAllStringFunc(highlight, func(item string) string {
			strs := regInnerFilter.FindStringSubmatch(item)
			if len(strs) == 4 {
				color := GetShellColor(strs[1])
				if color == "" {
					color = GetShellColor(strs[2])
				}
				if strs[1] == "1" || strs[1] == "" {
					return fmt.Sprintf(`<b style="color: %s !important;">%s</b>`, color, strs[3])
				} else {
					return fmt.Sprintf(`<b style="background-color: %s !important;color: slategray;">%s</b>`, color, strs[3])
				}
			}
			return item
		})
	}
	return highlight
}

func RegexHighlightShow(message string, filter string) string {
	highlight := message
	if filter != "" {
		regFilter, _ := regexp.Compile(filter)
		highlight = regFilter.ReplaceAllStringFunc(highlight, func(item string) string {
			return fmt.Sprintf(`<b style="color: #fb8c00 !important;">%s</b>`, item)
		})

	}
	return highlight
}

func GetShellColor(colorNo string) string {
	switch colorNo {
	case "40", "30":
		return "black"
	case "41", "31":
		return "red"
	case "42", "32":
		return "lightgreen"
	case "43", "33":
		return "yellow"
	case "44", "34":
		return "blue"
	case "45", "35":
		return "magenta"
	case "46", "36":
		return "cyan"
	case "47", "37":
		return "wheat"
	}
	return ""
}

func GetExpr(label string, value string) string {
	queryExpr := ""
	if label != "" && value != "" {
		queryExpr = fmt.Sprintf("%s=\"%s\"", label, value)
	}
	return queryExpr
}

func GetPodExpr(pod string) string {
	queryExpr := ""
	if pod != "" {
		queryExpr = fmt.Sprintf("k8s_pod_name=~\"%s\"", pod)
	}
	return queryExpr
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GenerateToken(uid int, username string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      uid,
		"username": username,
		"exp":      expire,
	})
	return token.SignedString([]byte("dagger-backend-secret"))
}

func LdapCheck(username string, password string) (bool, *models.User) {
	ldapHost, _ := runtime.Cfg.GetValue("ldap", "ldap_host")
	ldapPort, _ := runtime.Cfg.Int("ldap", "ldap_port")

	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))

	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("ldap connect error: %s", err))
		return false, nil
	}
	defer conn.Close()

	err = conn.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("ldap start error: %s", err))
		return false, nil
	}

	ldapBindUsername, _ := runtime.Cfg.GetValue("ldap", "ldap_bind_username")
	ldapBindPassword, _ := runtime.Cfg.GetValue("ldap", "ldap_bind_password")
	err = conn.Bind(ldapBindUsername, ldapBindPassword)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("ldap bind error: %s", err))
		return false, nil
	}

	base, _ := runtime.Cfg.GetValue("ldap", "ldap_base_dn")
	usernameKey, _ := runtime.Cfg.GetValue("ldap", "ldap_username_key")
	mailKey, _ := runtime.Cfg.GetValue("ldap", "ldap_mail_key")
	sql := ldap.NewSearchRequest(base,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(%s=%s)", usernameKey, username),
		[]string{fmt.Sprintf("%s", usernameKey), fmt.Sprintf("%s", mailKey)},
		nil)

	var cur *ldap.SearchResult
	if cur, err = conn.Search(sql); err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("ldap server search failed.: %s", err))
		return false, nil
	}

	if len(cur.Entries) == 0 {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("ldap not found user."))
		return false, nil
	}

	user := cur.Entries[0]
	err = conn.Bind(user.DN, password)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("check user error."))
		return false, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("%s", err))
		return false, nil
	}

	var userModel = models.User{}
	result := databases.DB.Model(&models.User{}).Where("username = ?", username).First(&userModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userModel = models.User{
			Username:    username,
			Password:    string(hash),
			IsActive:    true,
			IsSuperuser: false,
			Email:       user.GetAttributeValue("mail"),
			CreateAt:    time.Now().UTC(),
			LastLoginAt: time.Now().UTC(),
		}
		databases.DB.Create(&userModel)
	}

	return true, &userModel
}

func Md5(data string) string {
	keyByte := []byte(data)
	hashKey := md5.Sum(keyByte)
	md5Str := fmt.Sprintf("%x", hashKey)
	return md5Str
}

func StructLables(lables []models.LogLabel, name string) string {
	var elementArray []string

	var keys []string
	labelMap := make(map[string]string)
	for _, label := range lables {
		keys = append(keys, label.Key)
		labelMap[label.Key] = label.Value
	}
	keys = append(keys, "name")
	labelMap["name"] = name

	sort.Strings(keys)
	for _, k := range keys {
		elementArray = append(elementArray, fmt.Sprintf("%s=%s", k, labelMap[k]))
	}
	return strings.Join(elementArray, "&")
}

func ToYAML(s models.RuleYAML) (*bytes.Buffer, error) {
	d, err := yaml.Marshal(s)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(d)
	return b, nil
}

func GenerateYAML(rule models.LogRule) (string, error) {
	labelsMap := make(map[string]interface{})
	for _, label := range rule.Labels {
		labelsMap[label.Key] = label.Value
	}
	alertYAML := models.AlertYAML{
		Alert:  rule.Name,
		For:    0,
		Expr:   rule.LogQL,
		Labels: labelsMap,
		Annotations: map[string]interface{}{
			"description": rule.Description,
			"summary":     rule.Summary,
			"key":         rule.Key,
		},
	}

	file := models.RuleYAML{
		Name:  rule.ID,
		Rules: []models.AlertYAML{alertYAML},
	}
	b, err := ToYAML(file)
	if err != nil {
		return "", err
	}
	content := b.String()
	return content, nil
}

func String2Time(timeStr string, timeZone string) time.Time {
	local, _ := time.LoadLocation(timeZone)
	t, _ := time.ParseInLocation("2006-01-02T15:04:05Z", timeStr, local)
	return t
}

func TimeDateValueFormatter(v interface{}) string {
	if typed, isTyped := v.(float64); isTyped {
		return time.Unix(0, int64(typed)).Format("15:04:05")
	}
	return ""
}

func SplitDateTimeForMatrix(start string, end string) ([]int64, []string, int) {
	startIndex, _ := strconv.ParseInt(start[0:10], 10, 64)
	endIndex, _ := strconv.ParseInt(end[0:10], 10, 64)

	m := endIndex - startIndex
	interval := 1
	if m > 1000 && m < 10000 {
		interval = 10
	} else if m > 10000 && m < 100000 {
		interval = 100
	} else if m > 100000 && m < 1000000 {
		interval = 1000
	} else if m > 1000000 && m < 10000000 {
		interval = 10000
	} else if m > 10000000 && m < 100000000 {
		interval = 100000
	} else if m > 100000000 && m < 1000000000 {
		interval = 1000000
	} else if m > 1000000000 && m < 10000000000 {
		interval = 10000000
	}

	var index int64
	index = 0
	splitDateTimeArray := []int64{}
	splitValueArray := []string{}
	for {
		if index <= (endIndex - startIndex + 1) {
			if int(index)%interval == 0 {
				splitDateTimeArray = append(splitDateTimeArray, (startIndex+index)*1000)
			}
			splitValueArray = append(splitValueArray, "0")
			index++
			continue
		}
		break
	}

	return splitDateTimeArray, splitValueArray, interval
}
