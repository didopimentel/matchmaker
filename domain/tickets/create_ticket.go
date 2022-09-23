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

type CreateTicketInputPlayerParameters struct {
	Type  entities.MatchmakingTicketParameterType
	Value float64
}
type CreateTicketInput struct {
	PlayerId         string
	PlayerParameters []CreateTicketInputPlayerParameters
	MatchParameters  []entities.MatchmakingTicketParameter
}
type CreateTicketOutput struct {
	Ticket entities.MatchmakingTicket
}

// CreateTicket creates a matchmaking ticket for a given player with its current League and Table state
// as well as the parameter requirements to match with other players
func (c *CreateTicketUseCase) CreateTicket(ctx context.Context, input CreateTicketInput) (CreateTicketOutput, error) {
	ticket := entities.MatchmakingTicket{
		ID:              uuid.NewString(),
		PlayerId:        input.PlayerId,
		MatchParameters: input.MatchParameters,
		Status:          entities.MatchmakingStatus_Pending,
		CreatedAt:       time.Now().Unix(),
	}

	set := c.redisGateway.HSet(ctx, c.ticketsRedisSetName, input.PlayerId, ticket)
	if set.Err() != nil {
		log.Print(set.Err())
		return CreateTicketOutput{}, set.Err()
	}

	playerParameterMap := map[entities.MatchmakingTicketParameterType]float64{}
	for _, parameter := range input.PlayerParameters {
		playerParameterMap[parameter.Type] = parameter.Value
	}

	for _, parameter := range input.MatchParameters {
		score, ok := playerParameterMap[parameter.Type]
		if !ok {
			// if the player does not pass the parameter he won't be matched with other players who request it
			continue
		}

		cmd := c.redisGateway.ZAdd(ctx, string(parameter.Type), redis.Z{
			Score:  score,
			Member: input.PlayerId,
		})
		if cmd.Err() != nil {
			log.Print(cmd.Err())
		}
	}

	return CreateTicketOutput{
		Ticket: ticket,
	}, nil
}
