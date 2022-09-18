package tickets

import (
	"context"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type CreateTicketUseCaseRedisGateway interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type CreateTicketUseCaseMongoGateway interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
}

type CreateTicketUseCase struct {
	redisGateway CreateTicketUseCaseRedisGateway
	mongoGateway CreateTicketUseCaseMongoGateway
}

func NewCreateTicketUseCase(redisGateway CreateTicketUseCaseRedisGateway, mongoGateway CreateTicketUseCaseMongoGateway) *CreateTicketUseCase {
	return &CreateTicketUseCase{redisGateway: redisGateway, mongoGateway: mongoGateway}
}

type CreateTicketInput struct {
	PlayerID   string
	Parameters []entities.MatchmakingTicketParameter
}
type CreateTicketOutput struct {
	Ticket entities.MatchmakingTicket
}

func (c *CreateTicketUseCase) CreateTicket(ctx context.Context, input CreateTicketInput) (CreateTicketOutput, error) {
	ticket := entities.MatchmakingTicket{
		ID:         uuid.NewString(),
		Parameters: input.Parameters,
		PlayerID:   input.PlayerID,
	}

	var player entities.Player
	err := c.mongoGateway.FindOne(ctx, bson.M{"id": input.PlayerID}).Decode(&player)
	if err != nil {
		log.Print(err)
		return CreateTicketOutput{}, err
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
