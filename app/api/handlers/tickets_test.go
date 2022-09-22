package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/didopimentel/matchmaker/app/api/handlers"
	"github.com/didopimentel/matchmaker/app/api/handlers/mocks"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTicketsAPI_GetMatchmakingTicket(t *testing.T) {
	t.Parallel()
	playerId := uuid.NewString()

	ticket := entities.MatchmakingTicket{
		ID:            uuid.NewString(),
		PlayerId:      uuid.NewString(),
		League:        5,
		Table:         5,
		CreatedAt:     time.Now().Unix(),
		Status:        entities.MatchmakingStatus_Pending,
		GameSessionId: "",
		Parameters: []entities.MatchmakingTicketParameter{
			{
				Type:     entities.MatchmakingTicketParameterType_Table,
				Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
				Value:    5,
			},
			{
				Type:     entities.MatchmakingTicketParameterType_League,
				Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
				Value:    5,
			},
		},
	}

	ticketsUc := &mocks.TicketsAPIUseCasesMock{
		GetTicketFunc: func(ctx context.Context, input tickets.GetTicketInput) (tickets.GetTicketOutput, error) {
			return tickets.GetTicketOutput{
				Ticket: ticket,
			}, nil
		},
	}

	server := handlers.NewServer(handlers.UseCases{
		TicketsAPIUseCases: ticketsUc,
	})

	testServer := httptest.NewServer(server)
	defer testServer.Close()

	endpoint := fmt.Sprintf("matchmaking/players/%s/ticket", playerId)
	res, err := http.Get(fmt.Sprintf("%s/%s", testServer.URL, endpoint))
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()
	bytesResponse, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var response handlers.GetMatchmakingTicketResponse
	err = json.Unmarshal(bytesResponse, &response)
	require.NoError(t, err)

	require.Equal(t, ticket.PlayerId, response.Ticket.PlayerId)
	require.Equal(t, ticket.Status, response.Ticket.Status)
	require.Equal(t, ticket.Parameters, response.Ticket.Parameters)
	require.Equal(t, ticket.Table, response.Ticket.Table)
	require.Equal(t, ticket.League, response.Ticket.League)
	require.Equal(t, ticket.CreatedAt, response.Ticket.CreatedAt)
	require.Equal(t, ticket.ID, response.Ticket.ID)
}

func TestTicketsAPI_GetMatchmakingTicket_Failure(t *testing.T) {
	t.Parallel()

	t.Run("should fail when no ticket is found for player", func(t *testing.T) {
		ticketsUc := &mocks.TicketsAPIUseCasesMock{
			GetTicketFunc: func(ctx context.Context, input tickets.GetTicketInput) (tickets.GetTicketOutput, error) {
				return tickets.GetTicketOutput{}, tickets.TicketNotFoundErr
			},
		}

		server := handlers.NewServer(handlers.UseCases{
			TicketsAPIUseCases: ticketsUc,
		})

		testServer := httptest.NewServer(server)
		defer testServer.Close()

		endpoint := fmt.Sprintf("matchmaking/players/%s/ticket", uuid.NewString())
		res, err := http.Get(fmt.Sprintf("%s/%s", testServer.URL, endpoint))
		require.NoError(t, err)

		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestTicketsAPI_CreateMatchmakingTicket(t *testing.T) {
	t.Parallel()

	ticket := entities.MatchmakingTicket{
		ID:            uuid.NewString(),
		PlayerId:      uuid.NewString(),
		League:        5,
		Table:         5,
		CreatedAt:     time.Now().Unix(),
		Status:        entities.MatchmakingStatus_Pending,
		GameSessionId: "",
		Parameters: []entities.MatchmakingTicketParameter{
			{
				Type:     entities.MatchmakingTicketParameterType_Table,
				Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
				Value:    5,
			},
			{
				Type:     entities.MatchmakingTicketParameterType_League,
				Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
				Value:    5,
			},
		},
	}

	ticketsUc := &mocks.TicketsAPIUseCasesMock{
		CreateTicketFunc: func(ctx context.Context, input tickets.CreateTicketInput) (tickets.CreateTicketOutput, error) {
			return tickets.CreateTicketOutput{
				Ticket: ticket,
			}, nil
		},
	}

	server := handlers.NewServer(handlers.UseCases{
		TicketsAPIUseCases: ticketsUc,
	})

	testServer := httptest.NewServer(server)
	defer testServer.Close()

	request := handlers.CreateMatchmakingTicketRequest{
		Parameters: []entities.MatchmakingTicketParameter{
			{
				Type:     entities.MatchmakingTicketParameterType_Table,
				Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
				Value:    5,
			},
			{
				Type:     entities.MatchmakingTicketParameterType_League,
				Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
				Value:    5,
			},
		},
		PlayerId: uuid.NewString(),
		League:   5,
		Table:    5,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(fmt.Sprintf("%s/matchmaking/tickets", testServer.URL), "application/json", &buf)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, res.StatusCode)

	defer res.Body.Close()
	bytesResponse, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var response entities.MatchmakingTicket
	err = json.Unmarshal(bytesResponse, &response)
	require.NoError(t, err)

	require.Equal(t, ticket.PlayerId, response.PlayerId)
	require.Equal(t, ticket.Status, response.Status)
	require.Equal(t, ticket.Parameters, response.Parameters)
	require.Equal(t, ticket.Table, response.Table)
	require.Equal(t, ticket.League, response.League)
	require.Equal(t, ticket.CreatedAt, response.CreatedAt)
	require.Equal(t, ticket.ID, response.ID)
}

func TestTicketsAPI_CreateMatchmakingTicket_Failure(t *testing.T) {
	t.Parallel()

	t.Run("should fail when parameters are invalid", func(t *testing.T) {
		ticketsUc := &mocks.TicketsAPIUseCasesMock{
			CreateTicketFunc: func(ctx context.Context, input tickets.CreateTicketInput) (tickets.CreateTicketOutput, error) {
				return tickets.CreateTicketOutput{}, tickets.InvalidTicketParametersErr
			},
		}

		server := handlers.NewServer(handlers.UseCases{
			TicketsAPIUseCases: ticketsUc,
		})

		testServer := httptest.NewServer(server)
		defer testServer.Close()

		request := handlers.CreateMatchmakingTicketRequest{
			Parameters: []entities.MatchmakingTicketParameter{
				{
					Type:     entities.MatchmakingTicketParameterType_Table,
					Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
					Value:    5,
				},
				{
					Type:     entities.MatchmakingTicketParameterType_League,
					Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
					Value:    5,
				},
			},
			PlayerId: uuid.NewString(),
			League:   5,
			Table:    5,
		}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.Post(fmt.Sprintf("%s/matchmaking/tickets", testServer.URL), "application/json", &buf)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
