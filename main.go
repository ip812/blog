package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/godruoyi/go-snowflake"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/ip812/blog/config"
	"github.com/ip812/blog/database"
	"github.com/ip812/blog/logger"
)

func main() {
	cfg := config.New()
	log := logger.New(cfg)

	// https://snowsta.mp
	snowflake.SetStartTime(time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC))
	snowflake.SetMachineID(1)

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Endpoint,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Error("failed to connect to database: %s", err.Error())
	}
	defer db.Close()
	db.SetMaxOpenConns(10)
	queries := database.New(db)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Error("failed to set dialect: %s", err.Error())
	}
	if err := goose.Up(db, "sql/migrations"); err != nil {
		log.Error("failed to run migrations: %s", err.Error())
	}

	formDecoder := form.NewDecoder()
	formValidator := validator.New(validator.WithRequiredStructEnabled())

	hnd := Handler{
		config:        cfg,
		formDecoder:   formDecoder,
		formValidator: formValidator,
		db:            db,
		queries:       queries,
		log:           log,
	}

	mux := chi.NewRouter()
	mux.Handle("/static/*", hnd.StaticFiles())
	mux.With().Route("/p", func(mux chi.Router) {
		mux.Route("/public", func(mux chi.Router) {
		})
		mux.Route("/client", func(mux chi.Router) {
			// No handlers yet
		})
		mux.Route("/admin", func(mux chi.Router) {
			// No handlers yet
		})
	})
	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/public/v0", func(mux chi.Router) {
			mux.Route("/articles", func(mux chi.Router) {
			})
		})
		mux.Route("/client/v0", func(mux chi.Router) {
			// No handlers yet
		})
		mux.Route("/admin/v0", func(mux chi.Router) {
			// No handlers yet
		})
	})
	mux.Get("/healthz", hnd.Healthz)
	mux.NotFound(hnd.HomeRedirect)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      mux,
	}
	log.Info("server started on %s", cfg.App.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Error("cannot start server: %s", err.Error())
	}
}
