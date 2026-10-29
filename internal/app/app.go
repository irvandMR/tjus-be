package app

import (
	"tjus-be/internal/config"
)

func Run() {
	env := config.LoadEnvironment()
	db := config.InitDatabaseConfig(env)
	defer db.Close()

	server := NewServer(db, env)
	server.Start()
}
