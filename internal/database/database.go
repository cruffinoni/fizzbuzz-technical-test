package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cruffinoni/fizzbuzz/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	AddRequest(request *FizzBuzzRequest) error
	GetMostUsedRequest() (*MostUsedRequest, error)
}

type DB struct {
	instance *sqlx.DB
}

const (
	maxRetries = 5
	retryDelay = 10 * time.Second
)

func connectToDatabase(config *config.Database) (*sqlx.DB, error) {
	var (
		db        *sqlx.DB
		err       error
		retryTime = retryDelay
	)

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

		log.Printf("Failed to connect to the database, retrying in %s", retryTime)
		time.Sleep(retryTime)
		retryTime *= 2
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

func (db *DB) AddRequest(request *FizzBuzzRequest) error {
	err := db.instance.QueryRow(
		"INSERT INTO fizzbuzz.request (value_1, value_2, replace_1, replace_2, max) VALUES (?, ?, ?, ?, ?)",
		request.Int1,
		request.Int2,
		request.Str1,
		request.Str2,
		request.Limit,
	)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}

var ErrNoRequest = errors.New("no request found")

func (db *DB) GetMostUsedRequest() (*MostUsedRequest, error) {
	var request MostUsedRequest
	err := db.instance.QueryRow(
		"select fizzbuzz.request.value_1, fizzbuzz.request.value_2, count(*) as count from fizzbuzz.request group by fizzbuzz.request.value_1, fizzbuzz.request.value_2 order by count desc limit 1;",
	).Scan(&request.Int1, &request.Int2, &request.Hints)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRequest
		}
		return nil, err
	}
	return &request, nil
}
