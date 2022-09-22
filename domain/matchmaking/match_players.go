package matchmaking

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"log"
	"time"
)

type MatchPlayersUseCaseRedisGateway interface {
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

type MatchPlayerUseCaseConfig struct {
	MinCountPerMatch    int32
	MaxCountPerMatch    int32
	TicketsRedisSetName string
	MatchesRedisSetName string
	Timeout             time.Duration
	CountPerIteration   int64
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
	SessionID string
	PlayerIds []string
}

// MatchPlayers tries to match all tickets opened by players.
// If a player's ticket exceeds the expiration time, reduces by one the amount of players
// needed for a perfect match. After that, if no match is found, sets the ticket as expired,
// so it can no longer match with other players.
func (m *MatchPlayersUseCase) MatchPlayers(ctx context.Context) (MatchPlayersOutput, error) {
	var cursor uint64
	var tickets []string
	var err error

	log.Println("Matching Players...")
	var matchedSessions []PlayerSession
	alreadyMatchedPlayers := map[string]bool{}
	for {
		result := m.redisGateway.HScan(ctx, m.cfg.TicketsRedisSetName, cursor, "", m.cfg.CountPerIteration)
		tickets, cursor, err = result.Result()
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

			// We don't try to match with anyone if the ticket has expired
			if playerTicket.Status == entities.MatchmakingStatus_Expired {
				continue
			}

			hasExpired := time.Now().Unix() > playerTicket.CreatedAt+int64(m.cfg.Timeout.Seconds())

			maxCountForThisPlayer := m.cfg.MaxCountPerMatch
			// when has reached the time limit, we decrease the max amount for a perfect by 1
			if hasExpired && maxCountForThisPlayer-1 >= m.cfg.MinCountPerMatch {
				maxCountForThisPlayer--
			}

			var eligibleOpponents []string
			// Append the player
			eligibleOpponents = append(eligibleOpponents, playerTicket.PlayerId)
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
					if opponent == playerTicket.PlayerId {
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
			if int32(len(eligibleOpponents)) == maxCountForThisPlayer {
				// this could be an id or the address of a game server match
				gameSessionId := uuid.New().String()
				matchedSessions = append(matchedSessions, PlayerSession{PlayerIds: eligibleOpponents, SessionID: gameSessionId})
				for _, opponent := range eligibleOpponents {
					for _, parameter := range playerTicket.Parameters {
						if err = m.redisGateway.ZRem(ctx, string(parameter.Type), opponent).Err(); err != nil {
							return MatchPlayersOutput{}, err
						}
					}
					if err = m.redisGateway.HDel(ctx, m.cfg.TicketsRedisSetName, opponent).Err(); err != nil {
						return MatchPlayersOutput{}, err
					}
					alreadyMatchedPlayers[opponent] = true

					// creates a registry in Matches for each opponent
					playerTicket.Status = entities.MatchmakingStatus_Found
					playerTicket.GameSessionId = gameSessionId
					m.redisGateway.HSet(ctx, m.cfg.MatchesRedisSetName, opponent, playerTicket)
				}
				// sets the ticket as expired and removes from parameters sets, so it is not tried again
			} else if hasExpired {
				playerTicket.Status = entities.MatchmakingStatus_Expired
				if err = m.redisGateway.HSet(ctx, m.cfg.TicketsRedisSetName, playerTicket.PlayerId, playerTicket).Err(); err != nil {
					return MatchPlayersOutput{}, err
				}

				for _, parameter := range playerTicket.Parameters {
					if err = m.redisGateway.ZRem(ctx, string(parameter.Type), playerTicket.PlayerId).Err(); err != nil {
						return MatchPlayersOutput{}, err
					}
				}
			}

		}

		// Finished iterating through matchmaking tickets
		if cursor == 0 {
			break
		}
	}

	log.Println("Matched Players: ", matchedSessions)
	return MatchPlayersOutput{
		matchedSessions,
	}, nil
}
