package api

import (
	"net/http"
	"time"

	"github.com/ericbg27/RegistryAPI/db"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
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

type getUserRequest struct {
	UserName string `form:"user_name" binding:"required"`
}

type getUserResponse struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	UserName string `json:"user_name"`
}

func (s *Server) getUser(c *gin.Context) {
	var userReq getUserRequest

	if err := c.ShouldBindQuery(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"name":    "BadRequest",
			"message": "Incorrect parameters sent in request",
		})
		return
	}

	user, err := s.DbConnector.GetUser(userReq.UserName)
	if err != nil {
		notFoundErr, ok := err.(*db.NotFoundError)
		if ok {
			c.JSON(http.StatusNotFound, gin.H{
				"name":    "NotFound",
				"message": notFoundErr.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
	}

	userRes := &getUserResponse{
		FullName: user.FullName,
		Phone:    user.Phone,
		UserName: user.UserName,
	}

	c.JSON(http.StatusOK, userRes)
}

type getUsersRequest struct {
	PageIndex int `form:"page"`
	Offset    int `form:"offset"`
}

type getUsersUserResponse struct {
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type getUsersResponse struct {
	Users []*getUsersUserResponse `json:"users"`
}

const minOffset = 5

func (s *Server) getUsers(c *gin.Context) {
	var usersReq getUsersRequest

	if err := c.BindQuery(&usersReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"name":    "BadRequest",
			"message": "Incorrect parameters sent in request",
		})
		return
	}

	if usersReq.Offset == 0 {
		usersReq.Offset = minOffset
	}

	getUsersParams := db.GetUsersParams{
		PageIndex: usersReq.PageIndex,
		Offset:    usersReq.Offset,
	}

	users, err := s.DbConnector.GetUsers(getUsersParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
	}

	usersRes := &getUsersResponse{
		Users: []*getUsersUserResponse{},
	}
	for _, user := range users {
		userRes := &getUsersUserResponse{
			FullName:  user.FullName,
			Phone:     user.Phone,
			UserName:  user.UserName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		usersRes.Users = append(usersRes.Users, userRes)
	}

	c.JSON(http.StatusOK, usersRes)
}
