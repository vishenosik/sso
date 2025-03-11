package dgraph

import (
	"context"
	"fmt"

	"github.com/blacksmith-vish/sso/pkg/helpers/config"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

type dgraphClient struct {
	Client *dgo.Dgraph
}

type Config struct {
	Credentials config.Credentials
	Server      config.Server
}

func NewClient(
	ctx context.Context,
	config Config,
) *dgo.Dgraph {

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	connection, err := grpc.NewClient(addr, opts...)
	if err != nil {
		// TODO: handle error
	}

	client := dgo.NewDgraphClient(
		api.NewDgraphClient(connection),
	)

	err = client.Login(ctx, config.Credentials.User, config.Credentials.Password)
	if err != nil {
		// TODO: handle error
	}

	return client
}
