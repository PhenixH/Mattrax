package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/api"
	"github.com/mattrax/Mattrax/internal/authentication"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/middleware"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/mattrax/Mattrax/mdm"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var args mattrax.Arguments
	arg.MustParse(&args)
	// TODO: Verify arguments (eg. Domain is domain, cert paths exists, valid listen addr)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if args.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var tracer *zipkin.Tracer
	if args.Zipkin != "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal().Err(err).Msg("Error retrieving node hostname")
		}

		tracer, err = zipkin.NewTracer(reporterhttp.NewReporter(args.Zipkin), zipkin.WithSampler(zipkin.AlwaysSample),
			zipkin.WithLocalEndpoint(&model.Endpoint{ServiceName: hostname}))
		if err != nil {
			log.Fatal().Err(err).Msg("Error initialising Zipkin")
		}
	}

	dbconn, err := sql.Open("postgres", args.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initialising Postgres database connection")
	}
	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Error communicating with Postgres database")
	}

	q := db.New(dbconn)
	defer q.Close()

	// TODO: Check DB is working by querying

	var srv = &mattrax.Server{
		Args:         args,
		GlobalRouter: mux.NewRouter(),
		DB:           q,
		DBConn:       dbconn,
		Cache:        cache.New(5*time.Minute, 10*time.Minute),
	}
	if srv.Settings, err = settings.New(srv.DB); err != nil {
		log.Fatal().Err(err).Msg("Error starting settings service")
	}
	if srv.Cert, err = certificates.New(srv.DB); err != nil {
		log.Fatal().Err(err).Msg("Error starting certificates service")
	}
	if srv.Auth, err = authentication.New(srv.Cert, srv.Cache, srv.DB, args.Domain, args.Debug); err != nil {
		log.Fatal().Err(err).Msg("Error starting authentication service")
	}
	srv.GlobalRouter.Use(middleware.Logging())
	srv.GlobalRouter.Use(middleware.Headers())
	if tracer != nil {
		srv.GlobalRouter.Use(zipkinhttp.NewServerMiddleware(tracer, zipkinhttp.TagResponseSize(true)))
		srv.GlobalRouter.Use(middleware.ZipkinExtended)
	}
	srv.Router = srv.GlobalRouter.Schemes("https").Host(args.Domain).Subrouter()
	api.Mount(srv)
	mdm.Mount(srv)

	serve(args.Addr, args.Domain, args.TLSCert, args.TLSKey, nil, srv.GlobalRouter)
}
