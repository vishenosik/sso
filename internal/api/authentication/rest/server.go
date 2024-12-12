package authentication

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	embed "github.com/blacksmith-vish/sso"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	"github.com/blacksmith-vish/sso/pkg/api"

	"github.com/go-chi/chi/v5"
)

type Authentication interface {
	Login(
		ctx context.Context,
		request models.LoginRequest,
		appID string,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		request models.RegisterRequest,
	) (userID string, err error)

	IsAdmin(
		ctx context.Context,
		userID string,
	) (isAdmin bool, err error)
}

type authenticationAPI struct {
	log  *slog.Logger
	auth Authentication
}

type server = *authenticationAPI

func NewAuthenticationServer(
	log *slog.Logger,
	auth Authentication,
) *authenticationAPI {

	return &authenticationAPI{
		log:  log,
		auth: auth,
	}

}

// ping godoc
//
//	@Summary 	Регистрация пользователя
//	@Tags 		system
//	@Router 	/api/ping [get]
//	@Produce 	html
//	@Success 	200 {string}  string    "ok"
//	@Failure 	406 {string}  string    "not ok"
func (srv server) InitRouters(router *chi.Mux) {

	fs := http.FileServer(http.FS(embed.StaticFiles))

	router.Handle("/static/*", fs)

	router.Post("/register", srv.register())

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		(w).Header().Set("Content-Type", "text;charset-utf-8")

		templ, err := template.ParseFS(embed.StaticFiles, "static/authentication/*.html")

		if err != nil {
			fmt.Println(err)
			return
		}

		templ.ExecuteTemplate(w, "index.html", nil)
	})

	// Creating a New Router
	apiRouter := chi.NewRouter()

	router.Get(api.ApiV1("/users/{user_id}/is_admin"), srv.isAdmin())

	apiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "ping failed", http.StatusNotAcceptable)
		// w.Write([]byte("PONG"))
	})

	// Mounting the new Sub Router on the main router
	router.Mount("/api", apiRouter)

}
