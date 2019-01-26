package main

import (
	"flag"
	"fmt"
	"log"
	"context"
	"bufio"
	"os"
	"io"
	"strconv"
	"sync"
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

func listPlayers(ctx context.Context, client pb.GameClient, sessionId string)[]string {
	listPlayersResp, err := client.ListPlayers(context.Background(),
	                               &pb.PlayersListRequest{SessionId: sessionId,
		                       })
	if err != nil {
		log.Fatalf("failed to list all players for session: %v", err)
	}

	fmt.Println("Successfully pulled list of players for sessionId:", sessionId)

	playersList := listPlayersResp.GetPlayers()

	fmt.Println("ListSize:", len(playersList))

	return playersList
}

func enablePlayer(ctx context.Context, client pb.GameClient, sessionId string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter PlayerName to enable:")
	scanner.Scan()
	playerName := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	enablePlayerResp, err := client.EnablePlayer(context.Background(),
	                               &pb.EnablePlayerRequest{SessionId: sessionId,
		                       })
	if err != nil {
		log.Fatalf("failed to enable the player(%v) for session: %v", playerName, err)
	}

	fmt.Println("Successfully enabled the player (", playerName, ") for sessionId:", sessionId)

	fmt.Println("Enabled:", enablePlayerResp.GetPlayerEnabled())
}

func applyRules(ctx context.Context, client pb.GameClient, sessionId string) {
	rules := make([]pb.RulesType, 0)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("The Game Rules are:")
	for {
		fmt.Println("1: Single Row Match")
		fmt.Println("2: Single Col Match")
		fmt.Println("3: Single Diagonal Match")
		fmt.Println("4: Two Rows Match")
		fmt.Println("5: Two Cols Match")
		fmt.Println("6: Two Diagonal Match")
		fmt.Println("7: Full House")
		fmt.Println("0: Quit")
		fmt.Print("Enter Rules to enable:")
		scanner.Scan()
		ruleNo := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
		ruleNoInt, err := strconv.Atoi(ruleNo)
		if err != nil {
			log.Println(err)
			continue
		}
		if ruleNoInt < 0 && ruleNoInt > 7 {
			log.Println("Invalid selection (choose again)")
			continue
		}
		if ruleNoInt == 0 {
			break
		}
		rules = append(rules, pb.RulesType(ruleNoInt))
	}

	if len(rules) == 0 {
		return
	}

	rulesListResp, err := client.ApplyRules(context.Background(),
					&pb.RulesListRequest{SessionId: sessionId,
							     Rules: rules,
					})
	if err != nil {
		log.Fatalf("failed to apply the rules list for session: %v", err)
	}

	fmt.Println("Successfully applied the rules to sessionId:", sessionId)
	fmt.Println("Status:", rulesListResp.GetStatus())
}

func drawANumber(ctx context.Context, client pb.GameClient, sessionId string) {
	drawNumberResp, err := client.DrawANumber(context.Background(),
	                               &pb.DrawNumberRequest{SessionId: sessionId,
		                       })
	if err != nil {
		log.Fatalf("failed to draw a number for session: %v", err)
	}

	fmt.Println("Successfully drawn a number for sessionId:", sessionId)

	fmt.Println("DrawnNumber:", drawNumberResp.GetNumber())
}

func drawnNumbersList(ctx context.Context, client pb.GameClient, sessionId string) {
	drawnNumbersResp, err := client.DrawnNumbersList(context.Background(),
	                               &pb.DrawnNumbersListRequest{SessionId: sessionId,
		                       })
	if err != nil {
		log.Fatalf("failed to pull drawn numbers list for session: %v", err)
	}

	fmt.Println("Successfully pulled drawn numbers list for sessionId:", sessionId)

	dnList := drawnNumbersResp.GetNumbers()

	fmt.Println("DrawnNumbersList(size):", len(dnList))
	fmt.Println("List:", dnList)
}

var wg sync.WaitGroup

func attachToDraws(ctx context.Context, client pb.GameClient, sessionId string) {
	stream, err := client.AttachToDraws(context.Background(),
	                                     &pb.AttachRequest{SessionId: sessionId,
		                         })
	if err != nil {
		log.Fatalf("failed to attach to stream of drawn numbers for session: %v", err)
	}

	fmt.Println("Successfully attached to stream of drawn numbers for sessionId:", sessionId)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Receiving stream of draws...")
		for {
			drawNumberResp, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Received EOF...")
				break
			}
			if err != nil {
				log.Fatalf("%v.AttachToDrawss(_) = _, %v", client, err)
			}
			fmt.Println("Received(Drawn Number):", drawNumberResp.GetNumber())
		}
	}()
}

/*
AttachToDraws(ctx context.Context, in *AttachRequest, opts ...grpc.CallOption) (Game_AttachToDrawsClient, error)
AnnounceWinners(ctx context.Context, in *AnnounceWinnersRequest, opts ...grpc.CallOption) (Game_AnnounceWinnersClient, error)
*/

func stopGame(ctx context.Context, client pb.GameClient, sessionId string) {
	stopGameResp, err := client.StopGame(context.Background(),
	                               &pb.StopSessionRequest{SessionId: sessionId,
		                       })
	if err != nil {
		log.Fatalf("failed to stop the game for session: %v", err)
	}

	fmt.Println("Successfully stopped the game for sessionId:", sessionId)

	fmt.Println("Status:", stopGameResp.GetStatus())
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

	startSessionResp := startNewGame(context.Background(), client)
	log.Println(startSessionResp.GetSessionId())

	gameLink := getGameLink(context.Background(), client, startSessionResp.GetSessionId())
	log.Println(gameLink)

	sheet := addPlayer(context.Background(), client, startSessionResp.GetSessionId())
	log.Println(sheet)

	playersList := listPlayers(context.Background(), client, startSessionResp.GetSessionId())
	log.Println(playersList)

	enablePlayer(context.Background(), client, startSessionResp.GetSessionId())

	applyRules(context.Background(), client, startSessionResp.GetSessionId())

	drawANumber(context.Background(), client, startSessionResp.GetSessionId())

	drawnNumbersList(context.Background(), client, startSessionResp.GetSessionId())

	attachToDraws(context.Background(), client, startSessionResp.GetSessionId())

	stopGame(context.Background(), client, startSessionResp.GetSessionId())

	wg.Wait()

	defer conn.Close()
}
