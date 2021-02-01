package router

import (
	"dagger/backend/controllers"
	"dagger/backend/middlewares"
	"dagger/backend/runtime"
	session "dagger/backend/sessions"
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

	debug, _ := runtime.Cfg.Bool("global", "debug")
	if !debug {
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

	auth := router.Group("/api/v1/loki/auth")
	auth.POST("/login", controllers.Login)
	auth.POST("/register", controllers.Register)

	api := router.Group("/api/v1/loki", middlewares.JWTCheck())
	api.GET("/auth/userinfo", controllers.GetUserInfo)

	api.GET("/label/values", controllers.LokiLabelValues)
	api.GET("/labels", controllers.LokiLabels)
	api.GET("/query_range", controllers.LokiList)
	api.GET("/context", controllers.LokiContext)
	api.GET("/export", controllers.LokiExport)

	api.GET("/history", controllers.LokiHistoryList)
	api.POST("/history/create", controllers.LokiHistoryCreate)
	api.DELETE("/history/delete/:id", controllers.LokiHistoryDelete)

	api.GET("/snapshot", controllers.LokiSnapshotList)
	api.POST("/snapshot/create", controllers.LokiSnapshotCreate)
	api.DELETE("/snapshot/delete/:id", controllers.LokiSnapshotDelete)
	api.GET("/snapshot/detail/:id", controllers.LokiSnapshotDetail)

	alertEnabled, _ := runtime.Cfg.Bool("global", "alert_enabled")
	if alertEnabled {
		api.GET("/rule", controllers.LokiRuleList)
		api.POST("/rule/create", controllers.LokiRuleCreate)
		api.POST("/rule/update/:id", controllers.LokiRuleUpdate)
		api.DELETE("/rule/delete/:id", controllers.LokiRuleDelete)
		api.GET("/rule/download", controllers.LokiRuleDownload)

		api.GET("/group", controllers.LokiUserGroupList)
		api.POST("/group/create", controllers.LokiUserGroupCreate)
		api.DELETE("/group/delete/:id", controllers.LokiUserGroupDelete)
		api.POST("/group/update/:id", controllers.LokiUserGroupUpdate)
		api.POST("/group/join", controllers.LokiUserGroupJoin)
		api.POST("/group/leave", controllers.LokiUserGroupLeave)

		api.POST("/event/archive", controllers.LokiEventArchive)
		api.GET("/event", controllers.LokiEventList)
		api.GET("/event/details/:id", controllers.LokiEventDetailList)
	}

	api.GET("/user", controllers.LokiUserList)

	api.GET("/logql", controllers.TransformLogQL)

	setting := router.Group("/api/v1/loki/settings")
	setting.GET("/load", controllers.LoadSettings)

	dir, _ := os.Getwd()
	api.StaticFS("/static", http.Dir(dir+"/static"))

	ws := router.Group("/ws")
	ws.GET("/tail", controllers.LokiTail)

	apiv1 := router.Group("/api/v1")
	apiv1.POST("/alerts", controllers.LokiEventCollect)
	apiv2 := router.Group("/api/v2")
	apiv2.POST("/alerts", controllers.LokiEventCollect)

	swagger := router.Group("/swagger")
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(fmt.Sprintf(":%d", runtime.Port))
}
