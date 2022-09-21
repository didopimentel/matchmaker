package main

import (
	"context"
	"github.com/didopimentel/matchmaker/domain/tickets"
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

	removeExpiredTicketsUseCase := tickets.NewRemoveExpiredTicketsUseCase(redisClient, tickets.RemoveExpiredTicketsUseCaseConfig{
		TicketsRedisSetName: cfg.RedisTicketsSetName,
		TimeBeforeToRemove:  cfg.TicketsTimeBeforeToRemove,
		CountPerIteration:   cfg.RedisCountPerIteration,
	})

	err = gocron.Every(20).Seconds().Do(removeExpiredTicketsUseCase.RemoveExpiredTickets, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	<-gocron.Start()
}
