package main

import (
	"context"
	"github.com/didopimentel/matchmaker/api"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})
	mongoOpts := options.Client()
	mongoOpts.SetAuth(options.Credential{
		Username: "test",
		Password: "test",
	})
	mongoOpts.ApplyURI("mongodb://localhost:27017")
	mongoClient, err := mongo.NewClient(mongoOpts)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()

	ticketsAPI := api.NewTicketsAPI(&struct {
		*tickets.CreateTicketUseCase
		*tickets.GetTicketUseCase
	}{
		CreateTicketUseCase: tickets.NewCreateTicketUseCase(redisClient, mongoClient.Database("player_states").Collection("player_states")),
		GetTicketUseCase:    tickets.NewGetTicketUseCase(redisClient),
	})

	router.HandleFunc("/matchmaking/tickets", ticketsAPI.CreateMatchmakingTicket).Methods("POST")
	router.HandleFunc("/matchmaking/tickets/{id}", ticketsAPI.GetMatchmakingTicket).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
