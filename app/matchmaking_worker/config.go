package main

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	RedisAddress                   string        `mapstructure:"REDIS_ADDRESS"`
	RedisPassword                  string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB                        int           `mapstructure:"REDIS_DB"`
	RedisTicketsSetName            string        `mapstructure:"REDIS_TICKETS_SET_NAME"`
	RedisMatchesSetName            string        `mapstructure:"REDIS_MATCHES_SET_NAME"`
	RedisCountPerIteration         int64         `mapstructure:"REDIS_COUNT_PER_ITERATION"`
	MatchmakerMinPlayersPerSession int32         `mapstructure:"MATCHMAKER_MIN_PLAYERS_PER_SESSION"`
	MatchmakerMaxPlayersPerSession int32         `mapstructure:"MATCHMAKER_MAX_PLAYERS_PER_SESSION"`
	MatchmakerTimeout              time.Duration `mapstructure:"MATCHMAKER_TIMEOUT"`
	WorkerTimeScheduleInSeconds    uint64        `mapstructure:"WORKER_TIME_SCHEDULE_IN_SECONDS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
