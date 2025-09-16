package main

import (
	"embed"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/ip812/blog/articles"
	"github.com/ip812/blog/config"
	"github.com/ip812/blog/logger"
	"github.com/ip812/blog/templates/views"
	"github.com/ip812/blog/utils"
)

//go:embed static
var staticFS embed.FS

type Handler struct {
	config        *config.Config
	formDecoder   *form.Decoder
	formValidator *validator.Validate
	log           logger.Logger

	db DBWrapper
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

func (hnd *Handler) LandingPageView(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, views.LandingPage())
}

func (hnd *Handler) ArticlesView(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, views.Articles())
}

func (hnd *Handler) ArticleDetailsView(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Render(w, r, views.ArticleNotFound())
		return
	}

	if id == articles.ZeroTrustHomelabID {
		utils.Render(w, r, views.ArticleZeroTrustHomelab())
		return
	}

	utils.Render(w, r, views.ArticleNotFound())
}

func (hnd *Handler) ProjectsView(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, views.Projects())
}

func (hnd *Handler) CreateComment(w http.ResponseWriter, r *http.Request) error {
	return nil
}
