package api

import (
	"net/http"

	"github.com/ericbg27/RegistryAPI/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config    util.Config
	dbManager *gorm.DB
	router    *gin.Engine
}

func NewServer(dbManager *gorm.DB, config util.Config) (server *Server, err error) {
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
