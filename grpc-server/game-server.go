package main

import (
	"flag"
	"fmt"
	"log"
	"errors"
	"context"
	"net"
	"time"
	"math/rand"
	"google.golang.org/grpc"

	pb "github.com/neerajrush/go-bingo/proto"
)

var (
	ErrSessionNotActive = errors.New("Requested session is not active")
	ErrDuplicatePlayer  = errors.New("Player name already added to the game")
)

type BingoServer struct {
	sessionId string
	players []string
	drawnNumbers []int32
	winners []string
}

func (b *BingoServer) StartNewGame(ctx context.Context, in *pb.StartSessionRequest) (*pb.StartSessionResponse, error) {
	b.sessionId = in.GetGameName() + "-" + in.GetSecretPhrase()
	log.Printf("SessionId: %v", b.sessionId)
	return &pb.StartSessionResponse{SessionId: b.sessionId}, nil
}

func (b *BingoServer) GetGameLink(ctx context.Context, in *pb.GetSessionRequest) (*pb.GetSessionResponse, error) {
	log.Printf("SessionId: %v", in.GetSessionId())
	if in.GetSessionId() == b.sessionId {
		return &pb.GetSessionResponse{SessionId: b.sessionId, GameLink: "http://localhost:8001/" + b.sessionId}, nil
    	}
	return nil, ErrSessionNotActive
}

func (b *BingoServer) AddAPlayer(ctx context.Context, in *pb.AddPlayerRequest) (*pb.AddPlayerResponse, error) {
	log.Printf("SessionId: %v", in.GetSessionId())
	log.Printf("PlayerName: %v", in.GetName())
	for _, player := range b.players {
		if player == in.GetName() {
			return nil, ErrDuplicatePlayer
		}
	}
	b.players = append(b.players, in.GetName())
	sheet := make([]*pb.AddPlayerResponse_Columns, 5)
	for i,_ := range sheet {
		sheet[i] = &pb.AddPlayerResponse_Columns{Cols: make([]int32, 5),}
	}
	if in.GetSessionId() == b.sessionId {
		return &pb.AddPlayerResponse{SheetSize: 5, Sheet: sheet}, nil
    	}
	return nil, ErrSessionNotActive
}

func (b *BingoServer) ListPlayers(ctx context.Context, in *pb.PlayersListRequest) (*pb.PlayersListResponse, error) {
	log.Printf("SessionId: %v", in.GetSessionId())
	if in.GetSessionId() == b.sessionId {
		plResp := pb.PlayersListResponse{Players: make([]string, len(b.players)), }
		copy(plResp.Players, b.players)
		return &plResp, nil
    	}
	return nil, ErrSessionNotActive
}

func (b *BingoServer) EnablePlayer(ctx context.Context, in *pb.EnablePlayerRequest) (*pb.EnablePlayerResponse, error) {
	log.Printf("SessionId: %v", in.GetSessionId())
	log.Printf("PlayerName: %v", in.GetPlayerName())
	if in.GetSessionId() == b.sessionId {
		return &pb.EnablePlayerResponse{PlayerName: in.GetPlayerName(), PlayerEnabled: true,}, nil
    	}
	return nil, ErrSessionNotActive
}

func (b *BingoServer) ApplyRules(ctx context.Context, in *pb.RulesListRequest) (*pb.RulesListResponse, error) {
	log.Printf("SessionId: %v", in.GetSessionId())
	log.Printf("RulesList: %v", in.GetRules())
	if in.GetSessionId() == b.sessionId {
		return &pb.RulesListResponse{Status: true,}, nil
    	}
	return nil, ErrSessionNotActive
}

func (b *BingoServer) DrawANumber(ctx context.Context, in *pb.DrawNumberRequest) (*pb.DrawNumberResponse, error) {
	log.Printf("SessionId: %v", in.GetSessionId())
	if in.GetSessionId() == b.sessionId {
		return &pb.DrawNumberResponse{Number: 23,}, nil
    	}
	return nil, ErrSessionNotActive
}

func (b *BingoServer) AttachToDraws(in *pb.AttachRequest, stream pb.Game_AttachToDrawsServer) error {
	log.Printf("SessionId: %v", in.GetSessionId())
	if in.GetSessionId() != b.sessionId {
		return nil
    	}
	log.Println("About to stream all drawnNumbers..", b.drawnNumbers)
	for _, draw := range b.drawnNumbers {
		if err := stream.Send(&pb.DrawNumberResponse{Number: draw, }); err != nil {
	        	return err
            	}
        }
	return nil
}

func (b *BingoServer) DrawnNumbersList(ctx context.Context, in *pb.DrawnNumbersListRequest) (*pb.DrawnNumbersResponse, error) {
	log.Printf("SessionId: %v", in.GetSessionId())
	if in.GetSessionId() == b.sessionId {
		rand.Seed(time.Now().Unix() + 99999)
		for i := 0; i < 10; i++ {
			aNo := rand.Intn(100)
			b.drawnNumbers = append(b.drawnNumbers, int32(aNo))
		}
		return &pb.DrawnNumbersResponse{Numbers: b.drawnNumbers,}, nil
    	}
	return nil, ErrSessionNotActive
}

func (b *BingoServer) AnnounceWinners(in *pb.AnnounceWinnersRequest, stream pb.Game_AnnounceWinnersServer) error {
	log.Printf("SessionId: %v", in.GetSessionId())
	if in.GetSessionId() != b.sessionId {
		return nil
    	}
	if len(b.winners) == 0 && len(b.players) > 0 {
		b.winners = append(b.winners, b.players[0])
	}
	if len(b.winners) == 1 && len(b.players) > 1 {
		b.winners = append(b.winners, b.players[1])
	}
	log.Println("Winners are:", b.winners)
	for _, winner := range b.winners {
		if err := stream.Send(&pb.AnnounceWinnersResponse{ SessionId: b.sessionId, Player: winner, }); err != nil {
	        	return err
                }
        }
	fmt.Println("winners(size):", len(b.winners))
        return nil
}

func (b *BingoServer) StopGame(ctx context.Context, in *pb.StopSessionRequest) (*pb.StopSessionResponse, error) {
	log.Printf("SessionId: %v", in.GetSessionId())
	if in.GetSessionId() == b.sessionId {
		return &pb.StopSessionResponse{SessionId: in.GetSessionId(), Status: true,}, nil
    	}
	return nil, ErrSessionNotActive
}

var (
	port = flag.Int("port", 8000, "Game server port")
)

func newGameServer() *BingoServer {
	s := &BingoServer{sessionId: "", players: make([]string, 0), drawnNumbers: make([]int32, 0), winners: make([]string, 0),}
	return s
}

func init() {

}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGameServer(grpcServer, newGameServer())
	grpcServer.Serve(lis)
}
