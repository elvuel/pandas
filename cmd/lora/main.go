// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux/broker"
	"github.com/cloustone/pandas/mainflux/lora"
	"github.com/cloustone/pandas/mainflux/lora/api"
	"github.com/cloustone/pandas/mainflux/lora/mqtt"
	"github.com/cloustone/pandas/pkg/logger"
	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	r "github.com/go-redis/redis"

	"github.com/cloustone/pandas/mainflux/lora/redis"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	defHTTPPort       = "8180"
	defLoraMsgURL     = "tcp://localhost:1883"
	defNatsURL        = pandas.DefNatsURL
	defLogLevel       = "error"
	defESURL          = "localhost:6379"
	defESPass         = ""
	defESDB           = "0"
	defESConsumerName = "lora"
	defRouteMapURL    = "localhost:6379"
	defRouteMapPass   = ""
	defRouteMapDB     = "0"

	envHTTPPort       = "PD_LORA_ADAPTER_HTTP_PORT"
	envLoraMsgURL     = "PD_LORA_ADAPTER_MESSAGES_URL"
	envNatsURL        = "PD_NATS_URL"
	envLogLevel       = "PD_LORA_ADAPTER_LOG_LEVEL"
	envESURL          = "PD_THINGS_ES_URL"
	envESPass         = "PD_THINGS_ES_PASS"
	envESDB           = "PD_THINGS_ES_DB"
	envESConsumerName = "PD_LORA_ADAPTER_EVENT_CONSUMER"
	envRouteMapURL    = "PD_LORA_ADAPTER_ROUTE_MAP_URL"
	envRouteMapPass   = "PD_LORA_ADAPTER_ROUTE_MAP_PASS"
	envRouteMapDB     = "PD_LORA_ADAPTER_ROUTE_MAP_DB"

	loraServerTopic = "application/+/device/+/rx"

	thingsRMPrefix   = "thing"
	channelsRMPrefix = "channel"
)

type config struct {
	httpPort       string
	loraMsgURL     string
	natsURL        string
	logLevel       string
	esURL          string
	esPass         string
	esDB           string
	esConsumerName string
	routeMapURL    string
	routeMapPass   string
	routeMapDB     string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	rmConn := connectToRedis(cfg.routeMapURL, cfg.routeMapPass, cfg.routeMapDB, logger)
	defer rmConn.Close()

	esConn := connectToRedis(cfg.esURL, cfg.esPass, cfg.esDB, logger)
	defer esConn.Close()

	b, err := broker.New(cfg.natsURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer b.Close()

	thingRM := newRouteMapRepositoy(rmConn, thingsRMPrefix, logger)
	chanRM := newRouteMapRepositoy(rmConn, channelsRMPrefix, logger)

	mqttConn := connectToMQTTBroker(cfg.loraMsgURL, logger)

	svc := lora.New(b, thingRM, chanRM)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "lora_adapter",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "lora_adapter",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	go subscribeToLoRaBroker(svc, mqttConn, logger)
	go subscribeToThingsES(svc, esConn, cfg.esConsumerName, logger)

	errs := make(chan error, 2)

	go startHTTPServer(cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("LoRa adapter terminated: %s", err))
}

func loadConfig() config {
	return config{
		httpPort:       pandas.Env(envHTTPPort, defHTTPPort),
		loraMsgURL:     pandas.Env(envLoraMsgURL, defLoraMsgURL),
		natsURL:        pandas.Env(envNatsURL, defNatsURL),
		logLevel:       pandas.Env(envLogLevel, defLogLevel),
		esURL:          pandas.Env(envESURL, defESURL),
		esPass:         pandas.Env(envESPass, defESPass),
		esDB:           pandas.Env(envESDB, defESDB),
		esConsumerName: pandas.Env(envESConsumerName, defESConsumerName),
		routeMapURL:    pandas.Env(envRouteMapURL, defRouteMapURL),
		routeMapPass:   pandas.Env(envRouteMapPass, defRouteMapPass),
		routeMapDB:     pandas.Env(envRouteMapDB, defRouteMapDB),
	}
}

func connectToMQTTBroker(loraURL string, logger logger.Logger) mqttPaho.Client {
	opts := mqttPaho.NewClientOptions()
	opts.AddBroker(loraURL)
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetOnConnectHandler(func(c mqttPaho.Client) {
		logger.Info("Connected to Lora MQTT broker")
	})
	opts.SetConnectionLostHandler(func(c mqttPaho.Client, err error) {
		logger.Error(fmt.Sprintf("MQTT connection lost: %s", err.Error()))
		os.Exit(1)
	})

	client := mqttPaho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Lora MQTT broker: %s", token.Error()))
		os.Exit(1)
	}

	return client
}

func connectToRedis(redisURL, redisPass, redisDB string, logger logger.Logger) *r.Client {
	db, err := strconv.Atoi(redisDB)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to redis: %s", err))
		os.Exit(1)
	}

	return r.NewClient(&r.Options{
		Addr:     redisURL,
		Password: redisPass,
		DB:       db,
	})
}

func subscribeToLoRaBroker(svc lora.Service, mc mqttPaho.Client, logger logger.Logger) {
	mqtt := mqtt.NewBroker(svc, mc, logger)
	logger.Info("Subscribed to Lora MQTT broker")
	if err := mqtt.Subscribe(loraServerTopic); err != nil {
		logger.Error(fmt.Sprintf("Failed to subscribe to Lora MQTT broker: %s", err))
		os.Exit(1)
	}
}

func subscribeToThingsES(svc lora.Service, client *r.Client, consumer string, logger logger.Logger) {
	eventStore := redis.NewEventStore(svc, client, consumer, logger)
	logger.Info("Subscribed to Redis Event Store")
	if err := eventStore.Subscribe("mainflux.things"); err != nil {
		logger.Warn(fmt.Sprintf("Lora-adapter service failed to subscribe to Redis event source: %s", err))
	}
}

func newRouteMapRepositoy(client *r.Client, prefix string, logger logger.Logger) lora.RouteMapRepository {
	logger.Info(fmt.Sprintf("Connected to %s Redis Route-map", prefix))
	return redis.NewRouteMapRepository(client, prefix)
}

func startHTTPServer(cfg config, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.httpPort)
	logger.Info(fmt.Sprintf("lora-adapter service started, exposed port %s", cfg.httpPort))
	errs <- http.ListenAndServe(p, api.MakeHandler())
}
