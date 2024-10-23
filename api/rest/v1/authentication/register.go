package authentication_v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	errs "github.com/blacksmith-vish/sso/pkg/errors"
	"github.com/pkg/errors"
)

const (
	pathRegister string = "/register"
)

var (
	ErrUserExists error = errors.New("user already exists")
)

type RegisterRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

func (srv server) register() http.HandlerFunc {

	const errorMessage string = "registration failed"

	log := srv.log.With(
		slog.String("op", "authentication.register"),
	)

	// list of known errors compared to http codes
	errList := map[error]int{
		ErrUserExists: http.StatusNotAcceptable,
	}

	handler := errs.NewHandler(
		errs.WithMessage("registration failed"),
		errs.WithCodes(errList),
	)

	return func(w http.ResponseWriter, r *http.Request) {

		request := new(RegisterRequest)

		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			log.Error(err.Error())
		}

		serviceRequest := RegisterRequest{
			Nickname: r.Form.Get("nickname"),
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("pswd"),
		}

		_, err = srv.auth.Register(
			context.Background(),
			serviceRequest,
		)
		if err != nil {
			Err := handler.Handle(err)
			log.Error(errorMessage, Err.SlogAttrs()...)
			http.Error(w, Err.Error(), Err.Code())
			return
		}

		w.Write([]byte("Hello BITCH!"))
	}
}
