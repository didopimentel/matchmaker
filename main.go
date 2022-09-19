package main

import (
	"github.com/didopimentel/matchmaker/api"
	"github.com/didopimentel/matchmaker/domain/matchmaking"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})
	router := mux.NewRouter()

	ticketsRedisSetName := "tickets"

	// Tickets API
	ticketsAPI := api.NewTicketsAPI(&struct {
		*tickets.CreateTicketUseCase
		*tickets.GetTicketUseCase
	}{
		CreateTicketUseCase: tickets.NewCreateTicketUseCase(redisClient, ticketsRedisSetName),
		GetTicketUseCase:    tickets.NewGetTicketUseCase(redisClient, ticketsRedisSetName),
	})

	router.HandleFunc("/matchmaking/tickets", ticketsAPI.CreateMatchmakingTicket).Methods("POST")
	router.HandleFunc("/matchmaking/players/{id}/ticket", ticketsAPI.GetMatchmakingTicket).Methods("GET")

	// Matchmaking API
	matchmakingAPI := api.NewMatchmakingAPI(&struct {
		*matchmaking.MatchPlayersUseCase
	}{
		MatchPlayersUseCase: matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    2,
			MaxCountPerMatch:    4,
			TicketsRedisSetName: ticketsRedisSetName,
		}),
	})

	router.HandleFunc("/matchmaking/match-players", matchmakingAPI.MatchPlayers).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
