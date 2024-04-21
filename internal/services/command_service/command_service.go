package command_service

import (
	"bufio"
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
	dbCommand, err := postgres.FindCompletedCommandById(redisClient, storage, completedCommandId)
	if err != nil {
		slog.Error("error when get command result:", err)
		return models.CompletedCommandRequest{}, err
	}
	if !redis.IsCommandRunning(redisClient, completedCommandId) {
		redis.PutCompletedCommand(redisClient, dbCommand)
	}

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

	stdout, err := com.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = com.Start()
	postgres.SetPIDByCompletedCommandId(storage, completedCommandId, com.Process.Pid)
	fmt.Println("process pid", com.Process.Pid)
	if err != nil {
		slog.Error("stopping process error", err)
	}
	scanner := bufio.NewScanner(stdout)
	counter := 1
	for scanner.Scan() {
		m := scanner.Text()
		counter++
		err := redis.UpdateRunningCommandOutputById(redisClient, completedCommandId, m)
		if err != nil {
			slog.Error("error while updating redis data")
		}
		result, _ := redis.GetRunningCommandResultById(redisClient, completedCommandId)
		if counter%10 == 0 {
			postgres.UpdateCommandResultById(storage, completedCommandId, result)
		}
		fmt.Println(filename, m)
	}
	com.Wait()

	result, err := redis.GetRunningCommandResultById(redisClient, completedCommandId)
	if err != nil {
		slog.Error("error", err)
	}
	if !postgres.IsCommandStopped(storage, completedCommandId) {
		postgres.UpdateCompletedCommandById(storage, completedCommandId, result)
	}
	redis.DeleteRunningCommandResultById(redisClient, completedCommandId)

	//slog.Info(string(out))
}

func StopCommandById(storage *postgres.Storage, redisClient *redis.Redis, commandId int) (models.CompletedCommandRequest, error) {
	if !postgres.IsCommandRunning(storage, commandId) {
		slog.Error("error stopping command. status not running")
		return models.CompletedCommandRequest{}, fmt.Errorf("command %d status not running", commandId)
	}

	ppid, err := postgres.GetPPIDByCompletedCommandId(storage, commandId)
	if err != nil {
		return models.CompletedCommandRequest{}, err
	}
	com := exec.Command("pkill", "-P", strconv.Itoa(ppid))
	err = com.Start()
	if err != nil {
		return models.CompletedCommandRequest{}, err
	}

	result, err := redis.GetRunningCommandResultById(redisClient, commandId)
	if err != nil {
		slog.Error("error", err)
	}
	slog.Info("stopped command result:", result)
	err = postgres.UpdateStoppedCommandById(storage, commandId, result)
	if err != nil {
		return models.CompletedCommandRequest{}, err
	}
	redis.DeleteRunningCommandResultById(redisClient, commandId)
	completed, err := GetCompletedCommandById(storage, redisClient, commandId)
	return completed, err
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
