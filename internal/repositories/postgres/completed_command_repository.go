package postgres

import (
	"fmt"
	"log/slog"
	"postgrespro-executor-service/internal/models"
	"postgrespro-executor-service/internal/repositories/redis"
)

func CreateCompletedCommand(storage *Storage, command models.Command) (int, error) {
	row, err := storage.Db.Query(`insert into completed_commands (command_id, result, completed_at, status) values ($1,'', now(),'RUNNING') returning id`, command.Id)
	if err != nil {
		slog.Error("error while saving new completed command", err)
		return 0, err
	}

	var id int
	row.Next()
	err = row.Scan(&id)
	if err != nil {
		slog.Error("error while parsing completed command id after creating", err)
		return 0, err
	}
	return id, nil
}

func SetPIDByCompletedCommandId(storage *Storage, id int, pid int) {
	err := storage.Db.QueryRow(`update completed_commands set ppid = $1 where id = $2`, pid, id)
	if err != nil {
		slog.Error("error when set command pid", err)
	}
}

func IsCommandStopped(storage *Storage, id int) bool {
	var status string
	err := storage.Db.QueryRow(`select status from completed_commands where id = $1`, id).Scan(&status)
	if err != nil {
		slog.Error("error when get command pid")
		return false
	}
	return status == "STOPPED"
}
func IsCommandRunning(storage *Storage, id int) bool {
	var status string
	err := storage.Db.QueryRow(`select status from completed_commands where id = $1`, id).Scan(&status)
	if err != nil {
		slog.Error("error when get command pid")
		return false
	}
	return status == "RUNNING"
}

func GetPPIDByCompletedCommandId(storage *Storage, id int) (int, error) {
	var pid int
	err := storage.Db.QueryRow(`select ppid from completed_commands where id = $1`, id).Scan(&pid)
	if err != nil {
		slog.Error("error when get command pid")
		return -1, err
	}
	return pid, nil
}

func UpdateCompletedCommandById(storage *Storage, id int, result string) error {
	_, err := storage.Db.Query(`update completed_commands set result = $2, completed_at = now(), status = 'COMPLETED' where id = $1`, id, result)
	if err != nil {
		slog.Error("error while updating  completed command", id, err)
		return err
	}
	return nil
}

func UpdateStoppedCommandById(storage *Storage, id int, result string) error {
	_, err := storage.Db.Query(`update completed_commands set result = $2, completed_at = now(), status = 'STOPPED' where id = $1`, id, result)
	if err != nil {
		slog.Error("error while updating  completed command", id, err)
		return err
	}
	slog.Info("set status stopped!")
	return nil
}

func UpdateCommandResultById(storage *Storage, id int, result string) error {
	_, err := storage.Db.Query(`update completed_commands set result = $2 where id = $1`, id, result)
	if err != nil {
		slog.Error("error while updating command result", id, err)
		return err
	}
	return nil
}

func FindCompletedCommandById(redisClient *redis.Redis, storage *Storage, id int) (models.CompletedCommandRequest, error) {

	var cc models.CompletedCommandRequest
	var c models.Command
	err := storage.Db.QueryRow(`select cc.id,  c.id, c.code, c.description, c.created_at, cc.completed_at, cc.result, cc.status 
										from completed_commands cc join commands c on c.id = cc.command_id 
										where cc.id = $1`, id).Scan(&cc.Id, &c.Id, &c.Code, &c.Description, &c.CreatedAt, &cc.CompletedAt, &cc.Result, &cc.Status)
	if err != nil {
		slog.Error("error while get command by id", id, err)
		return models.CompletedCommandRequest{}, err
	}

	fmt.Println("is command running: ", redis.IsCommandRunning(redisClient, id))
	if redis.IsCommandRunning(redisClient, id) {
		cc.Result, _ = redis.GetRunningCommandResultById(redisClient, id)
		cc.Status = "RUNNING"
	}

	cc.Command = c
	return cc, nil
}
