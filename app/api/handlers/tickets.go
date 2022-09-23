package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

//go:generate moq -stub -pkg mocks -out mocks/tickets_api_use_cases.go . TicketsAPIUseCases

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
	MatchParameters  []entities.MatchmakingTicketParameter       `json:"MatchParameters"`
	PlayerId         string                                      `json:"PlayerId"`
	PlayerParameters []tickets.CreateTicketInputPlayerParameters `json:"PlayerParameters"`
}

// ValidateCreateMatchmakingTicket TODO: we might want to limit the amount of parameters since it might affect performance
func (api *TicketsAPI) ValidateCreateMatchmakingTicket(ctx context.Context, req CreateMatchmakingTicketRequest) error {
	if len(req.MatchParameters) == 0 {
		return tickets.InvalidTicketParametersErr
	}

	if len(req.PlayerParameters) == 0 {
		return tickets.InvalidPlayerParameters
	}

	return nil
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

	err = api.ValidateCreateMatchmakingTicket(ctx, req)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	output, err := api.uc.CreateTicket(ctx, tickets.CreateTicketInput{
		PlayerId:         req.PlayerId,
		PlayerParameters: req.PlayerParameters,
		MatchParameters:  req.MatchParameters,
	})
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

	AddHeaders(writer)
	writer.WriteHeader(http.StatusCreated)
	writer.Write(marshalledTicket)
	return
}

type GetMatchmakingTicketResponse struct {
	Ticket entities.MatchmakingTicket
}

func (api *TicketsAPI) GetMatchmakingTicket(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(request)

	playerId, ok := vars["id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := api.uc.GetTicket(ctx, tickets.GetTicketInput{PlayerId: playerId})
	if err != nil {
		if errors.Is(err, tickets.TicketNotFoundErr) {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(err.Error()))
			return
		}

		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := GetMatchmakingTicketResponse{
		Ticket: output.Ticket,
	}
	ticketBytes, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	AddHeaders(writer)

	_, err = writer.Write(ticketBytes)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
