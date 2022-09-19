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
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd
}

type CreateTicketUseCase struct {
	redisGateway CreateTicketUseCaseRedisGateway
}

func NewCreateTicketUseCase(redisGateway CreateTicketUseCaseRedisGateway) *CreateTicketUseCase {
	return &CreateTicketUseCase{redisGateway: redisGateway}
}

type CreateTicketInput struct {
	PlayerID   string
	League     int64
	Table      int64
	Parameters []entities.MatchmakingTicketParameter
}
type CreateTicketOutput struct {
	Ticket entities.MatchmakingTicket
}

func (c *CreateTicketUseCase) CreateTicket(ctx context.Context, input CreateTicketInput) (CreateTicketOutput, error) {
	ticket := entities.MatchmakingTicket{
		ID:         uuid.NewString(),
		PlayerID:   input.PlayerID,
		League:     input.League,
		Table:      input.Table,
		Parameters: input.Parameters,
	}

	// TODO: parameterize ttl
	set := c.redisGateway.HSet(ctx, "tickets", input.PlayerID, ticket)
	if set.Err() != nil {
		log.Print(set.Err())
		return CreateTicketOutput{}, set.Err()
	}

	cmd := c.redisGateway.ZAdd(ctx, string(entities.MatchmakingTicketParameterType_Table), redis.Z{
		Score:  float64(input.Table),
		Member: input.PlayerID,
	})
	if cmd.Err() != nil {
		log.Print(cmd.Err())
	}
	cmd = c.redisGateway.ZAdd(ctx, string(entities.MatchmakingTicketParameterType_League), redis.Z{
		Score:  float64(input.League),
		Member: input.PlayerID,
	})
	if cmd.Err() != nil {
		log.Print(cmd.Err())
	}

	return CreateTicketOutput{
		Ticket: ticket,
	}, nil
}
