package api

import (
	"net/http"

	"github.com/ericbg27/RegistryAPI/db"
	"github.com/ericbg27/RegistryAPI/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Config      util.Config
	DbConnector db.DBConnector
	Router      *gin.Engine
}

func NewServer(dbConnector db.DBConnector, config util.Config) (server *Server, err error) {
	server = &Server{
		DbConnector: dbConnector,
		Config:      config,
	}

	server.setupRouter()

	return
}

func (s *Server) setupRouter() {
	s.Router = gin.Default()

	s.Router.GET("/", s.healthCheck)
	s.Router.POST("/users", s.createUser)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{})
}

// Start runs the server listening on the specified port
func (s *Server) Start() {
	s.Router.Run(s.Config.ServerAddress)
}