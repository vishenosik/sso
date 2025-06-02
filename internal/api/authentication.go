package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	// pkg
	"github.com/go-chi/chi/v5"

	// internal pkg

	_http "github.com/vishenosik/gocherry/pkg/http"

	//internal

	"github.com/vishenosik/sso/internal/entities"
)

type Authentication interface {
	//
	LoginByEmail(
		ctx context.Context,
		email,
		password,
		appID string,
	) (string, error)
	//
	RegisterUser(
		ctx context.Context,
		user *entities.UserCreds,
	) (string, error)
	//
	IsAdmin(
		ctx context.Context,
		userID string,
	) (bool, error)
}

type AuthenticationHttpApi struct {
	log  *slog.Logger
	auth Authentication
}

func NewAuthenticationHttpApi(
	auth Authentication,
) *AuthenticationHttpApi {

	return &AuthenticationHttpApi{
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
func (auth *AuthenticationHttpApi) Routers(r chi.Router) {

	r.Group(func(r chi.Router) {
		r.Route(auth.registerUser())
		r.Route(routeUsers("get"), func(r chi.Router) {
			r.Get(_http.BlankRoute, func(w http.ResponseWriter, r *http.Request) {

				response := struct {
					Message string `json:"message"`
					Status  string `json:"status"`
				}{
					Message: "getter",
					Status:  "endpoints.get ok",
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "failed to encode response", http.StatusInternalServerError)
					return
				}
			})
		})
	})
}

func (auth *AuthenticationHttpApi) registerUser() (string, func(chi.Router)) {

	versionMiddleware, versionHandler := _http.DotVersionMiddlewareHandler("1.0")

	return routeUsers("register"), func(r chi.Router) {
		r.Use(
			versionMiddleware,
		)
		r.Post(_http.BlankRoute, versionHandler(_http.HandlersMap{
			"1.0": auth.registerUser_1_0(),
		}))
	}
}

func (auth *AuthenticationHttpApi) registerUser_1_0() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

var routeUsers = _http.MethodFunc("users")
