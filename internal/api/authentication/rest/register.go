package authentication

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/blacksmith-vish/sso/internal/services/authentication/models"

	"github.com/go-playground/validator/v10"

	"github.com/pkg/errors"
)

// Register godoc
// @Summary      Register a new user
// @Description  Registers a new user with the provided nickname, email, and password
// @Tags         authentication
// @Accept       x-www-form-urlencoded
// @Produce      plain
// @Param        nickname  formData  string  true  "User's nickname"
// @Param        email     formData  string  true  "User's email"
// @Param        pswd      formData  string  true  "User's password"
// @Success      200  {string}  string  "Hello BITCH!"
// @Failure      400  {string} string ""
// @Failure      406  {string} string ""
// @Failure      417  {string} string ""
// @Failure      500  {string} string ""
// @Router       /register [post]
func (srv server) register() http.HandlerFunc {

	log := srv.log.With(
		slog.String("op", "authentication.register"),
	)

	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			log.Error(err.Error())
		}

		serviceRequest := models.RegisterRequest{
			Nickname: r.Form.Get("nickname"),
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("pswd"),
		}

		if err := validator.New().Struct(serviceRequest); err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
			return
		}

		_, err = srv.auth.RegisterNewUser(
			context.Background(),
			serviceRequest,
		)

		if err != nil {
			if errors.Is(err, models.ErrUserExists) {

				log.Error("registration failed", slog.String("err", err.Error()))

				http.Error(w, errors.Wrap(err, "registration failed").Error(), http.StatusNotAcceptable)
				return
			}

			log.Error("registration failed", slog.String("err", err.Error()))

			http.Error(w, errors.Wrap(err, "registration failed").Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Hello BITCH!"))
	}
}
