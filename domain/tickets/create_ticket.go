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
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type CreateTicketUseCase struct {
	redisGateway CreateTicketUseCaseRedisGateway
}

func NewCreateTicketUseCase(redisGateway CreateTicketUseCaseRedisGateway) *CreateTicketUseCase {
	return &CreateTicketUseCase{redisGateway: redisGateway}
}

type CreateTicketInput struct {
	Parameters []entities.MatchmakingTicketParameter
}
type CreateTicketOutput struct {
	Ticket entities.MatchmakingTicket
}

func (c *CreateTicketUseCase) CreateTicket(ctx context.Context, input CreateTicketInput) (CreateTicketOutput, error) {
	ticket := entities.MatchmakingTicket{
		ID:         uuid.NewString(),
		Parameters: input.Parameters,
	}

	// TODO: parameterize ttl
	set := c.redisGateway.Set(ctx, ticket.ID, ticket, 5*time.Minute)
	if set.Err() != nil {
		log.Print(set.Err())
		return CreateTicketOutput{}, set.Err()
	}

	return CreateTicketOutput{
		Ticket: ticket,
	}, nil
}
