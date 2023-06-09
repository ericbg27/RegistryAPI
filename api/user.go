package api

import (
	"net/http"
	"time"

	"github.com/ericbg27/RegistryAPI/db"
	"github.com/ericbg27/RegistryAPI/util"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required,isPhone"`
	UserName string `json:"user_name" binding:"required,alphanum,min=6"`
	Password string `json:"password" binding:"required,min=6,validPassword"`
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
		return
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
		return
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

type loginUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	Token string `json:"token"`
}

func (s *Server) loginUser(c *gin.Context) {
	var loginReq loginUserRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"name":    "BadRequest",
			"message": "Incorrect parameters sent in request",
		})
		return
	}

	user, err := s.DbConnector.GetUser(loginReq.UserName)
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
		return
	}

	if !util.ComparePassword(user.Password, loginReq.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"name":    "Unauthorized",
			"message": "Wrong password sent in request",
		})
		return
	}

	token, err := s.Maker.CreateToken(user.UserName, s.Config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
		return
	}

	updateParams := db.UpdateUserParams{
		ID:         user.ID,
		FullName:   user.FullName,
		Phone:      user.Phone,
		Password:   user.Password,
		LoginToken: token,
	}

	if err = s.DbConnector.UpdateUser(updateParams); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
		return
	}

	loginRes := loginUserResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, loginRes)
}

type updateUserRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone" binding:"isPhone"`
}

func (s *Server) updateUser(c *gin.Context) {
	var updateUserParams updateUserRequest

	if err := c.BindJSON(&updateUserParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"name":    "BadRequest",
			"message": "Incorrect parameters sent in request",
		})
		return
	}

	user, ok := c.Keys["currentUser"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
		return
	}

	currentUser, ok := user.(*db.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
		return
	}

	updateParams := db.UpdateUserParams{
		ID:         currentUser.ID,
		FullName:   updateUserParams.FullName,
		Phone:      updateUserParams.Phone,
		Password:   currentUser.Password,
		LoginToken: currentUser.LoginToken,
	}

	if err := s.DbConnector.UpdateUser(updateParams); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"name":    "InternalServerError",
			"message": "Unexpected server error. Try again later",
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

type deleteUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
}

func (s *Server) deleteUser(c *gin.Context) {
	var deleteUserReq deleteUserRequest

	if err := c.ShouldBindJSON(&deleteUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"name":    "BadRequest",
			"message": "Incorrect parameters sent in request",
		})
		return
	}

	userReq, _ := c.Keys["currentUser"]
	loggedUser, _ := userReq.(*db.User)

	if !loggedUser.Admin && loggedUser.UserName != deleteUserReq.UserName {
		c.JSON(http.StatusForbidden, gin.H{
			"name":    "Forbidden",
			"message": "User is not allowed to access this resource",
		})
		return
	}

	if err := s.DbConnector.DeleteUser(deleteUserReq.UserName); err != nil {
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
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
