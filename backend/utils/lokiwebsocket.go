package utils

import (
	"dagger/backend/runtime"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WsMessage struct {
	MessageType int
	Data        []byte
}

func LokiWebsocketClient(params map[string]string) *websocket.Conn {

	paramArray := []string{}
	for k, v := range params {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", k, v))
	}

	scheme := "ws"
	if strings.Index(runtime.LokiServer, "https") > -1 {
		scheme = "wss"
	}

	u := url.URL{Scheme: scheme, Host: runtime.LokiServer[strings.Index(runtime.LokiServer, "/")+2:], Path: "/loki/api/v1/tail", RawQuery: strings.Join(paramArray, "&")}

	clientConnect, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("connect loki websocket %s error: %s", u.String(), err))
		return nil
	}
	return clientConnect
}

func WebSocketClientHandler(clientConnect *websocket.Conn, ch chan WsMessage, signal chan int) {
	timer := time.NewTimer(time.Millisecond * 500)
	for {
		select {
		case sg := <-signal:
			if sg == 0 {
				if timer != nil {
					timer.Stop()
				}
				close(ch)
				return
			}
		case <-timer.C:
			messageType, data, err := clientConnect.ReadMessage()
			if err != nil {
				close(ch)
				Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("read message from loki error, %s", err))
				timer.Stop()
				break
			}
			if len(data) > 0 {
				ch <- WsMessage{MessageType: messageType, Data: data}
			}
			timer.Reset(time.Microsecond * 500)
		}
	}
}

func LokiWebsocketServer(write http.ResponseWriter, request *http.Request) *websocket.Conn {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("unknown error, %s", reason))
		},
	}
	serverConnect, err := upgrader.Upgrade(write, request, nil)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("connect viewer websocket error: %s", err))
		return nil
	}
	return serverConnect
}

func WebSocketServerHandler(serverConnect *websocket.Conn, ch chan WsMessage, signal chan int) {
	for {
		messageType, data, err := serverConnect.ReadMessage()
		if err != nil {
			close(ch)
			Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("read message from viewer error, %s", err))
			return
		}
		if len(data) > 0 {
			if string(data) == "close" {
				signal <- 0
				ch <- WsMessage{MessageType: messageType, Data: data}
				close(ch)
				return
			}
			ch <- WsMessage{MessageType: messageType, Data: data}
		}
	}
}

func LokiWebsocketMessageConstruct(data []byte, filters []string) []byte {
	var queryResults []interface{}
	var message map[string]interface{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("Unmarshal message error, %s", err))
		return []byte{}
	}

	results := message["streams"]
	if results != nil {
		for _, result := range results.([]interface{}) {
			resultEle := result.(map[string]interface{})
			stream := resultEle["stream"].(map[string]interface{})
			values := resultEle["values"].([]interface{})
			for _, value := range values {
				item := make(map[string]interface{})
				item["stream"] = stream
				v := value.([]interface{})
				message := v[1].(string)
				if len(strings.Trim(message, "\n")) == 0 {
					continue
				}
				item["info"] = make(map[string]interface{})
				item["info"].(map[string]interface{})["timestamp"] = v[0].(string)
				timestamp, _ := strconv.ParseInt(v[0].(string)[0:10], 10, 64)
				item["info"].(map[string]interface{})["timestampstr"] = time.Unix(0, timestamp*int64(time.Millisecond)).Format("2006-01-02 15:04:05.000")
				item["info"].(map[string]interface{})["message"] = v[1].(string)
				item["info"].(map[string]interface{})["message"] = ShellHighlightShow(item["info"].(map[string]interface{})["message"].(string))
				for _, filter := range filters {
					item["info"].(map[string]interface{})["message"] = RegexHighlightShow(item["info"].(map[string]interface{})["message"].(string), filter)
				}

				// 正则匹配出日志类型
				logLevel := LogLevel(v[1].(string))
				item["info"].(map[string]interface{})["level"] = logLevel
				item["info"].(map[string]interface{})["animation"] = "background-color: yellow;transition: background-color 2s;"

				queryResults = append(queryResults, item)
			}
		}
	}
	data, _ = json.Marshal(queryResults)
	return data
}
