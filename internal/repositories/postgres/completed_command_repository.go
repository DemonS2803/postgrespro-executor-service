package postgres

import (
	"log/slog"
	"postgrespro-executor-service/internal/models"
)

func CreateCompletedCommand(storage *Storage, command models.Command) (int, error) {
	row, err := storage.Db.Query(`insert into completed_commands (command_id, result, completed_at) values ($1,null, null) returning id`, command.Id)
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

func UpdateCompletedCommandById(storage *Storage, id int, result string) error {
	_, err := storage.Db.Query(`update completed_commands set result = $2, completed_at = now() where id = $1`, id, result)
	if err != nil {
		slog.Error("error while updating  completed command", id, err)
		return err
	}
	return nil
}

func FindCompletedCommandById(storage *Storage, id int) (models.CompletedCommandRequest, error) {
	var cc models.CompletedCommandRequest
	var c models.Command
	err := storage.Db.QueryRow(`select cc.id,  c.id, c.code, c.description, c.created_at, cc.completed_at, cc.result 
										from completed_commands cc join commands c on c.id = cc.command_id 
										where cc.id = $1`, id).Scan(&cc.Id, &c.Id, &c.Code, &c.Description, &c.CreatedAt, &cc.CompletedAt, &cc.Result)
	if err != nil {
		slog.Error("error while get command by id", id, err)
		return models.CompletedCommandRequest{}, err
	}

	cc.Command = c
	return cc, nil
}
