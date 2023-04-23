package api

import (
	"net/http"

	"github.com/ericbg27/RegistryAPI/db"
	"github.com/gin-gonic/gin"
)

func (s *Server) isAdmin(c *gin.Context) {
	userReq, ok := c.Keys["currentUser"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"name":    "Unauthorized",
			"message": "User is not authorized to access this resource",
		})
		c.Abort()
		return
	}

	user, ok := userReq.(*db.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"name":    "Unauthorized",
			"message": "User is not authorized to access this resource",
		})
		c.Abort()
		return
	}

	if !user.Admin {
		c.JSON(http.StatusForbidden, gin.H{
			"name":    "Forbidden",
			"message": "User is not allowed to access this resource",
		})
		c.Abort()
		return
	}

	c.Next()
}
