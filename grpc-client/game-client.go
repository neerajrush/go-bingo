package main

import (
	"flag"
	"fmt"
	"log"
	"errors"
	"context"
	"net"
	"google.golang.org/grpc"

	pb "github.com/neerajrush/go-bingo/proto"
)

type BingoClient struct {
	GroupName string
	SecretPhrase string
}

var (
	port = flag.Int("port", 8000, "Game server port")
)

func newGameClient() *BingoClient {
	c := &BingoClient{GameName: "Bingo", SecretPhrase: "test",}
	return c
}

func main() {
	flag.Parse()
	conn, err := net.Dial(fmt.Sprintf("tcp://localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
