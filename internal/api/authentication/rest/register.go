package authentication

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/blacksmith-vish/sso/internal/services/authentication/models"

	"github.com/go-playground/validator/v10"

	"github.com/pkg/errors"
)

type Account struct {
	ID string
}

type HTTPError struct {
	err string
}

// getArticleById godoc
//
//	@Summary 	Регистрация пользователя
//	@Tags 		users
//	@Router 	/register/us/article_id [post]
//	@Accept 	json
//	@Param 		article_id path string true "Идентификатор баннера" example(5bs8879f-1e45-4043-a8f6-a8f7b934f45a)
//	@Param 		request body models.RegisterRequest true "Данные пользователя"
//	@Success 	200
//	@Failure 	406
//	@Failure 	417
//	@Failure 	500

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  Account
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /accounts/{id} [get]
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
