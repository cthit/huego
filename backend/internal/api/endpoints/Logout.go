package endpoints

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:   "/api",
		MaxAge: -1,
	})
	session.Clear()
	err := session.Save()

	if err != nil {
		log.Printf("Failed to save session: %v", err)
		c.String(http.StatusInternalServerError, "Failed to logout")
		return
	}

	c.Header("location", "/")
	c.String(http.StatusOK, "Logged out successfully")
}
