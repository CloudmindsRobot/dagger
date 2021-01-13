package controllers

import (
	"dagger/backend/models"
	"dagger/backend/utils"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//
// @Summary Used to do a query over a range of time and accepts the following query parameters in the URL
// @Description default limit 2000
// @Accept  json
// @Produce  json
// @Param   start path string true "The start time for the query as a nanosecond Unix epoch"
// @Param   end path string true "The end time for the query as a nanosecond Unix epoch"
// @Param   all path string false "The new query to all results"
// @Param   dsc path string true "The order to all results"
// @Param   filter path string false "The filter condition"
// @Param   level path string false "The log level"
// @Param   limit path string false "The max number of entries to return"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/list/ [get]
func LokiList(c *gin.Context) {
	filters := c.QueryArray("filters[]")
	level := c.DefaultQuery("level", "")
	pod := c.DefaultQuery("pod", "")
	dsc, _ := strconv.ParseBool(c.DefaultQuery("dsc", "true"))
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	queryExprArray := []string{}
	labels := utils.Labels(start, end)
	for _, label := range labels {
		if c.DefaultQuery(label.(string), "") != "" {
			queryExprArray = append(queryExprArray, utils.GetExpr(label.(string), c.DefaultQuery(label.(string), "")))
		}
	}

	if pod != "" {
		queryExprArray = append(queryExprArray, utils.GetPodExpr(pod))
	}

	if len(queryExprArray) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "缺少查询条件"})
		return
	}

	queryExpr := fmt.Sprintf("{%s}", strings.Join(queryExprArray, ","))
	for _, filter := range filters {
		_, err := regexp.Compile(filter)
		if err != nil {
			utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("regex compile error, %s", err))
			c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
			return
		}
		filter := strings.ReplaceAll(filter, "\\", "\\\\")
		filter = strings.ReplaceAll(filter, "\"", "\\\"")
		queryExpr = fmt.Sprintf("%s |~ \"%s\"", queryExpr, strings.Trim(filter, ""))
	}
	if level != "" {
		levelExpr := utils.GenerateLevelRegex(level)
		if levelExpr != "" {
			queryExpr = fmt.Sprintf("%s %s", queryExpr, levelExpr)
		}
	}

	middleStart := c.DefaultQuery("middleStart", "")
	if middleStart == "" {
		middleStart = start
	}
	middleEnd := c.DefaultQuery("middleEnd", "")
	if middleEnd == "" {
		middleEnd = end
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "500"))

	if limit > 50000 {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "最大支持50000条日志输出"})
		return
	}

	direction := "forward"
	if dsc {
		direction = "backward"
	}

	utils.Log4Zap(zap.InfoLevel).Info(fmt.Sprintf("query expr: %s", queryExpr))
	var queryResults []interface{}
	chartResult := make(map[string]interface{})
	podResults := []interface{}{}
	podSetStr := ""

	queryExpr = url.QueryEscape(queryExpr)
	result := utils.QueryRange(queryExpr, limit, middleStart, middleEnd, direction)

	resultType := result["resultType"]
	if resultType != nil && resultType.(string) == "matrix" {
		// 暂不支持matrix
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "暂不支持matrix类型查询"})
		return
	}

	results := result["result"]
	if results != nil {

		size := 20
		splitDateTimeArray, step := utils.SplitDateTime(start, end, size)
		chartResult["xAxis-data"] = splitDateTimeArray
		chartResult["yAxis-data"] = utils.InitSplitDateTime(size)

		for _, filter := range filters {
			filter = strings.ReplaceAll(filter, "\\\\", "\\")
		}
		for _, result := range results.([]interface{}) {
			resultEle := result.(map[string]interface{})
			stream := resultEle["stream"].(map[string]interface{})

			// pod信息
			podKey := ""
			for key := range stream {
				if strings.Index(key, "pod") > -1 {
					podKey = key
					break
				}
			}
			if podKey != "" && stream[podKey] != nil && strings.Index(podSetStr, stream[podKey].(string)) == -1 {
				podMap := make(map[string]interface{})
				podMap["text"] = stream[podKey]
				podMap["selected"] = false
				podResults = append(podResults, podMap)
				podSetStr += fmt.Sprintf("%s,", stream[podKey].(string))
			}

			values := resultEle["values"].([]interface{})
			for _, value := range values {
				item := make(map[string]interface{})
				item["stream"] = stream
				v := value.([]interface{})
				message := v[1].(string)
				// 保留换行符
				// if len(strings.Trim(message, "\n")) == 0 {
				// 	continue
				// }
				item["info"] = make(map[string]interface{})
				item["info"].(map[string]interface{})["timestamp"] = v[0].(string)
				timestamp, _ := strconv.ParseInt(v[0].(string)[0:13], 10, 64)
				item["info"].(map[string]interface{})["timestampstr"] = time.Unix(0, timestamp*int64(time.Millisecond)).Format("2006-01-02 15:04:05.000")
				item["info"].(map[string]interface{})["message"] = v[1].(string)
				item["info"].(map[string]interface{})["message"] = utils.ShellHighlightShow(item["info"].(map[string]interface{})["message"].(string))
				for _, filter := range filters {
					item["info"].(map[string]interface{})["message"] = utils.RegexHighlightShow(item["info"].(map[string]interface{})["message"].(string), filter)
				}

				// 正则匹配出日志类型
				logLevel := utils.LogLevel(message)
				item["info"].(map[string]interface{})["level"] = logLevel
				item["info"].(map[string]interface{})["animation"] = ""

				// 获取表格数据
				part := utils.TimeInPart(splitDateTimeArray, v[0].(string), step)
				if part >= 0 && part < size {
					chartResult["yAxis-data"].(map[string][]int)[logLevel][part]++
				}

				queryResults = append(queryResults, item)
			}
		}

	} else {
		c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}

	data := make(map[string]interface{})
	data["query"] = queryResults
	data["chart"] = chartResult
	data["pod"] = podResults

	c.JSON(200, data)
}

