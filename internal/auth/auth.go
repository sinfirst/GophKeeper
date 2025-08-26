package auth

import (
	"context"
	"fmt"
	"log"

	pb "github.com/sinfirst/GophKeeper/proto/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client представляет gRPC клиент для аутентификации
type Client struct {
	conn       *grpc.ClientConn
	authClient pb.GophKeeperClient
}

// RegisterResponse ответ регистрации
type RegisterResponse struct {
	Success bool
	UserId  string
	Message string
}

// NewClient создает новый gRPC клиент
func NewClient(serverAddr string) *Client {
	conn, err := grpc.NewClient(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Не удалось подключиться: %v", err)
		return &Client{}
	}

	client := pb.NewGophKeeperClient(conn)

	return &Client{
		conn:       conn,
		authClient: client,
	}
}

// Register отправляет запрос на регистрацию
func (c *Client) Register(ctx context.Context, username, password string) (*RegisterResponse, error) {
	// Заглушка - здесь будет реальный gRPC вызов
	fmt.Printf("Регистрация: %s/%s\n", username, password)

	// Пример реального вызова:
	// req := &proto.RegisterRequest{
	//     Username: username,
	//     Password: password,
	// }
	// resp, err := c.authClient.Register(ctx, req)

	// Имитация ответа от сервера
	return &RegisterResponse{
		Success: true,
		UserId:  "user-123",
		Message: "Registration successful",
	}, nil
}

// Close закрывает соединение
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
