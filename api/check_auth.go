package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) checkAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"name":    "Unauthorized",
			"message": "User is not authorized to access this resource",
		})
		c.Abort()
		return
	}

	token := tokenParts[1]
	payload, err := s.Maker.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"name":    "Unauthorized",
			"message": "User is not authorized to access this resource",
		})
		c.Abort()
		return
	}

	user, err := s.DbConnector.GetUser(payload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
		c.Abort()
		return
	}

	if user.LoginToken != token {
		c.JSON(http.StatusUnauthorized, gin.H{
			"name":    "Unauthorized",
			"message": "User is not authorized to access this resource",
		})
		c.Abort()
		return
	}

	c.Set("tokenPayload", token)
	c.Set("currentUser", user)
	c.Next()
}
