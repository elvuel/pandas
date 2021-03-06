// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/mainflux/broker"
	"github.com/cloustone/pandas/mainflux/coap"
	"github.com/cloustone/pandas/mainflux/coap/api"
	logger "github.com/cloustone/pandas/pkg/logger"
	thingsapi "github.com/cloustone/pandas/things/api/auth/grpc"
	gocoap "github.com/dustin/go-coap"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defPort          = "5683"
	defNatsURL       = pandas.DefNatsURL
	defThingsURL     = "localhost:8181"
	defLogLevel      = "error"
	defClientTLS     = "false"
	defCACerts       = ""
	defPingPeriod    = "12"
	defJaegerURL     = ""
	defThingsTimeout = "1" // in seconds

	envPort          = "PD_COAP_ADAPTER_PORT"
	envNatsURL       = "PD_NATS_URL"
	envThingsURL     = "PD_THINGS_URL"
	envLogLevel      = "PD_COAP_ADAPTER_LOG_LEVEL"
	envClientTLS     = "PD_COAP_ADAPTER_CLIENT_TLS"
	envCACerts       = "PD_COAP_ADAPTER_CA_CERTS"
	envPingPeriod    = "PD_COAP_ADAPTER_PING_PERIOD"
	envJaegerURL     = "PD_JAEGER_URL"
	envThingsTimeout = "PD_COAP_ADAPTER_THINGS_TIMEOUT"
)

type config struct {
	port          string
	natsURL       string
	thingsURL     string
	logLevel      string
	clientTLS     bool
	caCerts       string
	pingPeriod    time.Duration
	jaegerURL     string
	thingsTimeout time.Duration
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	conn := connectToThings(cfg, logger)
	defer conn.Close()

	thingsTracer, thingsCloser := initJaeger("things", cfg.jaegerURL, logger)
	defer thingsCloser.Close()

	cc := thingsapi.NewClient(conn, thingsTracer, cfg.thingsTimeout)
	respChan := make(chan string, 10000)

	b, err := broker.New(cfg.natsURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer b.Close()

	svc := coap.New(b, logger, cc, respChan)

	svc = api.LoggingMiddleware(svc, logger)

	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "coap_adapter",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "coap_adapter",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	errs := make(chan error, 2)

	go startHTTPServer(cfg.port, logger, errs)
	go startCOAPServer(cfg, svc, cc, respChan, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("CoAP adapter terminated: %s", err))
}

func loadConfig() config {
	tls, err := strconv.ParseBool(pandas.Env(envClientTLS, defClientTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envClientTLS)
	}

	pp, err := strconv.ParseInt(pandas.Env(envPingPeriod, defPingPeriod), 10, 64)
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envPingPeriod)
	}

	if pp < 1 || pp > 24 {
		log.Fatalf("Value of %s must be between 1 and 24", envPingPeriod)
	}

	timeout, err := strconv.ParseInt(pandas.Env(envThingsTimeout, defThingsTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envThingsTimeout, err.Error())
	}

	return config{
		thingsURL:     pandas.Env(envThingsURL, defThingsURL),
		natsURL:       pandas.Env(envNatsURL, defNatsURL),
		port:          pandas.Env(envPort, defPort),
		logLevel:      pandas.Env(envLogLevel, defLogLevel),
		clientTLS:     tls,
		caCerts:       pandas.Env(envCACerts, defCACerts),
		pingPeriod:    time.Duration(pp),
		jaegerURL:     pandas.Env(envJaegerURL, defJaegerURL),
		thingsTimeout: time.Duration(timeout) * time.Second,
	}
}

func connectToThings(cfg config, logger logger.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption
	if cfg.clientTLS {
		if cfg.caCerts != "" {
			tpc, err := credentials.NewClientTLSFromFile(cfg.caCerts, "")
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to load certs: %s", err))
				os.Exit(1)
			}
			opts = append(opts, grpc.WithTransportCredentials(tpc))
		}
	} else {
		logger.Info("gRPC communication is not encrypted")
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(cfg.thingsURL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to things service: %s", err))
		os.Exit(1)
	}
	return conn
}

func initJaeger(svcName, url string, logger logger.Logger) (opentracing.Tracer, io.Closer) {
	if url == "" {
		return opentracing.NoopTracer{}, ioutil.NopCloser(nil)
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler: &jconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: url,
			LogSpans:           true,
		},
	}.NewTracer()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to init Jaeger client: %s", err))
		os.Exit(1)
	}

	return tracer, closer
}

func startHTTPServer(port string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("CoAP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHTTPHandler())
}

func startCOAPServer(cfg config, svc coap.Service, auth mainflux.ThingsServiceClient, respChan chan<- string, l logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.port)
	l.Info(fmt.Sprintf("CoAP adapter service started, exposed port %s", cfg.port))
	errs <- gocoap.ListenAndServe("udp", p, api.MakeCOAPHandler(svc, auth, l, respChan, cfg.pingPeriod))
}
