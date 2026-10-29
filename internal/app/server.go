package app

import (
	"fmt"
	"tjus-be/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Engine *gin.Engine
	DB     *pgxpool.Pool
	Env    *config.Environment
}

func (s *Server) SetupRouters() {
	SetRouter(s)
}

func NewServer(db *pgxpool.Pool, env *config.Environment) *Server {
	e := gin.Default()

	server := &Server{
		Engine: e,
		DB:     db,
		Env:    env,
	}

	server.SetupRouters()

	return server
}

func (s *Server) Start() {
	fmt.Println("Starting the application...")
	s.Engine.Run(":" + s.Env.Port)
}
