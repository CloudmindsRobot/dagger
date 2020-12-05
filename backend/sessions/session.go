package sessions

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var Store = cookie.NewStore([]byte("dagger/backend-secret"))

func SetSession(c *gin.Context, username string) string {
	session := sessions.Default(c)
	sessionId := "dagger/backend-secret"
	session.Set(sessionId, username)
	session.Save()
	return sessionId
}

func GetSession(c *gin.Context, sessionId string) (username string) {
	session := sessions.Default(c)
	result := session.Get(sessionId)
	username = result.(string)
	return
}
