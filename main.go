package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/cenkalti/backoff/v5"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/godruoyi/go-snowflake"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/ip812/blog/config"
	"github.com/ip812/blog/logger"
	"github.com/ip812/blog/middleware"
	"github.com/ip812/blog/notifier"
	"github.com/ip812/blog/o11y"
	"github.com/ip812/blog/utils"
)

const (
	serviceName           = "blog"
	dbConnectTimeout      = 10 * time.Second
	dbMaxOpenConnections  = 10
	retryMaxElapsedTime   = 15 * time.Minute
	serverIdleTimeout     = 1 * time.Minute
	serverReadTimeout     = 10 * time.Second
	serverWriteTimeout    = 30 * time.Second
	serverShutdownTimeout = 10 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.New()
	log := logger.New(cfg)
	tracer, err := o11y.NewTracer(serviceName)
	if err != nil {
		log.Error("unable to initialize tracer due: %v", err)
	}

	// https://snowsta.mp
	startTime, _ := time.Parse(time.RFC3339, "2015-01-01T00:00:00Z")
	snowflake.SetStartTime(startTime)
	snowflake.SetMachineID(1)

	slacknotifier := notifier.NewSlack(cfg.Slack.BlogBotToken, log)

	swappableDB := NewSwappableDB()

	apiServer := startHTTPServer(cfg, log, tracer, swappableDB, slacknotifier)
	metricsServer := startMetricsServer(cfg, log)

	db, err := connectToDatabaseWithRetry(ctx, cfg, log)
	if err != nil {
		log.Error("exiting: could not connect to DB after retries: %s", err.Error())
		return
	}
	defer db.Close()

	swappableDB.Swap(db)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Error("failed to set dialect: %s", err.Error())
	}
	if err := goose.Up(db, "sql/migrations"); err != nil {
		log.Error("failed to run migrations: %s", err.Error())
	}

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Error("error shutting down server: %s", err.Error())
	} else {
		log.Info("server shutdown cleanly")
	}

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Error("error shutting down server: %s", err.Error())
	} else {
		log.Info("metrics server shutdown cleanly")
	}
}

type dbConnection struct {
	db *sql.DB
}

func connectToDatabaseWithRetry(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
	var conn dbConnection

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Endpoint,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	operation := func() (dbConnection, error) {
		connCtx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
		defer cancel()

		db, err := otelsql.Open(
			"postgres",
			connectionString,
			otelsql.WithAttributes(
				semconv.DBSystemNamePostgreSQL,
			),
		)
		if err != nil {
			log.Warn("failed to open the database connection: %v", err.Error())
			return conn, err
		}

		_, err = otelsql.RegisterDBStatsMetrics(
			db,
			otelsql.WithAttributes(
				semconv.DBSystemNamePostgreSQL,
			),
		)
		if err != nil {
			log.Warn("failed to register database metrics: %v", err.Error())
			return conn, err
		}

		if err := db.PingContext(connCtx); err != nil {
			log.Warn("failed to ping the database: %v", err.Error())
			return conn, err
		}

		db.SetMaxOpenConns(dbMaxOpenConnections)
		log.Info("connected to database")

		conn.db = db
		return conn, nil
	}

	_, err := backoff.Retry[dbConnection](
		ctx,
		operation,
		backoff.WithMaxElapsedTime(retryMaxElapsedTime),
	)

	return conn.db, err
}

func startHTTPServer(
	cfg *config.Config,
	log logger.Logger,
	tracer oteltrace.Tracer,
	db DBWrapper,
	slacknotifier *notifier.Slack,
) *http.Server {
	formDecoder := form.NewDecoder()
	formValidator := validator.New(validator.WithRequiredStructEnabled())

	handler := Handler{
		config:        cfg,
		formDecoder:   formDecoder,
		formValidator: formValidator,
		tracer:        tracer,
		slacknotifier: slacknotifier,
		db:            db,
		log:           log,
	}

	mux := chi.NewRouter()
	mux.Use(otelchi.Middleware(serviceName, otelchi.WithChiRoutes(mux)))
	mux.Use(middleware.TraceIDHeaderMiddleware)
	mux.Handle("/static/*", handler.StaticFiles())
	mux.With().Route("/p", func(mux chi.Router) {
		mux.Route("/public", func(mux chi.Router) {
			mux.Get("/landing-page", handler.LandingPageView)
			mux.Get("/articles", handler.ArticlesView)
			mux.Get("/articles/{id}", handler.ArticleDetailsView)
			mux.Get("/projects", handler.ProjectsView)
		})
	})

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/public/v0", func(mux chi.Router) {
			mux.Route("/articles", func(mux chi.Router) {
				mux.Post("/{id}/comments", utils.MakeTemplHandler(handler.CreateComment))
				mux.Get("/{id}/comments", utils.MakeTemplHandler(handler.GetAllCommentsByArticleID))
			})
		})
	})

	mux.Get("/healthz", handler.Healthz)
	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/p/public/landing-page", http.StatusFound)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		IdleTimeout:  serverIdleTimeout,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		Handler:      mux,
	}

	go func() {
		log.Info("server started on %s", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("cannot start server: %s", err.Error())
		}
	}()

	return server
}

func startMetricsServer(
	cfg *config.Config,
	log logger.Logger,
) *http.Server {
	mux := chi.NewRouter()

	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.MetricsPort),
		IdleTimeout:  serverIdleTimeout,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		Handler:      mux,
	}

	go func() {
		log.Info("metrics server started on %s", cfg.App.MetricsPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("cannot start metrics server: %s", err.Error())
		}
	}()

	return server
}
