package handlers

import (
	"github.com/gorilla/mux"
)

type UseCases struct {
	TicketsAPIUseCases
	MatchmakingAPIUseCases
}

func NewServer(ucs UseCases) *mux.Router {
	router := mux.NewRouter()

	// Tickets API
	ticketsAPI := NewTicketsAPI(ucs.TicketsAPIUseCases)

	router.HandleFunc("/matchmaking/tickets", ticketsAPI.CreateMatchmakingTicket).Methods("POST")
	router.HandleFunc("/matchmaking/players/{id}/ticket", ticketsAPI.GetMatchmakingTicket).Methods("GET")

	// Matchmaking API
	matchmakingAPI := NewMatchmakingAPI(ucs.MatchmakingAPIUseCases)

	router.HandleFunc("/matchmaking/match-players", matchmakingAPI.MatchPlayers).Methods("GET")
	return router
}
