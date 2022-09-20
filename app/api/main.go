package main

import (
	"github.com/didopimentel/matchmaker/app/api/handlers"
	"github.com/didopimentel/matchmaker/domain/matchmaking"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	cfg, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})
	router := mux.NewRouter()

	// Tickets API
	ticketsAPI := handlers.NewTicketsAPI(&struct {
		*tickets.CreateTicketUseCase
		*tickets.GetTicketUseCase
	}{
		CreateTicketUseCase: tickets.NewCreateTicketUseCase(redisClient, cfg.RedisTicketsSetName),
		GetTicketUseCase:    tickets.NewGetTicketUseCase(redisClient, cfg.RedisTicketsSetName, cfg.RedisMatchesSetName),
	})

	router.HandleFunc("/matchmaking/tickets", ticketsAPI.CreateMatchmakingTicket).Methods("POST")
	router.HandleFunc("/matchmaking/players/{id}/ticket", ticketsAPI.GetMatchmakingTicket).Methods("GET")

	// Matchmaking API
	matchmakingAPI := handlers.NewMatchmakingAPI(&struct {
		*matchmaking.MatchPlayersUseCase
	}{
		MatchPlayersUseCase: matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    cfg.MatchmakerMinPlayersPerSession,
			MaxCountPerMatch:    cfg.MatchmakerMaxPlayersPerSession,
			TicketsRedisSetName: cfg.RedisTicketsSetName,
			MatchesRedisSetName: cfg.RedisMatchesSetName,
		}),
	})

	router.HandleFunc("/matchmaking/match-players", matchmakingAPI.MatchPlayers).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
