package app

import (
	"fmt"
	"tjus-be/internal/config"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Engine *gin.Engine
	Cfg    *config.Config
	DB     *pgxpool.Pool
	Env    *config.Environment
	Redis  *redis.Client
}

func (s *Server) SetupConfiguration() {
	cfg := config.InitConfig()
	s.Cfg = cfg
	s.DB = cfg.DB
	s.Env = cfg.Env
	s.Redis = cfg.Redis
}

func (s *Server) SetupMiddlewares() {
	s.Engine.Use(gin.Recovery())
	s.Engine.Use(middleware.LoggerMiddleware(s.Env.IsProduction()))
}

func (s *Server) SetupRouters() {
	SetRouter(s)
}

func (s *Server) Close() {
	if s.Cfg != nil {
		s.Cfg.Close()
	}
}

func NewServer() *Server {
	server := &Server{}
	server.SetupConfiguration()

	if server.Env.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	server.Engine = gin.New()
	server.SetupMiddlewares()
	server.SetupRouters()

	return server
}

func (s *Server) Start() {
	fmt.Println("Starting the application...")
	s.Engine.Run(":" + s.Env.Port)
}
