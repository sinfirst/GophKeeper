package main

import (
	// ...
	"context"
	"encoding/json"
	"fmt"

	"github.com/sinfirst/GophKeeper/internal/config"
	"github.com/sinfirst/GophKeeper/internal/middleware/logging"
	"github.com/sinfirst/GophKeeper/internal/models"
	pb "github.com/sinfirst/GophKeeper/proto/gophkeeper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := config.NewConfig()
	logger := logging.NewLogger()
	conn, err := grpc.NewClient(config.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed to up client: %v", err)
	}
	defer conn.Close()

	c := pb.NewGophKeeperClient(conn)

	TestFunc(c)
}

func TestFunc(c pb.GophKeeperClient) {
	token, err := c.Register(context.Background(), &pb.AuthRequest{Username: "sinfirst2", Password: "qwerty12345"})
	if err != nil {
		fmt.Println(err, "1")
		fmt.Println("toke")
	}
	fmt.Println(token, err)
	req := models.LoginJSON{Login: "LoginForhhfdgReq", Password: "PasswordFogrReq"}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err, "1")
		fmt.Println("marshal")
	}
	id, err := c.StoreData(context.Background(), &pb.StoreRequest{Token: token.Token, Record: &pb.DataRecord{Type: "LOGIN", Data: jsonReq, Meta: "Test login req"}})
	if err != nil {
		fmt.Println(err, "1")
		fmt.Println("store")
	}
	record, err := c.RetrieveData(context.Background(), &pb.RetrieveRequest{Token: token.Token, Id: id.Id})
	if err != nil {
		fmt.Println(err, "2")
		fmt.Println("record")
	}
	fmt.Println(record.Record.Id, record.Record.Type, string(record.Record.Data), record.Record.Meta)
}
