package dgraph

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/vishenosik/sso/pkg/helpers/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

type dgraphClient struct {
	Client *dgo.Dgraph
}

type Config struct {
	Credentials config.Credentials
	GrpcServer  config.Server
}

func NewClient(
	ctx context.Context,
	config Config,
) (*dgo.Dgraph, error) {

	addr := fmt.Sprintf("%s:%d", config.GrpcServer.Host, config.GrpcServer.Port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	connection, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}

	client := dgo.NewDgraphClient(
		api.NewDgraphClient(connection),
	)

	err = client.Login(ctx, config.Credentials.User, config.Credentials.Password)
	if err != nil {
		return nil, err
	}

	return client, nil
}
