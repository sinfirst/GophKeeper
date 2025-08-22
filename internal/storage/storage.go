package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sinfirst/GophKeeper/internal/config"
	"github.com/sinfirst/GophKeeper/internal/models"
	"go.uber.org/zap"
)

// PGDB структура для хранения переменных
type PGDB struct {
	logger zap.SugaredLogger
	db     *pgxpool.Pool
	idData int
}

// NewPGDB конструктор для структуры
func NewPGDB(config config.Config, logger zap.SugaredLogger) *PGDB {
	db, err := pgxpool.New(context.Background(), config.DatabaseDsn)

	if err != nil {
		logger.Errorw("Problem with connecting to db ", err)
		return nil
	}
	return &PGDB{logger: logger, db: db, idData: 1}
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

func (p *PGDB) GetUserPassword(ctx context.Context, username string) (string, error) {
	var password string

	query := `SELECT user_password FROM users WHERE username = $1`
	row := p.db.QueryRow(ctx, query, username)
	err := row.Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (p *PGDB) StoreDataToDB(ctx context.Context, record models.Record, username string) (int, error) {
	query := `INSERT INTO records (id, type_record, user_data, meta, username)
				VALUES ($1, $2, $3, $4, $5)`
	_, err := p.db.Exec(ctx, query, p.idData, record.TypeRecord, record.Data, record.Meta, username)
	if err != nil {
		return p.idData, err
	}
	p.idData++
	return p.idData - 1, nil
}

func (p *PGDB) RetrieveDataFromDB(ctx context.Context, id int, username string) (models.Record, error) {
	var record models.Record
	var usernameFromBD string

	query := `SELECT type_record, user_data, meta, username FROM records WHERE id = $1`
	row := p.db.QueryRow(ctx, query, id)
	err := row.Scan(&record.TypeRecord, &record.Data, &record.Meta, &username)
	if usernameFromBD != username {
		return models.Record{}, fmt.Errorf("access denied")
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Record{}, fmt.Errorf("not found")
		}
		return models.Record{}, err
	}
	return record, nil
}
