package grpcApp

import (
	"io"
	"log/slog"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/blacksmith-vish/sso/internal/api/authentication/grpc/mocks"
	"github.com/blacksmith-vish/sso/internal/lib/config"
)

func Test_App_MustRun_PanicsOnError(t *testing.T) {
	t.Helper()
	t.Parallel()

	log := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}))
	conf := config.GRPCConfig{} // Invalid port to force an error
	conf.Port = 0
	authService := mocks.NewAuthentication(t) // Assuming a mock implementation

	app := NewGrpcApp(log, conf, authService)

	assert.Panics(t, func() {
		app.MustRun()
	})
}