//
// @Summary Retrieves the list of known values for a given label within a given time span. It accepts the following query parameters in the URL
// @Description limit 2000
// @Param   start path string true "The start time for the query as a nanosecond Unix epoch"
// @Param   end path string true "The end time for the query as a nanosecond Unix epoch"
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/labels/ [get]
func LokiLabels(c *gin.Context) {
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	values := utils.Labels(start, end)
	if values != nil {
		c.JSON(200, values)
	} else {
		c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}

}

//
// @Summary Retrieves the list of known values for a given label within a given time span. It accepts the following query parameters in the URL
// @Description limit 2000
// @Param   start path string true "The start time for the query as a nanosecond Unix epoch"
// @Param   end path string true "The end time for the query as a nanosecond Unix epoch"
// @Accept  json
// @Produce  json
// @Param   label path string true "The label value"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/label/values/ [get]
func LokiLabelValues(c *gin.Context) {
	label := c.DefaultQuery("label", "")
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	values := utils.LabelValues(label, start, end)
	if values != nil {
		c.JSON(200, values)
	} else {
		c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}
}

//
// @Summary Download loki log to log file and accepts the following query parameters in the URL
// @Description file log (max count 50000)
// @Accept  json
// @Produce  json
// @Param   start path string true "The start time for the query as a nanosecond Unix epoch"
// @Param   end path string true "The end time for the query as a nanosecond Unix epoch"
// @Param   filter path string false "The filter condition"
// @Param   level path string false "The log level"
// @Param   dsc path string true "The order to all results"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/export/ [get]
func LokiExport(c *gin.Context) {
	filters := c.QueryArray("filters[]")
	level := c.DefaultQuery("level", "")
	pod := c.DefaultQuery("pod", "")
	dsc, _ := strconv.ParseBool(c.DefaultQuery("dsc", "true"))

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	queryExprArray := []string{}
	labels := utils.Labels(start, end)
	for _, label := range labels {
		if c.DefaultQuery(label.(string), "") != "" {
			queryExprArray = append(queryExprArray, utils.GetExpr(label.(string), c.DefaultQuery(label.(string), "")))
		}
	}

	if pod != "" {
		queryExprArray = append(queryExprArray, utils.GetPodExpr(pod))
	}

	if len(queryExprArray) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "缺少查询条件"})
		return
	}

	queryExpr := fmt.Sprintf("{%s}", strings.Join(queryExprArray, ","))

	for _, filter := range filters {
		queryExpr = fmt.Sprintf("%s |~ \"%s\"", queryExpr, strings.Trim(filter, ""))
	}
	if level != "" {
		levelExpr := utils.GenerateLevelRegex(level)
		if levelExpr != "" {
			queryExpr = fmt.Sprintf("%s %s", queryExpr, levelExpr)
		}
	}

	direction := "forward"
	if dsc {
		direction = "backward"
	}

	limit := 5000
	length := 1

	dir, _ := os.Getwd()
	exportDir := fmt.Sprintf("%s/static/export", dir)
	cmd := fmt.Sprintf("mkdir -p %s", exportDir)
	_, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("mkdir error, %s", err))
		c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "创建文件下载目录失败"})
		return
	}

	filename := fmt.Sprintf("%s.log", time.Now().Format("20060102150405"))
	absolutePath := fmt.Sprintf("%s/static/export/%s", dir, filename)
	file, err := os.Create(absolutePath)
	if err != nil {
		utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("open loki csv file error, %s", err))
	}
	defer file.Close()

	file.WriteString("\xEF\xBB\xBF")

	res := make(map[string]interface{})
	res["exist"] = true

	index := 0
	utils.Log4Zap(zap.InfoLevel).Info(fmt.Sprintf("download expr: %s", queryExpr))
	queryExpr = url.QueryEscape(queryExpr)
	for {
		if index >= 10 {
			break
		}
		index++
		if length == 0 {
			break
		}

		result := utils.QueryRange(queryExpr, limit, start, end, direction)
		resultType := result["resultType"]
		if resultType != nil && resultType.(string) == "matrix" {
			// 暂不支持matrix
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "暂不支持matrix类型查询"})
			return
		}

		results := result["result"]
		if results != nil {
			messages := models.LokiMessages{}
			for _, result := range results.([]interface{}) {
				resultEle := result.(map[string]interface{})
				values := resultEle["values"].([]interface{})
				for _, value := range values {
					v := value.([]interface{})
					// 保留换行符
					// if len(strings.Trim(v[1].(string), "\n")) == 0 {
					// 	continue
					// }
					messages = append(messages, models.LokiMessage{Timestamp: v[0].(string), Message: v[1].(string)})
				}
			}
			length = len(messages)

			if length > 0 {
				if dsc {
					sort.Sort(messages)
					end = messages[len(messages)-1].Timestamp
				} else {
					sort.Sort(sort.Reverse(messages))
					start = messages[len(messages)-1].Timestamp
				}
				for _, message := range messages {
					file.WriteString(message.Message)
				}
			}
		} else {
			break
		}
	}

	res["download"] = filename
	c.JSON(200, res)
}

