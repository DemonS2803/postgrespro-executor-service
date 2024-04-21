package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"postgrespro-executor-service/internal/http-server/handlers/auth"
	"postgrespro-executor-service/internal/http-server/handlers/command"
	"postgrespro-executor-service/internal/repositories/postgres"
	"postgrespro-executor-service/internal/repositories/redis"
)

func Routes(redisClient *redis.Redis, db *postgres.Storage) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(auth.NewTokenHandler())

	router.Get("/ping", command.Ping())
	router.Post("/command", command.CreateCommand(db, redisClient))
	router.Get("/command/exec/{id}", command.ExecuteCommandController(db, redisClient))
	router.Get("/command/result/{id}", command.GetCompletedCommandByIdController(db, redisClient))
	router.Get("/command/stop/{id}", command.StopCommandController(db, redisClient))
	router.Get("/command/{id}", command.GetCommandByIdController(db, redisClient))
	router.Get("/commands", command.GetCommandList(db, redisClient))

	return router
}
