package main

import (
	"flag"
	"fmt"
	"log"
	"context"
	"google.golang.org/grpc"

	pb "github.com/neerajrush/go-bingo/proto"
)

type BingoGame struct {
	GameGroupName string
	SecretPhrase string
}


var (
	serverAddr = flag.String("server-addr", "127.0.0.1:8000", "Game grpc server port")
)

func newBingoGame() *BingoGame {
	c := &BingoGame{GameGroupName: "Bingo", SecretPhrase: "test",}
	return c
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Connected to the grpc server.", *serverAddr)
	client := pb.NewGameClient(conn)
	bingo := newBingoGame()

	sessionResp, err := client.StartNewGame(context.Background(),
	                               &pb.StartSessionRequest{GameName: bingo.GameGroupName, 
				                               SecretPhrase: bingo.SecretPhrase,
		                       })
	if err != nil {
		log.Fatalf("failed to create new bingo session: %v", err)
	}

	fmt.Println("Successfully created bingo game with sessionId:", sessionResp.GetSessionId())

	defer conn.Close()
}
