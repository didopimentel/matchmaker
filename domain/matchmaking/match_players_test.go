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

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInputs := []tickets.CreateTicketInput{
			{
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
		})
		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		require.Len(t, output.CreatedSessions, 1)

		// Should match 2nd and 3rd player
		for _, p := range output.CreatedSessions[0].PlayerIDs {
			if p != createTicketInputs[1].PlayerID && p != createTicketInputs[2].PlayerID {
				t.Error("wrong players matched")
			}
		}

		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

	t.Run("should match four players with parameter GREATER THAN", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInputs := []tickets.CreateTicketInput{
			{
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
		})
		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		require.Len(t, output.CreatedSessions, 1)

		require.Len(t, output.CreatedSessions[0].PlayerIDs, 4)
		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

	t.Run("should match three players with parameter SMALLER THAN", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInputs := []tickets.CreateTicketInput{
			{
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
		})
		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		require.Len(t, output.CreatedSessions, 1)

		require.Len(t, output.CreatedSessions[0].PlayerIDs, 3)
		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

	t.Run("should not match two players twice", func(t *testing.T) {
		cmd := redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())

		ticketsUseCase := tickets.NewCreateTicketUseCase(redisClient, "tickets")

		createTicketInputs := []tickets.CreateTicketInput{
			{
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
				PlayerID: uuid.NewString(),
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
		})
		_, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		output, err := matchPlayersUseCase.MatchPlayers(ctx)
		require.NoError(t, err)

		require.Len(t, output.CreatedSessions, 0)

		cmd = redisClient.FlushAll(ctx)
		require.NoError(t, cmd.Err())
	})

}
