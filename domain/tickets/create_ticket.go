package tickets

import (
	"context"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"log"
	"time"
)

type CreateTicketUseCaseRedisGateway interface {
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd
}

type CreateTicketUseCase struct {
	redisGateway        CreateTicketUseCaseRedisGateway
	ticketsRedisSetName string
}

func NewCreateTicketUseCase(redisGateway CreateTicketUseCaseRedisGateway, ticketsRedisSetName string) *CreateTicketUseCase {
	return &CreateTicketUseCase{redisGateway: redisGateway, ticketsRedisSetName: ticketsRedisSetName}
}

type CreateTicketInput struct {
	PlayerId   string
	League     int64
	Table      int64
	Parameters []entities.MatchmakingTicketParameter
}
type CreateTicketOutput struct {
	Ticket entities.MatchmakingTicket
}

// CreateTicket creates a matchmaking ticket for a given player with its current League and Table state
// as well as the parameter requirements to match with other players
func (c *CreateTicketUseCase) CreateTicket(ctx context.Context, input CreateTicketInput) (CreateTicketOutput, error) {
	if len(input.Parameters) == 0 {
		return CreateTicketOutput{}, InvalidTicketParametersErr
	}

	ticket := entities.MatchmakingTicket{
		ID:         uuid.NewString(),
		PlayerId:   input.PlayerId,
		League:     input.League,
		Table:      input.Table,
		Parameters: input.Parameters,
		Status:     entities.MatchmakingStatus_Pending,
		CreatedAt:  time.Now().Unix(),
	}

	set := c.redisGateway.HSet(ctx, c.ticketsRedisSetName, input.PlayerId, ticket)
	if set.Err() != nil {
		log.Print(set.Err())
		return CreateTicketOutput{}, set.Err()
	}

	cmd := c.redisGateway.ZAdd(ctx, string(entities.MatchmakingTicketParameterType_Table), redis.Z{
		Score:  float64(input.Table),
		Member: input.PlayerId,
	})
	if cmd.Err() != nil {
		log.Print(cmd.Err())
	}
	cmd = c.redisGateway.ZAdd(ctx, string(entities.MatchmakingTicketParameterType_League), redis.Z{
		Score:  float64(input.League),
		Member: input.PlayerId,
	})
	if cmd.Err() != nil {
		log.Print(cmd.Err())
	}

	return CreateTicketOutput{
		Ticket: ticket,
	}, nil
}
