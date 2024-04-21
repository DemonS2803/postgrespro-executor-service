package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"postgrespro-executor-service/internal/models"
	"strconv"
	"time"
)

func PutRequest(client *Redis, requestId string, commandId string) error {
	ctx := context.Background()

	err := client.Client.Set(ctx, fmt.Sprintf("request %s", requestId), commandId, 5*time.Minute).Err()
	if err != nil {
		slog.Error("error while saving requestID-completedId to redis")
		return err
	}
	return nil
}

func GetCompletedIdByRequestId(redisClient *Redis, requestId string) int {
	ctx := context.Background()
	req := redisClient.Client.Get(ctx, fmt.Sprintf("request %s", requestId))
	if err := req.Err(); err != nil {
		slog.Info("unable to GET data. error: %v", err)
		return -1
	}

	res, err := req.Result()
	if err != nil {
		slog.Info("unable to GET data. error: %v", err)
		return -1
	}

	id, err := strconv.Atoi(res)
	if err != nil {
		slog.Error("error while parsing request id", err)
		return -1
	}
	return id
}

func HasCompletedIdForRequest(redisClient *Redis, requestId string) bool {
	ctx := context.Background()
	req := redisClient.Client.Get(ctx, fmt.Sprintf("request %s", requestId))
	if err := req.Err(); err != nil {
		return false
	}

	return true
}

func GetCommand(client *Redis, commandId int) (models.Command, error) {
	ctx := context.Background()
	req := client.Client.Get(ctx, fmt.Sprintf("command %d", commandId))
	err := req.Err()
	if err != nil {
		slog.Error("error while get command", err)
		return models.Command{}, err
	}
	data, _ := req.Result()
	var command models.Command
	err = json.Unmarshal([]byte(data), &command)

	return command, err
}

func PutCommand(client *Redis, command models.Command) error {
	ctx := context.Background()
	data, err := json.Marshal(command)
	if err != nil {
		slog.Error("error when put command in redis", err)
		return err
	}

	req := client.Client.Set(ctx, fmt.Sprintf("command %d", command.Id), data, 5*time.Minute)
	fmt.Println(req.Err())

	return req.Err()
}

func GetCompletedCommand(client *Redis, completedCommandId int) (models.CompletedCommandRequest, error) {
	ctx := context.Background()
	req := client.Client.Get(ctx, fmt.Sprintf("result %d", completedCommandId))
	err := req.Err()
	if err != nil {
		slog.Error("error while get command result", err)
		return models.CompletedCommandRequest{}, err
	}
	data, _ := req.Result()
	var command models.CompletedCommandRequest
	err = json.Unmarshal([]byte(data), &command)

	return command, err
}

func PutCompletedCommand(client *Redis, command models.CompletedCommandRequest) error {
	ctx := context.Background()
	data, err := json.Marshal(command)
	if err != nil {
		slog.Error("error when put command result in redis", err)
		return err
	}

	req := client.Client.Set(ctx, fmt.Sprintf("result %d", command.Id), data, 5*time.Minute)
	fmt.Println(req.Err())

	return req.Err()
}

func GetCommandsList(client *Redis, limit int, offset int) ([]models.Command, error) {
	ctx := context.Background()
	req := client.Client.Get(ctx, fmt.Sprintf("commands list %d %d", limit, offset))
	err := req.Err()
	if err != nil {
		slog.Error("error while get commands list", err)
		return nil, err
	}
	data, _ := req.Result()
	var commands []models.Command
	err = json.Unmarshal([]byte(data), &commands)

	return commands, err
}

func PutCommandsList(client *Redis, limit int, offset int, command []models.Command) error {
	ctx := context.Background()
	data, err := json.Marshal(command)
	if err != nil {
		slog.Error("error when put commands list in redis", err)
		return err
	}

	req := client.Client.Set(ctx, fmt.Sprintf("commands list %d %d", limit, offset), data, 5*time.Minute)
	fmt.Println(req.Err())

	return req.Err()
}

func UpdateRunningCommandOutputById(client *Redis, commandId int, result string) error {
	ctx := context.Background()

	client.Client.Append(ctx, fmt.Sprintf("running %d output", commandId), result)
	return nil
}

func GetRunningCommandResultById(client *Redis, commandId int) (string, error) {
	ctx := context.Background()
	req := client.Client.Get(ctx, fmt.Sprintf("running %d output", commandId))
	return req.Result()
}

func DeleteRunningCommandResultById(client *Redis, commandId int) {
	ctx := context.Background()
	client.Client.Del(ctx, fmt.Sprintf("running %d output", commandId))
}

func IsCommandRunning(client *Redis, commandId int) bool {
	ctx := context.Background()
	req := client.Client.Exists(ctx, fmt.Sprintf("running %d output", commandId))
	return req.Val() == 1
}
