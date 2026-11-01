package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
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
	srv := &http.Server{
		Addr:    ":" + s.Env.Port,
		Handler: s.Engine,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Starting the application on port %s...", s.Env.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
