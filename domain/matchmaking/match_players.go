package matchmaking

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/go-redis/redis/v9"
)

type MatchPlayersUseCaseRedisGateway interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
}

type MatchPlayerUseCaseConfig struct {
	MinCountPerMatch    int32
	MaxCountPerMatch    int32
	TicketsRedisSetName string
}
type MatchPlayersUseCase struct {
	redisGateway MatchPlayersUseCaseRedisGateway
	cfg          MatchPlayerUseCaseConfig
}

func NewMatchPlayersUseCase(redisClient MatchPlayersUseCaseRedisGateway, config MatchPlayerUseCaseConfig) *MatchPlayersUseCase {
	return &MatchPlayersUseCase{redisGateway: redisClient, cfg: config}
}

type MatchPlayerInput struct {
	MinCount int32
	MaxCount int32
}

type MatchPlayersOutput struct {
	CreatedSessions []PlayerSession
}
type PlayerSession struct {
	PlayerIDs []string
}

func (m *MatchPlayersUseCase) MatchPlayers(ctx context.Context) (MatchPlayersOutput, error) {
	result := m.redisGateway.HScan(ctx, m.cfg.TicketsRedisSetName, 0, "", 10)

	var matchedSessions []PlayerSession
	alreadyMatchedPlayers := map[string]bool{}
	for {
		tickets, cursor, err := result.Result()
		if err != nil {
			return MatchPlayersOutput{}, err
		}

		for i := 0; i < len(tickets); i = i + 2 {
			if alreadyMatchedPlayers[tickets[i]] == true {
				continue
			}

			playerTicketBytes := []byte(tickets[i+1])

			var playerTicket entities.MatchmakingTicket
			err = json.Unmarshal(playerTicketBytes, &playerTicket)
			if err != nil {
				return MatchPlayersOutput{}, err
			}

			var eligibleOpponents []string
			// Append the player
			eligibleOpponents = append(eligibleOpponents, playerTicket.PlayerID)
			eligibleOpponentsCountMap := map[string]int{}
			for _, parameter := range playerTicket.Parameters {
				var result *redis.StringSliceCmd
				switch parameter.Operator {
				case entities.MatchmakingTicketParameterOperator_Equal:
					result = m.redisGateway.ZRangeByScore(ctx, string(parameter.Type), &redis.ZRangeBy{
						Min:   fmt.Sprint(parameter.Value),
						Max:   fmt.Sprint(parameter.Value),
						Count: int64(m.cfg.MaxCountPerMatch),
					})
				case entities.MatchmakingTicketParameterOperator_GreaterThan:
					result = m.redisGateway.ZRangeByScore(ctx, string(parameter.Type), &redis.ZRangeBy{
						Min:   fmt.Sprintf("(%d", parameter.Value),
						Max:   "+inf",
						Count: int64(m.cfg.MaxCountPerMatch),
					})
				case entities.MatchmakingTicketParameterOperator_SmallerThan:
					result = m.redisGateway.ZRangeByScore(ctx, string(parameter.Type), &redis.ZRangeBy{
						Min:   "0",
						Max:   fmt.Sprintf("(%d", parameter.Value),
						Count: int64(m.cfg.MaxCountPerMatch),
					})
				case entities.MatchmakingTicketParameterOperator_NotEqual:
					// TODO: support not equal operator
					continue
				default:
					// TODO: return error
					continue
				}

				// This will return the player ids of the eligible opponents
				foundOpponents, err := result.Result()
				if err != nil {
					return MatchPlayersOutput{}, err
				}

				for _, opponent := range foundOpponents {
					if opponent == playerTicket.PlayerID {
						continue
					}
					c, ok := eligibleOpponentsCountMap[opponent]
					if !ok {
						eligibleOpponentsCountMap[opponent] = 1
					} else {
						eligibleOpponentsCountMap[opponent] = c + 1
					}

					if eligibleOpponentsCountMap[opponent] == len(playerTicket.Parameters) {
						eligibleOpponents = append(eligibleOpponents, opponent)
					}

					if int32(len(eligibleOpponents)) == m.cfg.MaxCountPerMatch {
						break
					}
				}

			}

			// Found a match!
			if int32(len(eligibleOpponents)) >= m.cfg.MinCountPerMatch {
				matchedSessions = append(matchedSessions, PlayerSession{PlayerIDs: eligibleOpponents})
				for _, opponent := range eligibleOpponents {
					for _, parameter := range playerTicket.Parameters {
						if m.redisGateway.ZRem(ctx, string(parameter.Type), opponent).Err() != nil {
							return MatchPlayersOutput{}, err
						}
					}
					if m.redisGateway.HDel(ctx, m.cfg.TicketsRedisSetName, opponent).Err() != nil {
						return MatchPlayersOutput{}, err
					}
					alreadyMatchedPlayers[opponent] = true
				}
				// TODO: add match logic
			}

		}

		result = m.redisGateway.HScan(ctx, m.cfg.TicketsRedisSetName, 0, "", 10)
		// Finished iterating through matchmaking tickets
		if cursor == 0 {
			break
		}
	}

	return MatchPlayersOutput{
		matchedSessions,
	}, nil
}
