package api

import (
	"net/http"

	"github.com/ericbg27/RegistryAPI/db"
	"github.com/ericbg27/RegistryAPI/token"
	"github.com/ericbg27/RegistryAPI/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Config      util.Config
	DbConnector db.DBConnector
	Router      *gin.Engine
	Maker       token.Maker
}

func NewServer(dbConnector db.DBConnector, config util.Config, maker token.Maker) (server *Server, err error) {
	server = &Server{
		DbConnector: dbConnector,
		Config:      config,
		Maker:       maker,
	}

	server.setupRouter()

	return
}

func (s *Server) setupRouter() {
	s.Router = gin.Default()

	v1 := s.Router.Group("/v1")
	{
		v1.GET("/", s.healthCheck)

		v1User := v1.Group("/user")
		{
			v1User.GET("/", s.checkAuth, s.getUser)
			v1User.PUT("/", s.checkAuth, s.updateUser)
			v1User.POST("/", s.createUser)
			v1User.POST("/login", s.loginUser)
		}

		v1Users := v1.Group("/users")
		{
			v1Users.GET("/", s.getUsers)
		}
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{})
}

// Start runs the server listening on the specified port
func (s *Server) Start() {
	s.Router.Run(s.Config.ServerAddress)
}
