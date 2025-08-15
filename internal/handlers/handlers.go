package handlers

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/sinfirst/GophKeeper/internal/config"
	"github.com/sinfirst/GophKeeper/internal/middleware/auth"
)

type Storage interface {
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	AddUserToDB(ctx context.Context, username, password string) error
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
	return "", nil
}
