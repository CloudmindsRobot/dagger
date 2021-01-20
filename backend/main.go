package main

import (
	"dagger/backend/databases"
	router "dagger/backend/routers"
	"dagger/backend/runtime"
	"dagger/backend/utils"
	"fmt"
	"os"

	"go.uber.org/zap"
)

// @title dagger backend api
// @version 2.0.0
// @description this is dagger backend api server
// @BasePath /
func main() {
	if runtime.LokiServer == "" {
		runtime.LokiServer = os.Getenv("LOKI_SERVER")
		if runtime.LokiServer == "" {
			runtime.LokiServer, _ = runtime.Cfg.GetValue("loki", "address")
			if runtime.LokiServer == "" {
				utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("start server error: missing loki-server param or LOKI_SERVER env"))
				return
			}
		}
	}

	err := utils.CacheRule()
	if err != nil {
		utils.Log4Zap(zap.ErrorLevel).Error("init cache failed!")
		return
	}

	db, _ := databases.DB.DB()
	defer db.Close()

	router.InitRouter()
}
