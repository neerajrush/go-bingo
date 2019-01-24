package main

import (
	"fmt"
	"log"

	"github.com/neerajrush/go-bingo/proto"
)

type BingoServer interface {
	sessionId string
}

func (b *BingoServer) StartNewGame(ctx context.Context, in *StartSessionRequest) (*StartSessionResponse, error) {
	b.sessionId = in.GameName + in.SecretPhrase
	return &pb.StartSessionResponse{SessionId: b.sessionId}, nil
}

func (b *BingoServer) GetGameLink(ctx context.Context, in *GetSessionRequest) (*GetSessionResponse, error) {
}

func (b *BingoServer) AddAPlayer(ctx context.Context, in *AddPlayerRequest) (*AddPlayerResponse, error) {
}

func (b *BingoServer) ListPlayers(ctx context.Context, in *PlayersListRequest) (*PlayersListResponse, error) {
}

func (b *BingoServer) EnablePlayer(ctx context.Context, in *EnablePlayerRequest) (*EnablePlayerResponse, error) {
}

func (b *BingoServer) ApplyRules(ctx context.Context, in *RulesListRequest) (*RulesListResponse, error) {
}

func (b *BingoServer) DrawANumber(ctx context.Context, in *DrawNumberRequest) (*DrawNumberResponse, error) {
}

func (b *BingoServer) AttachToDraws(*AttachRequest, Game_AttachToDrawsServer) error {
}

func (b *BingoServer) DrawnNumbersList(ctx context.Context, in *DrawnNumbersListRequest) (*DrawnNumbersResponse, error) {
}

func (b *BingoServer) AnnounceWinners(*AnnounceWinnersRequest, Game_AnnounceWinnersServer) error {
}

func (b *BingoServer) StopGame(ctx context.Context, in *StopSessionRequest) (*StopSessionResponse, error) {
}

func main() {
}
