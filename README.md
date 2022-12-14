# Matchmaker

Matchmaker is a service made in Golang that provides matchmaking functionality.

## Requirements

- Go 1.18
- GNU Make 3.81
- Docker 20.10.14 & Docker Compose 1.29.2

## Installation

Clone the repository in your go apps directory
```bash
git clone https://github.com/didopimentel/matchmaker.git
```

## How it works

There are three components in this system.
- **API**: provides methods to Create and Retrieve tickets for a player.
- **Matchmaker**: a worker that runs every X seconds and tries to match players by using parameters sent when creating the ticket
- **Cleaner**: a worker that runs every Y seconds and removes all tickets with status **expired**.

The first step is to create a ticket through the API by calling the endpoint `POST /matchmaking/tickets`. This endpoint 
expects the following fields in the body:
- PlayerId: the player's unique id
- PlayerParameters: an array of the player parameters
  - Type: the parameter type.
  - Value: the value of the parameters for the player. 
- MatchParameters: an array of parameters
  - Type: the parameter type.
  - Operator: the parameter operator. Must be one of `=`, `<` or `>`.
  - Value: the value to be compared to.

Example:

```json
{
  "PlayerId": "aec2c1a2-770c-4903-a18c-b3816b384b21",
  "PlayerParameters": [
    {
      "Type": "league",
      "Value": 6
    },
    {
      "Type": "table",
      "Value": 6
    }
  ],
  "MatchParameters": [
    {
      "Type": "league",
      "Operator": "=",
      "Value": 6
    },
    {
      "Type": "table",
      "Operator": "=",
      "Value": 6
    }
  ]
}
```

That operation will create the ticket with the status `pending`. The ticket can now be retrieved by calling `/matchmaking/players/{playerId}/ticket`.

After creating the ticket it will be in the pool for the matchmaker to start trying to create matches. The
 matchmaker uses the environment variables `MATCHMAKER_MIN_PLAYERS_PER_SESSION` and `MATCHMAKER_MAX_PLAYERS_PER_SESSION`
in order to find the correct amount of players for a session. It will always try to find the maximum amount of player possible.

The `MATCHMAKER_TIMEOUT` env variable is responsible for determining when tickets must expire. When a ticket expires,
the matchmaker first tries to get "more flexible" and reduce the maximum number of player (perfect match) by 1. If still it no match
is found, it sets the ticket as expired and no longer tries to use that specific ticket for matches, removing the player who is the
owner of the ticket from the matchmaking pool.

### It is a match!

When the matchmaker finds the correct amount of players who for a session, it creates a unique id that represents a GameSession. It will
then update the tickets from all players that were matched and set the GameSessionId and the status of the ticket to `found`.

Now, whenever the ticket is retrieved, you can check that the status is `found` and grab the GameSessionId!

## Usage

Run tests with command `make test`

Run API with command: `make run-api`

Run Matchmaker with command: `make run-matchmaker`

Run Tickets Cleaner with command: `make run-cleaner`

## Implementation

The matchmaker makes use of Redis Sorted Sets to match players and Hashes to store tickets. When a ticket is created, it adds, a scored member for 
each PlayerParameter.
When the matchmaking worker runs, it scans through all tickets grabbing the MatchParameters. For each of those parameters it tries to find a scored
member (player) that matches the requirement of that ticket. If that another player meets all requirements it adds to a temporary "match". A match
is only found if enough players meet all the requirements.

![Matchmaker Architecture](static/matchmaker.png "Matchmaker Architecture")