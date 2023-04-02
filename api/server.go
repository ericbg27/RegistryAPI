package api

import (
	"net/http"

	"github.com/ericbg27/RegistryAPI/db"
	"github.com/ericbg27/RegistryAPI/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config    util.Config
	dbManager *db.DBManager
	router    *gin.Engine
}

func NewServer(dbManager *db.DBManager, config util.Config) (server *Server, err error) {
	server = &Server{
		dbManager: dbManager,
		config:    config,
	}

	server.setupRouter()

	return
}

func (s *Server) setupRouter() {
	s.router = gin.Default()

	s.router.GET("/", s.healthCheck)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{})
}

// Start runs the server listening on the specified port
func (s *Server) Start() {
	s.router.Run(s.config.ServerAddress)
}
