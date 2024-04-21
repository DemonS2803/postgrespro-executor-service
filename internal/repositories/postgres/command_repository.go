package postgres

import (
	"log/slog"
	"postgrespro-executor-service/internal/models"
)

func CountCommands(storage *Storage) int {
	var count int
	storage.Db.QueryRow(`select count(*) from commands`).Scan(&count)
	return count

}

func FindCommandById(storage *Storage, id int) (models.Command, error) {
	var com models.Command
	err := storage.Db.QueryRow(`select id, code, description, created_at from commands where id = $1`, id).Scan(&com.Id, &com.Code, &com.Description, &com.CreatedAt)
	if err != nil {
		slog.Error("error while get command by id", id, err)
		return models.Command{}, err
	}

	return com, nil
}

func FindCommandListByLimitAndOffset(storage *Storage, limit int, offset int) ([]models.Command, error) {
	offset = limit * offset
	rows, err := storage.Db.Query(`SELECT c.id, c.code, c.description, c.created_at
				FROM commands c
				ORDER BY created_at
				limit $1 
				offset $2;`, limit, offset)
	if err != nil {
		slog.Error("error while get commands list", err)
		return nil, err
	}

	var list []models.Command
	for rows.Next() {
		var command models.Command
		rows.Scan(&command.Id, &command.Code, &command.Description, &command.CreatedAt)
		list = append(list, command)
	}
	return list, nil
}

func CreateCommandAndReturn(storage *Storage, createCommand models.CreateCommandRequest) (models.Command, error) {
	var command models.Command
	err := storage.Db.QueryRow(`INSERT into commands (code, description, created_at) 
						values ($1, $2, now()) returning id, code, description, created_at`, createCommand.Code, createCommand.Description).Scan(&command.Id, &command.Code, &command.Description, &command.CreatedAt)

	if err != nil {
		slog.Error("error when insert new command", err)
		return models.Command{}, err
	}
	return command, nil
}
