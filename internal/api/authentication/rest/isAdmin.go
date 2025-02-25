package authentication

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	"github.com/go-chi/chi/v5"
)

func (srv server) isAdmin() http.HandlerFunc {

	const op = "authentication.http.IsAdmin"

	return func(w http.ResponseWriter, r *http.Request) {

		userID := chi.URLParam(r, "user_id")

		log := srv.log.With(
			slog.String("op", op),
			slog.String("user_id", userID),
		)

		isAdmin, err := srv.auth.IsAdmin(r.Context(), userID)

		if err != nil {
			log.Error("failed to check admin status", slog.String("err", err.Error()))

			switch {
			case err == models.ErrUserInvalidID:
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
			case err == models.ErrUserNotFound:
				http.Error(w, "User not found", http.StatusNotFound)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		response := struct {
			IsAdmin bool `json:"isAdmin"`
		}{
			IsAdmin: isAdmin,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to encode response", slog.String("err", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
