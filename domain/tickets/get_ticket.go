package tickets

import (
	"context"
	"encoding/json"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/go-redis/redis/v9"
	"log"
)

type GetTicketUseCaseRedisGateway interface {
	HGet(ctx context.Context, key, field string) *redis.StringCmd
}

type GetTicketUseCase struct {
	redisGateway GetTicketUseCaseRedisGateway
}

func NewGetTicketUseCase(redisGateway GetTicketUseCaseRedisGateway) *GetTicketUseCase {
	return &GetTicketUseCase{redisGateway: redisGateway}
}

type GetTicketInput struct {
	PlayerID string
}
type GetTicketOutput struct {
	Ticket entities.MatchmakingTicket
}

func (c *GetTicketUseCase) GetTicket(ctx context.Context, input GetTicketInput) (GetTicketOutput, error) {
	result := c.redisGateway.HGet(ctx, "tickets", input.PlayerID)
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
		Ticket: ticket,
	}, nil
}
