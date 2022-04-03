package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"wow/cache"
	"wow/config"
	"wow/logger"
	"wow/metrics"
	"wow/server"
	"wow/storage"
)

func main() {
	var configFile = flag.String("config", "./config.yaml", "Configuration file")

	flag.Parse()

	var conf config.ServerConfig

	err := config.Parse(*configFile, &conf)

	if err != nil {
		log.Fatalf("Read config file error: %s", err)
	}

	logger.Init(conf.LogLevel)

	quotesStorage, err := storage.NewFileStorage(conf.QuotesFilename)

	if err != nil {
		log.Fatalf("Init storage error: %s", err)
	}

	s := server.NewServer(&server.Options{
		Address:       conf.Server.Listen,
		ReadTimeout:   conf.Server.ReadTimeout,
		WriteTimeout:  conf.Server.WriteTimeout,
		DDoSRate:      uint64(conf.Server.DDoSRate),
		TargetBits:    byte(conf.Server.TargetBits),
		ChallengeTtl:  uint32(conf.Server.ChallengeTtl),
		QuotesStorage: quotesStorage,
		Cache:         cache.NewSimpleCache(uint32(conf.CacheTtl)),
		RateMeter:     metrics.NewRateMeter(conf.RateInterval, conf.RateSize),
	})

	if err := s.Start(); err != nil {
		log.Error(err)
	}
}
