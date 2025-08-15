package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sinfirst/GophKeeper/internal/config"
	"go.uber.org/zap"
)

// PGDB структура для хранения переменных
type PGDB struct {
	logger zap.SugaredLogger
	db     *pgxpool.Pool
}

// NewPGDB конструктор для структуры
func NewPGDB(config config.Config, logger zap.SugaredLogger) *PGDB {
	db, err := pgxpool.New(context.Background(), config.DatabaseDsn)

	if err != nil {
		logger.Errorw("Problem with connecting to db ", err)
		return nil
	}
	return &PGDB{logger: logger, db: db}
}

func (p *PGDB) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := p.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1
		)
	`, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking user existence: %w", err)
	}
	return exists, nil
}

func (p *PGDB) AddUserToDB(ctx context.Context, username, password string) error {
	var insertedUser string

	query := `
		INSERT INTO users (username, user_password)
		VALUES ($1, $2)
		ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
		RETURNING username
	`
	err := p.db.QueryRow(ctx, query, username, password).Scan(&insertedUser)

	if err != nil {
		return err
	}

	return nil
}
