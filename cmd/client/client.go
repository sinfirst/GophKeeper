package main

import (
	// ...
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sinfirst/GophKeeper/internal/models"
	"github.com/sinfirst/GophKeeper/internal/tui"
	pb "github.com/sinfirst/GophKeeper/proto/gophkeeper"
)

func main() {
	app := tui.NewAuthApp()
	if err := app.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
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
