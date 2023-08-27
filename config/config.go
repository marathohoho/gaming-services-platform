package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type PlatformConfig struct {
	ApiServerHost   string `env:"API_SERVER" env-default:"localhost:5100"`
	UserServer      string `env:"USER_SERVER" env-default:"localhost:5101"`
	WalletServer    string `env:"WALLET_SERVER" env-default:"localhost:5102"`
	GrpcServer      string `env:"GRPC_SERVER" env-default:"localhost:5103"`
	WebsocketServer string `env:"WEBSOCKET_SERVER" env-default:"localhost:5104"`

	Redis struct {
		Addr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" env-default:""`
		Db       int    `env:"REDIS_PAYMENTCORE_DB" env-default:"2"`
	}

	UserServerHost   string `env:"USER_SERVER_HOST" env-default:"http://localhost:5101"`
	WalletServerHost string `env:"WALLET_SERVER_HOST" env-default:"http://localhost:5102"`
	GrpcServerHost   string `env:"GRPC_SERVER_HOST" env-default:"http://localhost:5103"`
}

func Init() PlatformConfig {
	var cfg PlatformConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read configuration")
	}

	return cfg
}
