package main

import (
	"context"
	"github.com/didopimentel/matchmaker/domain/matchmaking"
	"github.com/go-redis/redis/v9"
	"github.com/jasonlvhit/gocron"
	"log"
)

func main() {
	cfg, err := LoadConfig("./app/matchmaking_worker")
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	matchmakingUseCase := matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
		MinCountPerMatch:    cfg.MatchmakerMinPlayersPerSession,
		MaxCountPerMatch:    cfg.MatchmakerMaxPlayersPerSession,
		TicketsRedisSetName: cfg.RedisTicketsSetName,
		MatchesRedisSetName: cfg.RedisMatchesSetName,
		Timeout:             cfg.MatchmakerTimeout,
	})

	err = gocron.Every(20).Seconds().Do(matchmakingUseCase.MatchPlayers, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	<-gocron.Start()
}
