package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/authn"
	api "github.com/cloustone/pandas/authn/api"
	grpcapi "github.com/cloustone/pandas/authn/api/grpc"
	httpapi "github.com/cloustone/pandas/authn/api/http"
	"github.com/cloustone/pandas/authn/jwt"
	"github.com/cloustone/pandas/authn/postgres"
	"github.com/cloustone/pandas/authn/tracing"
	mfidp "github.com/cloustone/pandas/authn/uuid"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/logger"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	opentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defLogLevel      = "error"
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "mainflux"
	defDBPass        = "mainflux"
	defDBName        = "authn"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defHTTPPort      = "8180"
	defGRPCPort      = "8181"
	defSecret        = "authn"
	defServerCert    = ""
	defServerKey     = ""
	defJaegerURL     = ""

	envLogLevel      = "PD_AUTHN_LOG_LEVEL"
	envDBHost        = "PD_AUTHN_DB_HOST"
	envDBPort        = "PD_AUTHN_DB_PORT"
	envDBUser        = "PD_AUTHN_DB_USER"
	envDBPass        = "PD_AUTHN_DB_PASS"
	envDBName        = "PD_AUTHN_DB"
	envDBSSLMode     = "PD_AUTHN_DB_SSL_MODE"
	envDBSSLCert     = "PD_AUTHN_DB_SSL_CERT"
	envDBSSLKey      = "PD_AUTHN_DB_SSL_KEY"
	envDBSSLRootCert = "PD_AUTHN_DB_SSL_ROOT_CERT"
	envHTTPPort      = "PD_AUTHN_HTTP_PORT"
	envGRPCPort      = "PD_AUTHN_GRPC_PORT"
	envSecret        = "PD_AUTHN_SECRET"
	envServerCert    = "PD_AUTHN_SERVER_CERT"
	envServerKey     = "PD_AUTHN_SERVER_KEY"
	envJaegerURL     = "PD_JAEGER_URL"
)

type config struct {
	logLevel   string
	dbConfig   postgres.Config
	httpPort   string
	grpcPort   string
	secret     string
	serverCert string
	serverKey  string
	jaegerURL  string
	resetURL   string
}

type tokenConfig struct {
	hmacSampleSecret []byte // secret for signing token
	tokenDuration    string // token in duration in min
}

func StartAuthnService() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	tracer, closer := initJaeger("authn", cfg.jaegerURL, logger)
	defer closer.Close()

	dbTracer, dbCloser := initJaeger("authn_db", cfg.jaegerURL, logger)
	defer dbCloser.Close()

	svc := newService(db, dbTracer, cfg.secret, logger)
	errs := make(chan error, 2)

	go startHTTPServer(tracer, svc, cfg.httpPort, cfg.serverCert, cfg.serverKey, logger, errs)
	go startGRPCServer(tracer, svc, cfg.grpcPort, cfg.serverCert, cfg.serverKey, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Authentication service terminated: %s", err))
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
		logLevel:   pandas.Env(envLogLevel, defLogLevel),
		dbConfig:   dbConfig,
		httpPort:   pandas.Env(envHTTPPort, defHTTPPort),
		grpcPort:   pandas.Env(envGRPCPort, defGRPCPort),
		secret:     pandas.Env(envSecret, defSecret),
		serverCert: pandas.Env(envServerCert, defServerCert),
		serverKey:  pandas.Env(envServerKey, defServerKey),
		jaegerURL:  pandas.Env(envJaegerURL, defJaegerURL),
	}

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
		logger.Error(fmt.Sprintf("Failed to init Jaeger: %s", err))
		os.Exit(1)
	}

	return tracer, closer
}

func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sqlx.DB {
	db, err := postgres.Connect(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func newService(db *sqlx.DB, tracer opentracing.Tracer, secret string, logger logger.Logger) authn.Service {
	database := postgres.NewDatabase(db)
	repo := tracing.New(postgres.New(database), tracer)

	idp := mfidp.New()
	t := jwt.New(secret)
	svc := authn.New(repo, idp, t)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "authn",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "authn",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func startHTTPServer(tracer opentracing.Tracer, svc authn.Service, port string, certFile string, keyFile string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if certFile != "" || keyFile != "" {
		logger.Info(fmt.Sprintf("Authentication service started using https, cert %s key %s, exposed port %s", certFile, keyFile, port))
		errs <- http.ListenAndServeTLS(p, certFile, keyFile, httpapi.MakeHandler(svc, tracer))
		return
	}
	logger.Info(fmt.Sprintf("Authentication service started using http, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, tracer))

}

func startGRPCServer(tracer opentracing.Tracer, svc authn.Service, port string, certFile string, keyFile string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	var server *grpc.Server
	if certFile != "" || keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to load authn certificates: %s", err))
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("Authentication gRPC service started using https on port %s with cert %s key %s", port, certFile, keyFile))
		server = grpc.NewServer(grpc.Creds(creds))
	} else {
		logger.Info(fmt.Sprintf("Authentication gRPC service started using http on port %s", port))
		server = grpc.NewServer()
	}

	mainflux.RegisterAuthNServiceServer(server, grpcapi.NewServer(tracer, svc))
	logger.Info(fmt.Sprintf("Authentication gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
