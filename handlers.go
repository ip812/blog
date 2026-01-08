package main

import (
	"embed"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/godruoyi/go-snowflake"
	"github.com/ip812/blog/articles"
	"github.com/ip812/blog/config"
	"github.com/ip812/blog/database"
	"github.com/ip812/blog/logger"
	"github.com/ip812/blog/notifier"
	"github.com/ip812/blog/status"
	"github.com/ip812/blog/templates/components"
	"github.com/ip812/blog/templates/views"
	"github.com/ip812/blog/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
        opsNewCommentsReceived = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "blog_new_comments_received_total",
                Help: "Total number of new comments received",
        })
)

//go:embed static
var staticFS embed.FS

type Handler struct {
	config        *config.Config
	formDecoder   *form.Decoder
	formValidator *validator.Validate
	tracer        oteltrace.Tracer
	slacknotifier *notifier.Slack
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

func (hnd *Handler) ProjectsView(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, views.Projects())
}

func (hnd *Handler) ArticleDetailsView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.Render(w, r, views.ArticleNotFound())
		return
	}

	if id == articles.ZeroTrustHomelabID {
		utils.Render(w, r, views.ArticleZeroTrustHomelab())
		return
	}

	if id == articles.ZeroTrustHomelabV2ID {
		utils.Render(w, r, views.ArticleZeroTrustHomelabV2())
		return
	}

	if id == articles.AnsiblePlusTailsclaleEqualGreatComboID {
		utils.Render(w, r, views.ArticleAnsiblePlusTailscaleEqualGreatCombo())
		return
	}

	if id == articles.DeferDeepDiveID {
		utils.Render(w, r, views.ArticleDeferDeepDive())
		return
	}

	utils.Render(w, r, views.ArticleNotFound())
}

func getOrSetUsername(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie(CookieKey)
	if err != nil {
		username := generateUsername()
		http.SetCookie(w, &http.Cookie{
			Name:     CookieKey,
			Value:    username,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})
		return username
	}
	return c.Value
}

func (hnd *Handler) CreateComment(w http.ResponseWriter, r *http.Request) error {
	articleID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		status.AddToast(w, status.WarningStatusBadRequest(status.WarnNotNumbericID))
		return utils.Render(w, r, components.NoComments())
	}

	err = r.ParseForm()
	if err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrParsingFrom))
		return utils.Render(w, r, components.NoComments())
	}
	var props components.CommentInputFormProps
	err = hnd.formDecoder.Decode(&props, r.Form)
	if err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrDecodingForm))
		return utils.Render(w, r, components.NoComments())
	}

	_, span := hnd.tracer.Start(
		r.Context(),
		"CreateComment(",
		oteltrace.WithAttributes(attribute.String("article", string(props.ArticleID))),
	)
	defer span.End()

	username := getOrSetUsername(w, r)

	db, err := hnd.db.DB()
	if err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrDB))
		return utils.Render(w, r, components.NoComments())
	}

	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrDB))
		return utils.Render(w, r, components.NoComments())
	}
	defer tx.Rollback()

	queries := database.New(tx)

	_, err = queries.CreateComment(r.Context(), database.CreateCommentParams{
		ID:        int64(snowflake.ID()),
		ArticleID: int64(articleID),
		Username:  username,
		Content:   props.Content,
	})
	if err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrCreateArticleComment))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return utils.Render(w, r, components.NoComments())
	}

	comments, err := queries.GetAllCommentsByArticleID(r.Context(), int64(articleID))
	if err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrGetAllArticleComments))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return utils.Render(w, r, components.NoComments())
	}

	if err := tx.Commit(); err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrDB))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return utils.Render(w, r, components.NoComments())
	}

	if len(comments) == 0 {
		hnd.log.Warn("no comments found after creating a comment")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return utils.Render(w, r, components.NoComments())
	}

	opsNewCommentsReceived.Inc()
	hnd.log.Info("comment created successfully for article ID %d", articleID)
	if err := hnd.slacknotifier.SendMsg(
		hnd.config.Slack.GeneralChannelID,
		fmt.Sprintf(
			"New comment for article *%s* was added:\n>%s\n>%s",
			articles.GetByID(articleID).Name,
			"https://blog.ip812.com/p/public/articles/"+strconv.FormatUint(articleID, 10),
			props.Content,
		),
	); err != nil {
		hnd.log.Error("failed to send Slack notification for new comment: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	commentProps := []components.CommentProps{}
	for _, c := range comments {
		commentProps = append(commentProps, components.CommentProps{
			ID:        uint64(c.ID),
			Username:  c.Username,
			AvatarURL: getAvatarURL(c.Username),
			Content:   c.Content,
		})
	}

	return utils.Render(w, r, components.Comments(commentProps))
}

func (hnd *Handler) GetAllCommentsByArticleID(w http.ResponseWriter, r *http.Request) error {
	db, err := hnd.db.DB()
	if err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrDB))
		return utils.Render(w, r, components.NoComments())
	}

	queries := database.New(db)

	articleID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		status.AddToast(w, status.WarningStatusBadRequest(status.WarnNotNumbericID))
		return utils.Render(w, r, components.NoComments())
	}

	comments, err := queries.GetAllCommentsByArticleID(r.Context(), int64(articleID))
	if err != nil {
		status.AddToast(w, status.ErrorInternalServerError(status.ErrGetAllArticleComments))
		return utils.Render(w, r, components.NoComments())
	}

	if len(comments) == 0 {
		return utils.Render(w, r, components.NoComments())
	}

	commentProps := []components.CommentProps{}
	for _, c := range comments {
		commentProps = append(commentProps, components.CommentProps{
			ID:        uint64(c.ID),
			Username:  c.Username,
			AvatarURL: getAvatarURL(c.Username),
			Content:   c.Content,
		})
	}

	return utils.Render(w, r, components.Comments(commentProps))
}
