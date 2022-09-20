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
	PlayerID string
}
type GetTicketOutput struct {
	Status        entities.MatchmakingStatus
	GameSessionId string
	Ticket        entities.MatchmakingTicket
}

func (c *GetTicketUseCase) GetTicket(ctx context.Context, input GetTicketInput) (GetTicketOutput, error) {
	result := c.redisGateway.HGet(ctx, c.matchesRedisSetName, input.PlayerID)
	var gameSessionId string
	if result.Err() != nil {
		if !errors.Is(result.Err(), redis.Nil) {
			log.Print(result.Err())
			return GetTicketOutput{}, result.Err()
		}
	} else {
		var err error
		gameSessionId, err = result.Result()
		if err != nil {
			return GetTicketOutput{}, err
		}
	}

	if gameSessionId != "" {
		return GetTicketOutput{
			Status:        entities.MatchmakingStatus_Found,
			GameSessionId: gameSessionId,
			Ticket:        entities.MatchmakingTicket{},
		}, nil
	}

	// TODO: what to do when expired?
	result = c.redisGateway.HGet(ctx, c.ticketsRedisSetName, input.PlayerID)
	if result.Err() != nil {
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
		Status: entities.MatchmakingStatus_Pending,
		Ticket: ticket,
	}, nil
}
