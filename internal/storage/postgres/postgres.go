package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/errs"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

const maxRetries = 10
const retryDelay = 3 * time.Second

type Storage struct {
	db  *sql.DB
	log *zap.Logger
}

func NewStorage(postgresConf config.PostgresConfig, log *zap.Logger) (*Storage, error) {
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresConf.Host, postgresConf.Port, postgresConf.User, postgresConf.Password, postgresConf.DBName)

	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", url)
		if err != nil {
			log.Error("error opening connection to postgres", zap.Error(err))
			return nil, fmt.Errorf("error opening connection to postgres: %w", err)
		}

		err = db.Ping()
		if err == nil {
			break
		}

		log.Warn("failed to ping postgres, retrying...", zap.Error(err))
		time.Sleep(retryDelay)
	}

	if db == nil {
		return nil, errors.New("failed to connect to postgres")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres after %d retries: %w", maxRetries, err)
	}

	createTableStmt := `
    CREATE TABLE IF NOT EXISTS urlshortener (
        short_url TEXT NOT NULL PRIMARY KEY,
        url TEXT NOT NULL UNIQUE
    )`

	_, err = db.Exec(createTableStmt)
	if err != nil {
		return nil, fmt.Errorf("error executing create table statement: %w", err)
	}

	return &Storage{db: db, log: log}, nil
}

func (s *Storage) Put(url, shortURL string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	query := `INSERT INTO urlshortener (url, short_url) VALUES ($1, $2)`
	s.log.Info("storage.put", zap.String("url", url), zap.String("short-url", shortURL))

	_, err = tx.Exec(query, url, shortURL)
	if err != nil {
		_ = tx.Rollback()
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return errs.ErrURLIsExist
		}

		return fmt.Errorf("error executing insert statement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (s *Storage) Get(shortURL string) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("error starting transaction: %w", err)
	}

	query := `SELECT url FROM urlshortener WHERE short_url = $1`
	s.log.Info("storage.get", zap.String("short-url", shortURL))

	var url string
	err = tx.QueryRow(query, shortURL).Scan(&url)
	if err != nil {
		_ = tx.Rollback()

		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrURLIsNotExist
		}

		return "", fmt.Errorf("error scanning row: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("error committing transaction: %w", err)
	}

	return url, nil
}