//
// @Summary Get loki log context from grafana loki and accepts the following query parameters in the URL
// @Description limit 2000
// @Accept  json
// @Produce  json
// @Param   start path string true "The start time for the query as a nanosecond Unix epoch"
// @Param   end path string true "The end time for the query as a nanosecond Unix epoch"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/context/ [get]
func LokiContext(c *gin.Context) {
	queryExprArray := []string{}

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	labels := utils.Labels(start, end)
	for _, label := range labels {
		if c.DefaultQuery(label.(string), "") != "" {
			queryExprArray = append(queryExprArray, utils.GetExpr(label.(string), c.DefaultQuery(label.(string), "")))
		}
	}

	if len(queryExprArray) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "缺少查询条件"})
		return
	}

	queryExpr := fmt.Sprintf("{%s}", strings.Join(queryExprArray, ","))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	direction := c.DefaultQuery("direction", "")
	if direction == "next" {
		direction = "forward"
	} else {
		direction = "backward"
	}

	utils.Log4Zap(zap.InfoLevel).Info(fmt.Sprintf("context expr: %s", queryExpr))
	queryExpr = url.QueryEscape(queryExpr)
	result := utils.QueryRange(queryExpr, limit, start, end, direction)
	queryResults := []interface{}{}

	resultType := result["resultType"]
	if resultType != nil && resultType.(string) == "matrix" {
		// 暂不支持matrix
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "暂不支持matrix类型查询"})
		return
	}

	results := result["result"]
	if results != nil {
		for _, result := range results.([]interface{}) {
			resultEle := result.(map[string]interface{})
			values := resultEle["values"].([]interface{})
			for _, value := range values {
				item := make(map[string]interface{})
				v := value.([]interface{})
				// 保留换行符
				// if len(strings.Trim(v[1].(string), "\n")) == 0 {
				// 	continue
				// }
				// 正则匹配出日志类型
				logLevel := utils.LogLevel(v[1].(string))
				item["timestamp"] = v[0]
				item["level"] = logLevel
				item["message"] = utils.ShellHighlightShow(v[1].(string))
				queryResults = append(queryResults, item)
			}
		}
	} else {
		c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}

	c.JSON(200, queryResults)
}

