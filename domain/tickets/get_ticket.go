package tickets

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/go-redis/redis/v9"
	"log"
)

type GetTicketUseCaseRedisGateway interface {
	HGet(ctx context.Context, key, field string) *redis.StringCmd
}

type GetTicketUseCase struct {
	redisGateway        GetTicketUseCaseRedisGateway
	ticketsRedisSetName string
	matchesRedisSetName string
}

func NewGetTicketUseCase(redisGateway GetTicketUseCaseRedisGateway, ticketsRedisSetName, matchesRedisSetName string) *GetTicketUseCase {
	return &GetTicketUseCase{redisGateway: redisGateway, ticketsRedisSetName: ticketsRedisSetName, matchesRedisSetName: matchesRedisSetName}
}

type GetTicketInput struct {
	PlayerId string
}
type GetTicketOutput struct {
	Ticket entities.MatchmakingTicket
}

func (c *GetTicketUseCase) GetTicket(ctx context.Context, input GetTicketInput) (GetTicketOutput, error) {
	result := c.redisGateway.HGet(ctx, c.matchesRedisSetName, input.PlayerId)
	if result.Err() != nil {
		if !errors.Is(result.Err(), redis.Nil) {
			log.Print(result.Err())
			return GetTicketOutput{}, result.Err()
		}
	} else {
		// The match was already found for the player
		var ticketBytes []byte
		err := result.Scan(&ticketBytes)
		if err != nil {
			return GetTicketOutput{}, err
		}

		var ticket entities.MatchmakingTicket
		err = json.Unmarshal(ticketBytes, &ticket)
		if err != nil {
			return GetTicketOutput{}, err
		}
		return GetTicketOutput{
			Ticket: ticket,
		}, nil
	}

	result = c.redisGateway.HGet(ctx, c.ticketsRedisSetName, input.PlayerId)
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			return GetTicketOutput{}, TicketNotFoundErr
		}

		log.Print(result.Err())
		return GetTicketOutput{}, result.Err()
	}

	var ticketBytes []byte
	err := result.Scan(&ticketBytes)
	if err != nil {
		return GetTicketOutput{}, err
	}

	var ticket entities.MatchmakingTicket
	err = json.Unmarshal(ticketBytes, &ticket)
	if err != nil {
		return GetTicketOutput{}, err
	}

	return GetTicketOutput{
		Ticket: ticket,
	}, nil
}
