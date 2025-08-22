package handlers

import (
	"context"
	"fmt"

	"go.uber.org/zap"
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
	RetrieveDataFromDB(ctx context.Context, id int, username string) (models.Record, error)
}

type Handler struct {
	storage Storage
	config  config.Config
	logger  zap.SugaredLogger
}

func NewHandler(storage Storage, config config.Config, logger zap.SugaredLogger) Handler {
	handler := Handler{storage: storage, config: config, logger: logger}
	return handler
}

func (h *Handler) Register(ctx context.Context, login, password string) (string, error) {
	exist, err := h.storage.CheckUsernameExists(ctx, login)
	if err != nil {
		h.logger.Errorf("err: %v", err)
		return "", err
	}

	if exist {
		return "", fmt.Errorf("conflict")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorf("err: %v", err)
		return "", err
	}

	err = h.storage.AddUserToDB(ctx, login, string(hashedPassword))
	if err != nil {
		h.logger.Errorf("err: %v", err)
		return "", err
	}

	token, err := auth.BuildJWTString(login)
	if err != nil {
		h.logger.Errorf("err: %v", err)
		return "", err
	}

	return token, nil
}

func (h *Handler) Login(ctx context.Context, login, password string) (string, error) {
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
		h.logger.Errorf("err: %v", err)
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
	username, err := auth.CheckToken(token)
	if err != nil {
		return models.Record{}, fmt.Errorf("unauthenticated")
	}
	return h.storage.RetrieveDataFromDB(ctx, id, username)
}
