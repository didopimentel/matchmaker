package matchmaking

import "context"

type MatchPlayersUseCaseRedisGateway interface {
}

type MatchPlayersUseCase struct {
	redisClient MatchPlayersUseCaseRedisGateway
}

func NewMatchPlayersUseCase(redisClient MatchPlayersUseCaseRedisGateway) *MatchPlayersUseCase {
	return &MatchPlayersUseCase{redisClient: redisClient}
}

type MatchPlayersOutput struct{}

func (m *MatchPlayersUseCase) MatchPlayers(ctx context.Context) (MatchPlayersOutput, error) {

	return MatchPlayersOutput{}, nil
}
