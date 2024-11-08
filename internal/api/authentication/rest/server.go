package authentication

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	embed "github.com/blacksmith-vish/sso"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"

	"github.com/go-chi/chi/v5"
)

type Authentication interface {
	RegisterNewUser(
		ctx context.Context,
		request models.RegisterRequest,
	) (userID string, err error)
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

func (srv server) InitRouters(router *chi.Mux) {

	fs := http.FileServer(http.FS(embed.StaticFiles))
	// swagger:route GET /profile
	//
	// Gets profile of user
	//
	//     Produces:
	//     - application/json
	//     - application/x-protobuf
	//
	//     Schemes: http, https, ws, wss
	//
	//     Security:
	//       api_key:
	//       oauth: read, write
	//
	//     Responses:
	//       default: genericError
	//       200: someResponse
	//       422: validationError
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

	apiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// Mounting the new Sub Router on the main router
	router.Mount("/api", apiRouter)

}
