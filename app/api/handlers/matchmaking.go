package handlers

import (
	"context"
	"encoding/json"
	"github.com/didopimentel/matchmaker/domain/matchmaking"
	"log"
	"net/http"
)

type MatchmakingAPIUseCases interface {
	MatchPlayers(ctx context.Context) (matchmaking.MatchPlayersOutput, error)
}

type MatchmakingAPI struct {
	uc MatchmakingAPIUseCases
}

func NewMatchmakingAPI(uc MatchmakingAPIUseCases) *MatchmakingAPI {
	return &MatchmakingAPI{uc: uc}
}

func (api *MatchmakingAPI) MatchPlayers(writer http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()

	output, err := api.uc.MatchPlayers(ctx)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshalledSessions, err := json.Marshal(output.CreatedSessions)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	AddHeaders(writer)
	writer.WriteHeader(http.StatusOK)
	writer.Write(marshalledSessions)
	return
}
