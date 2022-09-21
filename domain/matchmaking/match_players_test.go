package matchmaking_test

import (
	"context"
	"github.com/didopimentel/matchmaker/domain/entities"
	"github.com/didopimentel/matchmaker/domain/matchmaking"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMatchPlayersUseCase_MatchPlayers(t *testing.T) {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})

	t.Run("should match two players with parameter EQUAL", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsRedisSetName := "tickets"
		matchesRedisSetName := "matches"
		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, ticketsRedisSetName)
		getTicketsUseCase := tickets.NewGetTicketUseCase(redisClient, ticketsRedisSetName, matchesRedisSetName)

		createTicketInputs := []tickets.CreateTicketInput{
			{
				PlayerId: uuid.NewString(),
				League:   5,
				Table:    6,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    5,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    6,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   7,
				Table:    8,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    7,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    8,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   7,
				Table:    8,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    7,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    8,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   10,
				Table:    11,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    10,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    11,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   15,
				Table:    16,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    15,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    16,
					},
				},
			},
		}

		for _, input := range createTicketInputs {
			_, err := ticketsUseCase.CreateTicket(ctx, input)
			require.NoError(t, err)
		}

		getTicketOutput, err := getTicketsUseCase.GetTicket(ctx, tickets.GetTicketInput{PlayerId: createTicketInputs[1].PlayerId})
		require.NoError(t, err)

		require.Equal(t, entities.MatchmakingStatus_Pending, getTicketOutput.Ticket.Status)

		matchPlayersUseCase := matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    2,
			MaxCountPerMatch:    2,
			TicketsRedisSetName: ticketsRedisSetName,
			MatchesRedisSetName: matchesRedisSetName,
		})
		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		require.Len(t, output.CreatedSessions, 1)

		// Should match 2nd and 3rd player
		for _, p := range output.CreatedSessions[0].PlayerIds {
			if p != createTicketInputs[1].PlayerId && p != createTicketInputs[2].PlayerId {
				t.Error("wrong players matched")
			}
		}

		getTicketOutput, err = getTicketsUseCase.GetTicket(ctx, tickets.GetTicketInput{PlayerId: output.CreatedSessions[0].PlayerIds[0]})
		require.NoError(t, err)

		gameSessionId := getTicketOutput.Ticket.GameSessionId
		require.Equal(t, entities.MatchmakingStatus_Found, getTicketOutput.Ticket.Status)
		require.NotEqual(t, "", getTicketOutput.Ticket.GameSessionId)

		getTicketOutput, err = getTicketsUseCase.GetTicket(ctx, tickets.GetTicketInput{PlayerId: output.CreatedSessions[0].PlayerIds[1]})
		require.NoError(t, err)

		require.Equal(t, entities.MatchmakingStatus_Found, getTicketOutput.Ticket.Status)
		require.NotEqual(t, "", getTicketOutput.Ticket.GameSessionId)

		require.Equal(t, gameSessionId, getTicketOutput.Ticket.GameSessionId)

		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

	t.Run("should match four players with parameter GREATER THAN", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInputs := []tickets.CreateTicketInput{
			{
				PlayerId: uuid.NewString(),
				League:   5,
				Table:    6,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_GreaterThan,
						Value:    5,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_GreaterThan,
						Value:    6,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   7,
				Table:    8,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    7,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    8,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   10,
				Table:    11,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    10,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    11,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   15,
				Table:    16,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    15,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    16,
					},
				},
			},
		}

		for _, input := range createTicketInputs {
			_, err := ticketsUseCase.CreateTicket(ctx, input)
			require.NoError(t, err)
		}

		matchPlayersUseCase := matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    2,
			MaxCountPerMatch:    4,
			TicketsRedisSetName: "tickets",
			MatchesRedisSetName: "matches",
			Timeout:             time.Minute,
		})
		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		require.Len(t, output.CreatedSessions, 1)

		require.Len(t, output.CreatedSessions[0].PlayerIds, 4)
		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

	t.Run("should match three players with parameter SMALLER THAN", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInputs := []tickets.CreateTicketInput{
			{
				PlayerId: uuid.NewString(),
				League:   5,
				Table:    6,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    5,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    6,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   7,
				Table:    8,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    7,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    8,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   10,
				Table:    11,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
						Value:    10,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_SmallerThan,
						Value:    11,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   15,
				Table:    16,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    15,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    16,
					},
				},
			},
		}

		for _, input := range createTicketInputs {
			_, err := ticketsUseCase.CreateTicket(ctx, input)
			require.NoError(t, err)
		}

		matchPlayersUseCase := matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    2,
			MaxCountPerMatch:    3,
			TicketsRedisSetName: "tickets",
			MatchesRedisSetName: "matches",
			Timeout:             time.Minute,
		})
		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		require.Len(t, output.CreatedSessions, 1)

		require.Len(t, output.CreatedSessions[0].PlayerIds, 3)
		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

	t.Run("should not match two players twice", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInputs := []tickets.CreateTicketInput{
			{
				PlayerId: uuid.NewString(),
				League:   5,
				Table:    6,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    5,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    6,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   7,
				Table:    8,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    7,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    8,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   7,
				Table:    8,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    7,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    8,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   10,
				Table:    11,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    10,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    11,
					},
				},
			},
			{
				PlayerId: uuid.NewString(),
				League:   15,
				Table:    16,
				Parameters: []entities.MatchmakingTicketParameter{
					{
						Type:     entities.MatchmakingTicketParameterType_League,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    15,
					},
					{
						Type:     entities.MatchmakingTicketParameterType_Table,
						Operator: entities.MatchmakingTicketParameterOperator_Equal,
						Value:    16,
					},
				},
			},
		}

		for _, input := range createTicketInputs {
			_, err := ticketsUseCase.CreateTicket(ctx, input)
			require.NoError(t, err)
		}

		matchPlayersUseCase := matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    2,
			MaxCountPerMatch:    2,
			TicketsRedisSetName: "tickets",
			MatchesRedisSetName: "matches",
			Timeout:             time.Minute,
		})
		_, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		require.Len(t, output.CreatedSessions, 0)

		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

	t.Run("should find an imperfect match in first try after expiring", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInput := tickets.CreateTicketInput{
			PlayerId: uuid.NewString(),
			League:   5,
			Table:    6,
			Parameters: []entities.MatchmakingTicketParameter{
				{
					Type:     entities.MatchmakingTicketParameterType_League,
					Operator: entities.MatchmakingTicketParameterOperator_Equal,
					Value:    5,
				},
				{
					Type:     entities.MatchmakingTicketParameterType_Table,
					Operator: entities.MatchmakingTicketParameterOperator_Equal,
					Value:    6,
				},
			},
		}

		_, err := ticketsUseCase.CreateTicket(ctx, createTicketInput)
		require.NoError(t, err)

		matchPlayersUseCase := matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    2,
			MaxCountPerMatch:    3,
			TicketsRedisSetName: "tickets",
			MatchesRedisSetName: "matches",
			Timeout:             time.Second * 3,
		})

		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)
		require.Len(t, output.CreatedSessions, 0)

		// wait 5 seconds for expiration
		time.Sleep(time.Second * 5)

		createAnotherTicketInput := tickets.CreateTicketInput{
			PlayerId: uuid.NewString(),
			League:   5,
			Table:    6,
			Parameters: []entities.MatchmakingTicketParameter{
				{
					Type:     entities.MatchmakingTicketParameterType_League,
					Operator: entities.MatchmakingTicketParameterOperator_Equal,
					Value:    5,
				},
				{
					Type:     entities.MatchmakingTicketParameterType_Table,
					Operator: entities.MatchmakingTicketParameterOperator_Equal,
					Value:    6,
				},
			},
		}

		// create another ticket that should match the requirements
		_, err = ticketsUseCase.CreateTicket(ctx, createAnotherTicketInput)
		require.NoError(t, err)

		output, err = matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		// TODO: check error
		require.Len(t, output.CreatedSessions, 1)

		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

	t.Run("should not match expired tickets", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInput := tickets.CreateTicketInput{
			PlayerId: uuid.NewString(),
			League:   5,
			Table:    6,
			Parameters: []entities.MatchmakingTicketParameter{
				{
					Type:     entities.MatchmakingTicketParameterType_League,
					Operator: entities.MatchmakingTicketParameterOperator_Equal,
					Value:    5,
				},
				{
					Type:     entities.MatchmakingTicketParameterType_Table,
					Operator: entities.MatchmakingTicketParameterOperator_Equal,
					Value:    6,
				},
			},
		}

		_, err := ticketsUseCase.CreateTicket(ctx, createTicketInput)
		require.NoError(t, err)

		matchPlayersUseCase := matchmaking.NewMatchPlayersUseCase(redisClient, matchmaking.MatchPlayerUseCaseConfig{
			MinCountPerMatch:    2,
			MaxCountPerMatch:    2,
			TicketsRedisSetName: "tickets",
			MatchesRedisSetName: "matches",
			Timeout:             time.Second * 3,
		})

		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)
		require.Len(t, output.CreatedSessions, 0)

		// wait 5 seconds for expiration
		time.Sleep(time.Second * 5)

		// run again since the first run after expiration tries to find an imperfect match
		output, err = matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)
		require.Len(t, output.CreatedSessions, 0)

		createAnotherTicketInput := tickets.CreateTicketInput{
			PlayerId: uuid.NewString(),
			League:   5,
			Table:    6,
			Parameters: []entities.MatchmakingTicketParameter{
				{
					Type:     entities.MatchmakingTicketParameterType_League,
					Operator: entities.MatchmakingTicketParameterOperator_Equal,
					Value:    5,
				},
				{
					Type:     entities.MatchmakingTicketParameterType_Table,
					Operator: entities.MatchmakingTicketParameterOperator_Equal,
					Value:    6,
				},
			},
		}

		// create another ticket that should match the requirements
		_, err = ticketsUseCase.CreateTicket(ctx, createAnotherTicketInput)
		require.NoError(t, err)

		output, err = matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		// TODO: check error
		require.Len(t, output.CreatedSessions, 0)

		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

}
