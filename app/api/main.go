package main

import (
	"github.com/didopimentel/matchmaker/app/api/handlers"
	"github.com/didopimentel/matchmaker/domain/matchmaking"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/go-redis/redis/v9"
	"log"
	"net/http"
)

func main() {
	cfg, err := LoadConfig("./app/api")
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	ticketsAPIUseCases := &struct {
		*tickets.CreateTicketUseCase
		*tickets.GetTicketUseCase
	}{
		CreateTicketUseCase: tickets.NewCreateTicketUseCase(redisClient, cfg.RedisTicketsSetName),
		GetTicketUseCase:    tickets.NewGetTicketUseCase(redisClient, cfg.RedisTicketsSetName, cfg.RedisMatchesSetName),
	}
	matchmakingAPIUseCases := &struct {
		*matchmaking.MatchPlayersUseCase
	}{
		MatchPlayersUseCase: matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    cfg.MatchmakerMinPlayersPerSession,
			MaxCountPerMatch:    cfg.MatchmakerMaxPlayersPerSession,
			TicketsRedisSetName: cfg.RedisTicketsSetName,
			MatchesRedisSetName: cfg.RedisMatchesSetName,
		}),
	}

	apiUseCases := handlers.UseCases{
		TicketsAPIUseCases:     ticketsAPIUseCases,
		MatchmakingAPIUseCases: matchmakingAPIUseCases,
	}

	server := handlers.NewServer(apiUseCases)

	log.Fatal(http.ListenAndServe(":8000", server))
}
