package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log/slog"
	"os"
	"postgrespro-executor-service/internal/utils/error_check"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "main_db"
)

type Storage struct {
	Db *sql.DB
}

func New() (*Storage, error) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	slog.Info("ready to connect psql")
	// open database
	db, err := sql.Open("postgres", psqlconn)
	error_check.CheckError(err)

	// check db
	err = db.Ping()
	error_check.CheckError(err)

	storage := &Storage{Db: db}

	slog.Info("postgres success connected")

	// create tables
	initQuery := readInitFile()
	_, err = db.Query(initQuery)
	error_check.CheckError(err)

	slog.Info("success read init sql file")

	// setup data
	count := CountCommands(storage)
	if count == 0 {
		fillQuery := readFillFile()
		_, err = db.Query(fillQuery)
		if err != nil {
			slog.Warn("error while fill test data. continue with that data or rebuild postgres container")

		}
	}

	return storage, nil
}

func readInitFile() string {
	file, err := os.Open("resources/database/init.sql")
	slog.Info("in read init", err)
	error_check.CheckError(err)
	b, err := io.ReadAll(file)
	error_check.CheckError(err)

	return string(b)
}

func readFillFile() string {
	file, err := os.Open("resources/database/fill.sql")
	error_check.CheckError(err)
	b, err := io.ReadAll(file)
	error_check.CheckError(err)
	return string(b)
}
