package router

import (
	"dagger/backend/gin/controllers"
	"dagger/backend/gin/middlewares"
	"dagger/backend/gin/runtime"
	session "dagger/backend/gin/sessions"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "dagger/backend/docs"
)

func InitRouter() {
	if !runtime.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(sessions.Sessions("loki-backend-go-session", session.Store))
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s [INFO] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(gin.Recovery())

	api := router.Group("/api/v1/loki")

	api.POST("/auth/login", controllers.Login)
	api.POST("/auth/register", controllers.Register)
	api.GET("/auth/userinfo", middlewares.JWTCheck(), controllers.GetUserInfo)

	api.GET("/label/values", middlewares.JWTCheck(), controllers.LokiLabelValues)
	api.GET("/labels", middlewares.JWTCheck(), controllers.LokiLabels)
	api.GET("/query_range", middlewares.JWTCheck(), controllers.LokiList)
	api.GET("/context", middlewares.JWTCheck(), controllers.LokiContext)
	api.GET("/export", middlewares.JWTCheck(), controllers.LokiExport)

	api.GET("/history", middlewares.JWTCheck(), controllers.LokiHistoryList)
	api.POST("/history/create", middlewares.JWTCheck(), controllers.LokiHistoryCreate)
	api.DELETE("/history/delete/:id", middlewares.JWTCheck(), controllers.LokiHistoryDelete)

	api.GET("/snapshot", middlewares.JWTCheck(), controllers.LokiSnapshotList)
	api.POST("/snapshot/create", middlewares.JWTCheck(), controllers.LokiSnapshotCreate)
	api.DELETE("/snapshot/delete/:id", middlewares.JWTCheck(), controllers.LokiSnapshotDelete)
	api.GET("/snapshot/detail/:id", middlewares.JWTCheck(), controllers.LokiSnapshotDetail)

	api.GET("/settings/load", controllers.LoadSettings)

	dir, _ := os.Getwd()
	api.StaticFS("/static", http.Dir(dir+"/static"))

	ws := router.Group("/ws")
	ws.GET("/tail", controllers.LokiTail)

	swagger := router.Group("/swagger")
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(fmt.Sprintf(":%d", runtime.Port))
}
