package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/blacksmith-vish/sso/internal/app/config"

	authentication_v1 "github.com/blacksmith-vish/sso/internal/gen/grpc/v1/authentication"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Conf       *config.Config
	AuthClient authentication_v1.AuthenticationClient
}

func newConfig() *config.Config {
	return config.EnvConfig()
}

func New(t *testing.T) (context.Context, *Suite) {

	t.Helper()
	t.Parallel()

	conf := newConfig()

	ctx, cancelCtx := context.WithTimeout(context.Background(), conf.AuthenticationService.TokenTTL)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(
		grpcAddress(conf),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Conf:       conf,
		AuthClient: authentication_v1.NewAuthenticationClient(cc),
	}

}

func grpcAddress(conf *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(int(conf.GrpcConfig.Port)))
}
