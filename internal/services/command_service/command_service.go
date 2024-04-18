package command_service

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"postgrespro-executor-service/internal/models"
	"postgrespro-executor-service/internal/repositories/postgres"
	"postgrespro-executor-service/internal/repositories/redis"
	"strconv"
)

func GetCommandById(storage *postgres.Storage, redisClient *redis.Redis, commandId int) (models.Command, error) {
	command, err := redis.GetCommand(redisClient, commandId)
	if err == nil {
		return command, nil
	}
	dbCommand, err := postgres.FindCommandById(storage, commandId)
	if err != nil {
		slog.Error("error when get command:", err)
		return models.Command{}, err
	}
	redis.PutCommand(redisClient, dbCommand)

	return dbCommand, nil
}

func CreateCommand(storage *postgres.Storage, redisClient *redis.Redis, createCommand models.CreateCommandRequest) (models.Command, error) {
	command, err := postgres.CreateCommandAndReturn(storage, createCommand)
	if err != nil {
		slog.Error("error when insert new command", err)
		return models.Command{}, err
	}
	redis.PutCommand(redisClient, command)

	return command, nil
}

func GetCommandsListByLimitAndOffset(storage *postgres.Storage, redisClient *redis.Redis, limit int, offset int) ([]models.Command, error) {
	list, err := redis.GetCommandsList(redisClient, limit, offset)
	if err == nil {
		return list, nil
	}

	list, err = postgres.FindCommandListByLimitAndOffset(storage, limit, offset)
	if err != nil {
		slog.Error("error while get command list", err)
		return nil, err
	}
	redis.PutCommandsList(redisClient, limit, offset, list)
	return list, nil
}

func GetCompletedCommandById(storage *postgres.Storage, redisClient *redis.Redis, completedCommandId int) (models.CompletedCommandRequest, error) {
	command, err := redis.GetCompletedCommand(redisClient, completedCommandId)
	if err == nil {
		return command, nil
	}
	dbCommand, err := postgres.FindCompletedCommandById(storage, completedCommandId)
	if err != nil {
		slog.Error("error when get command result:", err)
		return models.CompletedCommandRequest{}, err
	}
	redis.PutCompletedCommand(redisClient, dbCommand)

	return dbCommand, nil
}

func ExecuteCommand(storage *postgres.Storage, redisClient *redis.Redis, commandId int, requestId string) {
	//slog.Info("starting command", commandId)

	dbCommand, err := postgres.FindCommandById(storage, commandId)
	if err != nil {
		slog.Error("error when get command:", err)
		redis.PutRequest(redisClient, requestId, "-1")
		return
	}

	completedCommandId, err := postgres.CreateCompletedCommand(storage, dbCommand)
	redis.PutRequest(redisClient, requestId, strconv.Itoa(completedCommandId))
	filename, err := CreateTempFile(dbCommand.Id, dbCommand.Code)
	if err != nil {
		slog.Error("error when create tmp file:", err)
		postgres.UpdateCompletedCommandById(storage, completedCommandId, "server error")
		return
	}

	com := exec.Command("/bin/bash", fmt.Sprintf("./resources/scripts/tmp/%s", filename))
	out, err := com.Output()
	if err != nil {
		slog.Error("error during running command", err)
		postgres.UpdateCompletedCommandById(storage, completedCommandId, err.Error())
	} else {
		postgres.UpdateCompletedCommandById(storage, completedCommandId, string(out))
	}

	//slog.Info(string(out))
}

func CreateTempFile(commandId int, code string) (string, error) {
	filename := fmt.Sprintf("tmp%d.sh", commandId)

	if IsTempFileExists(filename) {
		slog.Info("file ", filename, "already exists")
		return filename, nil
	}
	file, err := os.Create(fmt.Sprintf("./resources/scripts/tmp/%s", filename))
	os.Chmod(fmt.Sprintf("./resources/scripts/tmp/%s", filename), 777)
	slog.Info("created file", filename)
	defer file.Close()
	n, err := file.Write([]byte(code))
	if err != nil {
		return "", err
	}
	slog.Info("wrote", n, "bytes to file", filename)
	return filename, nil
}

func IsTempFileExists(filename string) bool {
	_, err := os.Stat(fmt.Sprintf("./resources/scripts/tmp/%s", filename))

	return err == nil
}
