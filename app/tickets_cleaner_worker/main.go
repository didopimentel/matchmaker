package main

import (
	"context"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/go-redis/redis/v9"
	"github.com/jasonlvhit/gocron"
	"log"
)

func main() {
	cfg, err := LoadConfig("./app/tickets_cleaner_worker")
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

	err = gocron.Every(cfg.WorkerTimeScheduleInSeconds).Seconds().Do(removeExpiredTicketsUseCase.RemoveExpiredTickets, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	<-gocron.Start()
}
