package database

import (
	"fmt"
	"log"
	"time"

	"github.com/cruffinoni/fizzbuzz/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	AddRequest(request *FizzBuzzRequest) (int64, error)
}

type DB struct {
	instance *sqlx.DB
}

const (
	maxRetries = 5
	retryDelay = 5 * time.Second
)

func connectToDatabase(config *config.Database) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempt to connect to the database (try %d/%d)", i+1, maxRetries)
		db, err = sqlx.Open(
			"mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.Username, config.Password, config.Host, config.Port),
		)
		if err == nil {
			if err = db.Ping(); err == nil {
				log.Printf("Succesfully connected to the database")
				return db, nil
			}
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to the database after %d retries: %w", maxRetries, err)
}

func NewDB(config *config.Database) (*DB, error) {
	dbConnection, err := connectToDatabase(config)
	if err != nil {
		return nil, err
	}
	return &DB{instance: dbConnection}, nil
}

func (db *DB) AddRequest(request *FizzBuzzRequest) (int64, error) {
	var id int64
	err := db.instance.QueryRow(
		"INSERT INTO fizzbuzz.request (value_1, value_2, replace_1, replace_2, max) VALUES (?, ?, ?, ?, ?)",
		request.Int1,
		request.Int2,
		request.Str1,
		request.Str2,
		request.Limit,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
