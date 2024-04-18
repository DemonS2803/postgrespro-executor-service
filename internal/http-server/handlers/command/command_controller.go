package command

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"postgrespro-executor-service/internal/models"
	"postgrespro-executor-service/internal/repositories/postgres"
	"postgrespro-executor-service/internal/repositories/redis"
	"postgrespro-executor-service/internal/services/command_service"
	resp "postgrespro-executor-service/internal/utils/response"
	"strconv"
	"time"
)

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("success ping")
		resp.Send200Success(w, r)
	}
}

func ExecuteCommandController(storage *postgres.Storage, redisClient *redis.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commandId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			slog.Error("error while parsing request id", err)
			resp.Send400Error(w, r)
			return
		}

		requestId := middleware.GetReqID(r.Context())
		go command_service.ExecuteCommand(storage, redisClient, commandId, requestId)
		for !redis.HasCompletedIdForRequest(redisClient, requestId) {
			time.Sleep(time.Millisecond)
		}

		completedId := redis.GetCompletedIdByRequestId(redisClient, requestId)
		if completedId < 1 {
			resp.Send404Error(w, r)
			return
		}

		render.JSON(w, r, models.CompletedCommandRequest{
			Id: completedId,
		})
	}
}

func GetCommandByIdController(storage *postgres.Storage, redisClient *redis.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			slog.Error("error while parsing request id", err)
			resp.Send400Error(w, r)
			return
		}
		command, err := command_service.GetCommandById(storage, redisClient, id)
		if err != nil {
			resp.Send404Error(w, r)
			return
		}

		render.JSON(w, r, command)
	}
}

func GetCompletedCommandByIdController(storage *postgres.Storage, redisClient *redis.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			slog.Error("error while parsing request id", err)
			resp.Send400Error(w, r)
			return
		}
		command, err := command_service.GetCompletedCommandById(storage, redisClient, id)
		if err != nil {
			resp.Send404Error(w, r)
			return
		}

		render.JSON(w, r, command)
	}
}

func GetCommandList(storage *postgres.Storage, redisClient *redis.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 10
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			offset = 0
		}

		list, err := command_service.GetCommandsListByLimitAndOffset(storage, redisClient, limit, offset)
		if err != nil {
			resp.Send404Error(w, r)
			return
		}

		render.JSON(w, r, list)
	}
}

func CreateCommand(storage *postgres.Storage, redisClient *redis.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createCommand models.CreateCommandRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&createCommand)
		validate := validator.New()
		err2 := validate.Struct(createCommand)
		if err != nil || err2 != nil {
			slog.Error("error while parsing request body", err)
			resp.Send400Error(w, r)
			return
		}

		command, err := command_service.CreateCommand(storage, redisClient, createCommand)
		if err != nil {
			resp.Send500Error(w, r)
			return
		}
		render.JSON(w, r, command)
	}
}
