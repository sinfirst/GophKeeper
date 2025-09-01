package client

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/sinfirst/GophKeeper/internal/models"
	pb "github.com/sinfirst/GophKeeper/proto/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Client представляет gRPC клиент для аутентификации
type Client struct {
	conn   *grpc.ClientConn
	client pb.GophKeeperClient
	token  string
}

// NewClient создает новый gRPC клиент
func NewClient(serverAddr string) *Client {
	conn, err := grpc.NewClient(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Can't connect with server: %v", err)
	}

	client := pb.NewGophKeeperClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}
}

// Register отправляет запрос на регистрацию
func (c *Client) Register(ctx context.Context, username, password string) error {
	resp, err := c.client.Register(ctx, &pb.AuthRequest{Username: username, Password: password})
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.AlreadyExists:
			return fmt.Errorf("пользователь с таким логином уже существует")
		case codes.Internal:
			return fmt.Errorf("ошибка сервера")
		}
	}
	c.token = resp.Token
	return nil
}

func (c *Client) Login(ctx context.Context, username, password string) error {
	resp, err := c.client.Login(ctx, &pb.AuthRequest{Username: username, Password: password})
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.Unauthenticated:
			return fmt.Errorf("неверный пароль")
		case codes.NotFound:
			return fmt.Errorf("пользователь с таким логином не найден")
		case codes.Internal:
			return fmt.Errorf("ошибка сервера")
		}
	}
	c.token = resp.Token
	return nil
}

func (c *Client) StoreData(ctx context.Context, typeRecord, meta string, data []byte) (int, error) {
	record := &pb.DataRecord{Type: typeRecord, Data: data, Meta: meta}
	resp, err := c.client.StoreData(ctx, &pb.StoreRequest{Token: c.token, Record: record})
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.Unauthenticated:
			return 0, fmt.Errorf("войдите в аккаунт перед выполнением запроса")
		case codes.Internal:
			return 0, fmt.Errorf("ошибка сервера")
		}
	}
	return int(resp.Id), nil
}

func (c *Client) RetrieveData(ctx context.Context, id string) (models.Record, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return models.Record{}, fmt.Errorf("введите число")
	}

	resp, err := c.client.RetrieveData(ctx, &pb.RetrieveRequest{Token: c.token, Id: intID})
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.Unauthenticated:
			return models.Record{}, fmt.Errorf("войдите в аккаунт перед выполнением запроса")
		case codes.PermissionDenied:
			return models.Record{}, fmt.Errorf("в доступе отказано")
		case codes.NotFound:
			return models.Record{}, fmt.Errorf("данные с таким id не найдены")
		case codes.Internal:
			return models.Record{}, fmt.Errorf("ошибка сервера")
		}
	}
	return models.Record{Id: int(resp.Record.Id), TypeRecord: resp.Record.Type, Data: resp.Record.Data, Meta: resp.Record.Meta}, nil
}

func (c *Client) UpdateData(ctx context.Context, id, meta string, data []byte) error {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("введите число")
	}

	_, err = c.client.UpdateData(ctx, &pb.UpdateResponse{Token: c.token, Id: intID, Meta: meta, Data: data})
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.Unauthenticated:
			return fmt.Errorf("войдите в аккаунт перед выполнением запроса")
		case codes.PermissionDenied:
			return fmt.Errorf("в доступе отказано")
		case codes.NotFound:
			return fmt.Errorf("данные с таким id не найдены")
		case codes.Internal:
			return fmt.Errorf("ошибка сервера")
		}
	}
	return nil
}

func (c *Client) ListData(ctx context.Context) ([]models.Record, error) {
	var records []models.Record

	resp, err := c.client.ListData(ctx, &pb.ListRequest{Token: c.token})
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.Unauthenticated:
			return nil, fmt.Errorf("войдите в аккаунт перед выполнением запроса")
		case codes.PermissionDenied:
			return nil, fmt.Errorf("в доступе отказано")
		case codes.NotFound:
			return nil, fmt.Errorf("данные с таким id не найдены")
		case codes.Internal:
			return nil, fmt.Errorf("ошибка сервера")
		}
	}
	for _, i := range resp.Records {
		records = append(records, models.Record{
			Id:         int(i.Id),
			TypeRecord: i.Type,
			Data:       i.Data,
			Meta:       i.Meta,
		})
	}
	return records, nil
}

func (c *Client) DeleteData(ctx context.Context, id string) error {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("введите число")
	}

	_, err = c.client.DeleteData(ctx, &pb.DeleteRequest{Token: c.token, Id: intID})
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.Unauthenticated:
			return fmt.Errorf("войдите в аккаунт перед выполнением запроса")
		case codes.PermissionDenied:
			return fmt.Errorf("в доступе отказано")
		case codes.NotFound:
			return fmt.Errorf("данные с таким id не найдены")
		case codes.Internal:
			return fmt.Errorf("ошибка сервера")
		}
	}
	return nil
}

func (c *Client) GetVersion(ctx context.Context) (models.VersionBuild, error) {
	resp, err := c.client.GetVersion(ctx, &emptypb.Empty{})
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.Internal:
			return models.VersionBuild{}, fmt.Errorf("ошибка сервера")
		}
	}
	return models.VersionBuild{Version: resp.Ver.Version, Date: resp.Ver.Date}, nil
}

// Close закрывает соединение
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
