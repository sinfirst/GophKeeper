package client

import (
	"context"
	"log"

	pb "github.com/sinfirst/GophKeeper/proto/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client представляет gRPC клиент для аутентификации
type Client struct {
	conn   *grpc.ClientConn
	client pb.GophKeeperClient
	token  string
}

// Response ответ сервера
type Response struct {
	Success bool
	Message string
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
func (c *Client) Register(ctx context.Context, username, password string) (*Response, error) {
	resp, _ := c.client.Register(ctx, &pb.AuthRequest{Username: username, Password: password})
	c.token = resp.Token

	return nil, nil
}

// Close закрывает соединение
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
