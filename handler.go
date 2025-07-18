package main

import (
	"database/sql"
	"embed"
	"net/http"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/ip812/blog/config"
	"github.com/ip812/blog/database"
	"github.com/ip812/blog/logger"
)

//go:embed static
var staticFS embed.FS

type Handler struct {
	config        *config.Config
	formDecoder   *form.Decoder
	formValidator *validator.Validate
	db            *sql.DB
	queries       *database.Queries
	log           logger.Logger
}

func (hnd *Handler) StaticFiles() http.Handler {
	if hnd.config.App.Env == config.Local {
		hnd.log.Info("serving static files from local directory")
		return http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	}

	hnd.log.Info("serving static files from embedded FS")
	return http.StripPrefix("/", http.FileServer(http.FS(staticFS)))
}

func (hnd *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

func (hnd *Handler) HomeRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/p/public/home", http.StatusFound)
}
