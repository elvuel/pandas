// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux/broker"
	"github.com/cloustone/pandas/mainflux/transformers/senml"
	"github.com/cloustone/pandas/mainflux/writers"
	"github.com/cloustone/pandas/mainflux/writers/api"
	"github.com/cloustone/pandas/mainflux/writers/postgres"
	"github.com/cloustone/pandas/pkg/logger"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	svcName = "postgres-writer"
	sep     = ","

	defNatsURL         = pandas.DefNatsURL
	defLogLevel        = "error"
	defPort            = "9104"
	defDBHost          = "postgres"
	defDBPort          = "5432"
	defDBUser          = "mainflux"
	defDBPass          = "mainflux"
	defDBName          = "messages"
	defDBSSLMode       = "disable"
	defDBSSLCert       = ""
	defDBSSLKey        = ""
	defDBSSLRootCert   = ""
	defSubjectsCfgPath = "/config/subjects.toml"

	envNatsURL         = "PD_NATS_URL"
	envLogLevel        = "PD_POSTGRES_WRITER_LOG_LEVEL"
	envPort            = "PD_POSTGRES_WRITER_PORT"
	envDBHost          = "PD_POSTGRES_WRITER_DB_HOST"
	envDBPort          = "PD_POSTGRES_WRITER_DB_PORT"
	envDBUser          = "PD_POSTGRES_WRITER_DB_USER"
	envDBPass          = "PD_POSTGRES_WRITER_DB_PASS"
	envDBName          = "PD_POSTGRES_WRITER_DB_NAME"
	envDBSSLMode       = "PD_POSTGRES_WRITER_DB_SSL_MODE"
	envDBSSLCert       = "PD_POSTGRES_WRITER_DB_SSL_CERT"
	envDBSSLKey        = "PD_POSTGRES_WRITER_DB_SSL_KEY"
	envDBSSLRootCert   = "PD_POSTGRES_WRITER_DB_SSL_ROOT_CERT"
	envSubjectsCfgPath = "PD_POSTGRES_WRITER_SUBJECTS_CONFIG"
)

type config struct {
	natsURL         string
	logLevel        string
	port            string
	dbConfig        postgres.Config
	subjectsCfgPath string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	b, err := broker.New(cfg.natsURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer b.Close()

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	repo := newService(db, logger)
	st := senml.New()
	if err = writers.Start(b, repo, st, svcName, cfg.subjectsCfgPath, logger); err != nil {
		logger.Error(fmt.Sprintf("Failed to create Postgres writer: %s", err))
	}

	errs := make(chan error, 2)

	go startHTTPServer(cfg.port, errs, logger)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Postgres writer service terminated: %s", err))
}

func loadConfig() config {
	dbConfig := postgres.Config{
		Host:        pandas.Env(envDBHost, defDBHost),
		Port:        pandas.Env(envDBPort, defDBPort),
		User:        pandas.Env(envDBUser, defDBUser),
		Pass:        pandas.Env(envDBPass, defDBPass),
		Name:        pandas.Env(envDBName, defDBName),
		SSLMode:     pandas.Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     pandas.Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      pandas.Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: pandas.Env(envDBSSLRootCert, defDBSSLRootCert),
	}

	return config{
		natsURL:         pandas.Env(envNatsURL, defNatsURL),
		logLevel:        pandas.Env(envLogLevel, defLogLevel),
		port:            pandas.Env(envPort, defPort),
		dbConfig:        dbConfig,
		subjectsCfgPath: pandas.Env(envSubjectsCfgPath, defSubjectsCfgPath),
	}
}

func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sqlx.DB {
	db, err := postgres.Connect(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func newService(db *sqlx.DB, logger logger.Logger) writers.MessageRepository {
	svc := postgres.New(db)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "postgres",
			Subsystem: "message_writer",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "postgres",
			Subsystem: "message_writer",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func startHTTPServer(port string, errs chan error, logger logger.Logger) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Postgres writer service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(svcName))
}
