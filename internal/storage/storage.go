package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sinfirst/GophKeeper/internal/config"
	"github.com/sinfirst/GophKeeper/internal/models"
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
	var id int
	query := `INSERT INTO records (type_record, user_data, meta, username)
				VALUES ($1, $2, $3, $4)`
	_, err := p.db.Exec(ctx, query, record.TypeRecord, record.Data, record.Meta, username)
	if err != nil {
		return 0, err
	}
	query = `SELECT id FROM records WHERE username = $1`
	row := p.db.QueryRow(ctx, query, username)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PGDB) GetUserByDataID(ctx context.Context, id int) (string, error) {
	var username string
	query := `SELECT username FROM records WHERE id = $1`
	row := p.db.QueryRow(ctx, query, id)
	err := row.Scan(&username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("not found")
		}
		return "", err
	}
	return username, nil

}
func (p *PGDB) RetrieveDataFromDB(ctx context.Context, id int) (models.Record, error) {
	var record models.Record

	query := `SELECT type_record, user_data, meta FROM records WHERE id = $1`
	row := p.db.QueryRow(ctx, query, id)
	err := row.Scan(&record.TypeRecord, &record.Data, &record.Meta)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Record{}, fmt.Errorf("not found")
		}
		return models.Record{}, err
	}
	return record, nil
}

func (p *PGDB) UpdateDataInDB(ctx context.Context, id int, meta string, data []byte) error {
	query := `UPDATE records SET user_data = $1, meta = $2 
			WHERE id = $3`
	result, err := p.db.Exec(ctx, query, data, meta, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

func (p *PGDB) GetListData(ctx context.Context, username string) ([]models.Record, error) {
	var records []models.Record
	query := `SELECT id, type_record, user_data, meta 
			FROM records WHERE username = $1`
	rows, err := p.db.Query(ctx, query, username)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record models.Record

		err := rows.Scan(&record.Id, &record.TypeRecord, &record.Data, &record.Meta)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func (p *PGDB) CheckRecordExist(ctx context.Context, id int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM records WHERE id = $1
		)`
	err := p.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking existence: %w", err)
	}
	return exists, nil
}
func (p *PGDB) DeleteDataFromDB(ctx context.Context, id int) error {
	query := `DELETE FROM records
				WHERE id = $1`

	_, err := p.db.Exec(ctx, query, id)

	if err != nil {
		p.logger.Errorw("Problem with deleting from db: ", err)
		return err
	}
	return nil
}

// InitMigrations инициализация миграций
func InitMigrations(conf config.Config, logger zap.SugaredLogger) error {
	logger.Infow("Start migrations")
	db, err := sql.Open("pgx", conf.DatabaseDsn)

	if err != nil {
		logger.Errorw("Error with connection to DB: ", err)
		return err
	}

	defer db.Close()

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "migrations")

	err = goose.Up(db, migrationsPath)
	if err != nil {
		logger.Errorw("Error with migrations: ", err)
		return err
	}
	return nil
}
