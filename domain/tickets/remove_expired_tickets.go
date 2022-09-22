package tickets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/go-redis/redis/v9"
	"time"
)

type RemoveExpiredTicketsUseCaseRedisGateway interface {
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd
}

type RemoveExpiredTicketsUseCaseConfig struct {
	TicketsRedisSetName string
	TimeBeforeToRemove  time.Duration
	CountPerIteration   int64
}
type RemoveExpiredTicketsUseCase struct {
	redisGateway RemoveExpiredTicketsUseCaseRedisGateway
	cfg          RemoveExpiredTicketsUseCaseConfig
}

func NewRemoveExpiredTicketsUseCase(redisGateway RemoveExpiredTicketsUseCaseRedisGateway, cfg RemoveExpiredTicketsUseCaseConfig) *RemoveExpiredTicketsUseCase {
	return &RemoveExpiredTicketsUseCase{redisGateway: redisGateway, cfg: cfg}
}

type RemoveExpiredTicketsOutput struct {
	ExpiredTicketsCount int64
}

// RemoveExpiredTickets removes expired tickets that were created before a certain time.
func (c *RemoveExpiredTicketsUseCase) RemoveExpiredTickets(ctx context.Context) (RemoveExpiredTicketsOutput, error) {
	var cursor uint64
	var tickets []string
	var err error
	var count int64
	for {
		result := c.redisGateway.HScan(ctx, c.cfg.TicketsRedisSetName, cursor, "", c.cfg.CountPerIteration)
		tickets, cursor, err = result.Result()
		if err != nil {
			return RemoveExpiredTicketsOutput{}, err
		}
		for i := 0; i < len(tickets); i = i + 2 {
			playerTicketBytes := []byte(tickets[i+1])
			var playerTicket entities.MatchmakingTicket
			err = json.Unmarshal(playerTicketBytes, &playerTicket)
			if err != nil {
				return RemoveExpiredTicketsOutput{}, err
			}

			// Removes if the ticket is expired and the time has passed the threshold
			if playerTicket.Status == entities.MatchmakingStatus_Expired && playerTicket.CreatedAt < time.Now().Add(-c.cfg.TimeBeforeToRemove).Unix() {
				if err = c.redisGateway.HDel(ctx, c.cfg.TicketsRedisSetName, playerTicket.PlayerId).Err(); err != nil {
					return RemoveExpiredTicketsOutput{}, err
				}
				count++
			}
		}

		// Finished iterating through matchmaking tickets
		if cursor == 0 {
			break
		}
	}

	fmt.Println("Tickets Cleaned: ", count)
	return RemoveExpiredTicketsOutput{
		ExpiredTicketsCount: count,
	}, nil
}
