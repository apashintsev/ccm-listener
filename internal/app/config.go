package app

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"os"
)

type appConfig struct {
	Logger struct {
		LogLvl string // debug, info, error
	}

	NetworkConfig *liteclient.GlobalConfig

	TargetContractAddress *address.Address

	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SslMode  string
		Timezone string
	}
}

var CFG *appConfig = &appConfig{}

func InitConfig() error {
	godotenv.Load(".env")

	CFG.Logger.LogLvl = os.Getenv("LOG_LVL")

	jsonConfig, err := os.Open("network-config.json") //лучше взять по адресу https://ton.org/testnet-global.config.json динамически
	if err != nil {
		return err
	}

	if target, err := address.ParseAddr(os.Getenv("ADDRESS_TO_LISTEN")); err != nil {
		return err
	} else {
		CFG.TargetContractAddress = target
	}

	if err := json.NewDecoder(jsonConfig).Decode(&CFG.NetworkConfig); err != nil {
		return err
	}
	defer jsonConfig.Close()

	CFG.Postgres.Host = os.Getenv("POSTGRES_HOST")
	CFG.Postgres.Port = os.Getenv("POSTGRES_PORT")
	CFG.Postgres.User = os.Getenv("POSTGRES_USER")
	CFG.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	CFG.Postgres.Name = os.Getenv("POSTGRES_DB")
	CFG.Postgres.SslMode = os.Getenv("POSTGRES_SSLMODE")
	CFG.Postgres.Timezone = os.Getenv("POSTGRES_TIMEZONE")

	return nil
}
