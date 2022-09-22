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
- playerId: the player's unique id
- league: the player's current league
- table: the player's current table
- parameters: an array of parameters
  - type: the parameter type. Must be one of `league` or `table`.
  - operator: the parameter operator. Must be one of `=`, `<` or `>`.
  - value: the value to be compared to.

Example:

```json
{
  "playerId": "aec2c1a2-770c-4903-a18c-b3816b384b21",          
  "league": 6,               
  "table": 6,
  "parameters": [
    {
      "type": "league",
      "operator": "=",
      "value": 6
    },
    {
      "type": "table",
      "operator": "=",
      "value": 6
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


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)