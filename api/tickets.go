package api

import (
	"context"
	"encoding/json"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type TicketsAPIUseCases interface {
	CreateTicket(ctx context.Context, input tickets.CreateTicketInput) (tickets.CreateTicketOutput, error)
	GetTicket(ctx context.Context, input tickets.GetTicketInput) (tickets.GetTicketOutput, error)
}

type TicketsAPI struct {
	uc TicketsAPIUseCases
}

func NewTicketsAPI(uc TicketsAPIUseCases) *TicketsAPI {
	return &TicketsAPI{uc: uc}
}

type CreateMatchmakingTicketRequest struct {
	Parameters []entities.MatchmakingTicketParameter `json:"parameters"`
	PlayerID   string                                `json:"player_id"`
}

func (api *TicketsAPI) CreateMatchmakingTicket(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req CreateMatchmakingTicketRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := api.uc.CreateTicket(ctx, tickets.CreateTicketInput{Parameters: req.Parameters, PlayerID: req.PlayerID})
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshalledTicket, err := json.Marshal(output.Ticket)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write(marshalledTicket)
	return
}

func (api *TicketsAPI) GetMatchmakingTicket(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(request)

	playerId, ok := vars["id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := api.uc.GetTicket(ctx, tickets.GetTicketInput{PlayerID: playerId})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	ticketBytes, err := json.Marshal(output.Ticket)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(ticketBytes)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
