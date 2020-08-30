package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/server/http"
	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/storages/psqlstore"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = logger.ApplyConfig(config.Logger.Path, config.Logger.Level)
	if err != nil {
		log.Fatalf("can't setup logger: %s", err)
	}

	if err := Migrate(config.DB.URL); err != nil {
		log.Fatalf("can't perform migrations: %s", err)
	}

	storage, _ := psqlstore.NewPSQLStore(config.DB.URL)
	calendar := app.New(storage)

	server := internalhttp.NewServer(
		calendar,
		config.Server.Host+":"+config.Server.Port,
		config.Server.ShutDownTimeOut)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals
		signal.Stop(signals)

		if err := server.Stop(); err != nil {
			log.Fatalf("%s", err)
		}
	}()

	if err := server.Start(); err != nil {
		log.Fatalf("%s", err)
	}
}
