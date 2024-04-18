package main

import (
	"log/slog"
	"net/http"
	"postgrespro-executor-service/internal/http-server/router"
	"postgrespro-executor-service/internal/repositories/postgres"
	"postgrespro-executor-service/internal/repositories/redis"
	"postgrespro-executor-service/internal/utils/error_check"
)

func main() {
	slog.Info("starting app..")

	db, err := postgres.New()
	error_check.CheckError(err)

	redisClient, err := redis.New()
	error_check.CheckError(err)

	err = http.ListenAndServe(":8080", router.Routes(redisClient, db))

	slog.Info("shutdown server")
}
