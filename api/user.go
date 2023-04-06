package api

import (
	"net/http"

	"github.com/ericbg27/RegistryAPI/db"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required"` // TODO: Validação de telefone
	UserName string `json:"user_name" binding:"required,alphanum,min=6"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) createUser(c *gin.Context) {
	var userReq createUserRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{
			"name":    "BadRequest",
			"message": "Incorrect parameters sent in request",
		})
		return
	}

	var userParams = db.CreateUserParams{
		FullName: userReq.FullName,
		Phone:    userReq.Phone,
		UserName: userReq.UserName,
		Password: userReq.Password,
	}

	_, err := s.DbConnector.CreateUser(userParams)
	if err != nil {
		dbErr, ok := err.(*db.BadInputError)
		if ok {
			c.JSON(http.StatusBadRequest, &gin.H{
				"name":    "AlreadyExists",
				"message": dbErr.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
		return
	}

	c.JSON(http.StatusCreated, &gin.H{
		"message": "User created successfully",
	})
	return
}
