package main

import (
	"postgrespro-executor-service/internal/models"
	"postgrespro-executor-service/internal/repositories/redis"
	"postgrespro-executor-service/internal/utils/error_check"
	"testing"
	"time"
)

//func TestRouterFunctionality()  {
//
//}

func TestRedisConnect(t *testing.T) {
	redisClient, err := redis.New()
	error_check.CheckError(err)

	command := models.Command{
		Id:          0,
		Code:        "data",
		Description: "test",
		CreatedAt:   time.Time{},
	}

	err = redis.PutCommand(redisClient, command)
	if err != nil {
		t.Error("error when put data in redis", err)
	}

	gotCommand, err := redis.GetCommand(redisClient, command.Id)
	if err != nil {
		t.Error("error when get data in redis", err)
	}

	if command != gotCommand {
		t.Error("error after comparing got command", err)
	}
}

func TestPostgresConnect(t *testing.T) {

}