//
// @Summary WebSocket endpoint that will stream log messages based on a query. It accepts the following query parameters in the URL
// @Description limit 2000
// @Accept  json
// @Produce  json
// @Param   start path string true "The end time for the query as a nanosecond Unix epoch"
// @Param   pod path string false "The pod filter condition to perform"
// @Param   filter path string false "The filter condition"
// @Param   level path string false "The log level"
// @Param   limit path string false "The max number of entries to return"
// @Success 200 {string} string	"[]"
// @Router /ws/tail/ [get]
func LokiTail(c *gin.Context) {
	level := c.DefaultQuery("level", "")
	pod := c.DefaultQuery("pod", "")

	filtersStr := c.DefaultQuery("filters", "")
	filters := strings.Split(filtersStr, ",")
	start := c.DefaultQuery("start", "")

	queryExprArray := []string{}

	t := time.Now().Unix()
	end := fmt.Sprintf("%d000000000", t)
	labels := utils.Labels(start, end)
	for _, label := range labels {
		if c.DefaultQuery(label.(string), "") != "" {
			queryExprArray = append(queryExprArray, utils.GetExpr(label.(string), c.DefaultQuery(label.(string), "")))
		}
	}

	if pod != "" {
		queryExprArray = append(queryExprArray, utils.GetPodExpr(pod))
	}

	if len(queryExprArray) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "缺少查询条件"})
		return
	}

	queryExpr := fmt.Sprintf("{%s}", strings.Join(queryExprArray, ","))
	for _, filter := range filters {
		_, err := regexp.Compile(filter)
		if err != nil {
			utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("regex compile error, %s", err))
			c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
			return
		}
		filter := strings.ReplaceAll(filter, "\\", "\\\\")
		filter = strings.ReplaceAll(filter, "\"", "\\\"")
		queryExpr = fmt.Sprintf("%s |~ \"%s\"", queryExpr, strings.Trim(filter, ""))
	}
	if level != "" {
		levelExpr := utils.GenerateLevelRegex(level)
		if levelExpr != "" {
			queryExpr = fmt.Sprintf("%s %s", queryExpr, levelExpr)
		}
	}

	queryExpr = url.QueryEscape(queryExpr)

	params := make(map[string]string)
	params["query"] = queryExpr
	params["start"] = start
	params["limit"] = "500"
	params["delay_for"] = "0"

	clientConnect := utils.LokiWebsocketClient(params)
	if clientConnect == nil {
		return
	}
	defer clientConnect.Close()

	serverConnect := utils.LokiWebsocketServer(c.Writer, c.Request)
	if serverConnect == nil {
		return
	}
	defer serverConnect.Close()

	chanSendMessage := make(chan utils.WsMessage)
	chanReceiveMessage := make(chan utils.WsMessage)
	chanSignal := make(chan int)

	go utils.WebSocketClientHandler(clientConnect, chanSendMessage, chanSignal)
	go utils.WebSocketServerHandler(serverConnect, chanReceiveMessage, chanSignal)

	for {
		select {
		case wsClientMessage, ok := <-chanSendMessage:
			if !ok {
				return
			}
			data := utils.LokiWebsocketMessageConstruct(wsClientMessage.Data, filters)
			err := serverConnect.WriteMessage(wsClientMessage.MessageType, data)
			if err != nil {
				utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("send message to viewer error, %s", err))
				return
			}
		case wsServerMessage, ok := <-chanReceiveMessage:
			if !ok {
				return
			}
			data := string(wsServerMessage.Data)
			if data == "close" {
				return
			}
		}
	}
}
