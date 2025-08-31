package handlers

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/sinfirst/GophKeeper/internal/config"
	"github.com/sinfirst/GophKeeper/internal/middleware/auth"
	"github.com/sinfirst/GophKeeper/internal/models"
)

type Storage interface {
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	AddUserToDB(ctx context.Context, username, password string) error
	GetUserPassword(ctx context.Context, username string) (string, error)
	StoreDataToDB(ctx context.Context, record models.Record, username string) (int, error)
	RetrieveDataFromDB(ctx context.Context, id int) (models.Record, error)
	GetUserByDataID(ctx context.Context, id int) (string, error)
	UpdateDataInDB(ctx context.Context, id int, meta string, data []byte) error
	GetListData(ctx context.Context, username string) ([]models.Record, error)
	CheckRecordExist(ctx context.Context, id int) (bool, error)
	DeleteDataFromDB(ctx context.Context, id int) error
}

type Handler struct {
	storage Storage
	config  config.Config
}

func NewHandler(storage Storage, config config.Config) Handler {
	handler := Handler{storage: storage, config: config}
	return handler
}

func (h *Handler) Register(ctx context.Context, login, password string) (string, error) {
	exist, err := h.storage.CheckUsernameExists(ctx, login)
	if err != nil {
		return "", err
	}

	if exist {
		return "", fmt.Errorf("conflict")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	err = h.storage.AddUserToDB(ctx, login, string(hashedPassword))
	if err != nil {
		return "", err
	}

	token, err := auth.BuildJWTString(login)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *Handler) Login(ctx context.Context, login, password string) (string, error) {
	exist, err := h.storage.CheckUsernameExists(ctx, login)
	if err != nil {
		return "", err
	}

	if exist {
		return "", fmt.Errorf("not found")
	}

	passwordFromBD, err := h.storage.GetUserPassword(ctx, login)
	if err != nil {
		return "", fmt.Errorf("unauthenticated")
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordFromBD), []byte(password))
	if err != nil {
		return "", fmt.Errorf("unauthenticated")
	}
	token, err := auth.BuildJWTString(login)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *Handler) StoreData(ctx context.Context, token string, record models.Record) (int, error) {
	username, err := auth.CheckToken(token)
	if err != nil {
		return 0, fmt.Errorf("unauthenticated")
	}
	return h.storage.StoreDataToDB(ctx, record, username)
}

func (h *Handler) RetrieveData(ctx context.Context, token string, id int) (models.Record, error) {
	_, err := h.checkAccess(ctx, token, id)
	if err != nil {
		return models.Record{}, err
	}
	return h.storage.RetrieveDataFromDB(ctx, id)
}

func (h *Handler) UpdateData(ctx context.Context, token, meta string, id int, data []byte) error {
	_, err := h.checkAccess(ctx, token, id)
	if err != nil {
		return err
	}
	return h.storage.UpdateDataInDB(ctx, id, meta, data)

}

func (h *Handler) ListData(ctx context.Context, token string) ([]models.Record, error) {
	username, err := auth.CheckToken(token)
	if err != nil {
		return nil, fmt.Errorf("unauthenticated")
	}
	return h.storage.GetListData(ctx, username)
}

func (h *Handler) DeleteData(ctx context.Context, token string, id int) error {
	_, err := h.checkAccess(ctx, token, id)
	if err != nil {
		return err
	}
	exist, err := h.storage.CheckRecordExist(ctx, id)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("not found")
	}
	return h.storage.DeleteDataFromDB(ctx, id)
}

func (h *Handler) checkAccess(ctx context.Context, token string, id int) (string, error) {
	username, err := auth.CheckToken(token)
	if err != nil {
		return username, fmt.Errorf("unauthenticated")
	}

	usernameFromBD, err := h.storage.GetUserByDataID(ctx, id)
	if username != usernameFromBD {
		return username, fmt.Errorf("access denied")
	}
	if err != nil {
		return "", err
	}
	return username, nil
}
