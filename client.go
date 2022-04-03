package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"time"
	"wow/client"
	"wow/config"
	"wow/logger"
)

func worker(c *client.Client, maxDelay int) {
	for {
		quote, err := c.GetQuote()

		if err != nil {
			log.Error(err)
		} else {
			fmt.Println(quote)
		}

		time.Sleep(time.Duration(rand.Intn(maxDelay)) * time.Millisecond)
	}
}

func main() {
	var configFile = flag.String("config", "./config.yaml", "Configuration file")

	flag.Parse()

	var conf config.ClientConfig

	err := config.Parse(*configFile, &conf)

	if err != nil {
		log.Fatalf("Read config file error: %s", err)
	}

	logger.Init(conf.LogLevel)

	serverAddr := conf.Client.ServerAddr

	if envServerAddr := os.Getenv("SERVER_ADDR"); envServerAddr != "" {
		serverAddr = envServerAddr
	}

	c := client.Client{
		ServerAddr:   serverAddr,
		ConnTimeout:  time.Duration(conf.Client.ConnTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.Client.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.Client.WriteTimeout) * time.Second,
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < conf.ParallelRequests; i++ {
		go worker(&c, conf.NextQuoteDelayMs)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
