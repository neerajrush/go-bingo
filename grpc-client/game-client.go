package main

import (
	"flag"
	"fmt"
	"log"
	"context"
	"bufio"
	"os"
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

func startNewGame(ctx context.Context, client pb.GameClient) *pb.StartSessionResponse {
	bingo := BingoGame{}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter GameGroupName:")
	scanner.Scan()
	bingo.GameGroupName = scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	fmt.Print("Enter SecretPhrase:")
	scanner.Scan()
	bingo.SecretPhrase = scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	sessionResp, err := client.StartNewGame(context.Background(),
	                               &pb.StartSessionRequest{GameName: bingo.GameGroupName, 
				                               SecretPhrase: bingo.SecretPhrase,
		                       })
	if err != nil {
		log.Fatalf("failed to create new bingo session: %v", err)
	}

	fmt.Println("Successfully created bingo game with sessionId:", sessionResp.GetSessionId())

	return sessionResp
}

func getGameLink(ctx context.Context, client pb.GameClient, sessionId string) string {
	sessionResp, err := client.GetGameLink(context.Background(),
	                               &pb.GetSessionRequest{SessionId: sessionId,
		                       })
	if err != nil {
		log.Fatalf("failed to get session link: %v", err)
	}

	fmt.Println("Successfully retrieved game link for sessionId:", sessionResp.GetSessionId())

	return sessionResp.GetGameLink()
}

func addPlayer(ctx context.Context, client pb.GameClient, sessionId string) [][]int32 {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter PlayerName:")
	scanner.Scan()
	playerName := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	addPlayerResp, err := client.AddAPlayer(context.Background(),
	                               &pb.AddPlayerRequest{SessionId: sessionId,
							    Name: playerName,
		                       })
	if err != nil {
		log.Fatalf("failed to add new player to the game session: %v", err)
	}

	fmt.Println("Successfully added a new player (", playerName , ") to sessionId:", sessionId)

	fmt.Println("SheetSize:", addPlayerResp.GetSheetSize())

	sheet := make([][]int32, addPlayerResp.GetSheetSize())

	svrSheet := addPlayerResp.GetSheet()

	for i, rows := range svrSheet {
		sheet[i] = make([]int32, addPlayerResp.GetSheetSize())
		for j, val := range rows.Cols {
			sheet[i][j] = val
		}
	}

	return sheet
}

/*
	AddAPlayer(ctx context.Context, in *AddPlayerRequest, opts ...grpc.CallOption) (*AddPlayerResponse, error)
	ListPlayers(ctx context.Context, in *PlayersListRequest, opts ...grpc.CallOption) (*PlayersListResponse, error)
	EnablePlayer(ctx context.Context, in *EnablePlayerRequest, opts ...grpc.CallOption) (*EnablePlayerResponse, error)
	ApplyRules(ctx context.Context, in *RulesListRequest, opts ...grpc.CallOption) (*RulesListResponse, error)
	DrawANumber(ctx context.Context, in *DrawNumberRequest, opts ...grpc.CallOption) (*DrawNumberResponse, error)
	AttachToDraws(ctx context.Context, in *AttachRequest, opts ...grpc.CallOption) (Game_AttachToDrawsClient, error)
	DrawnNumbersList(ctx context.Context, in *DrawnNumbersListRequest, opts ...grpc.CallOption) (*DrawnNumbersResponse, error)
	AnnounceWinners(ctx context.Context, in *AnnounceWinnersRequest, opts ...grpc.CallOption) (Game_AnnounceWinnersClient, error)
	StopGame(ctx context.Context, in *StopSessionRequest, opts ...grpc.CallOption) (*StopSessionResponse, error)
	*/

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

	startSessionResp := startNewGame(context.Background(), client)
	log.Println(startSessionResp.GetSessionId())

	gameLink := getGameLink(context.Background(), client, startSessionResp.GetSessionId())
	log.Println(gameLink)

	sheet := addPlayer(context.Background(), client, startSessionResp.GetSessionId())
	log.Println(sheet)

	defer conn.Close()
}
