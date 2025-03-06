package dgraph

import (
	"context"
	"fmt"

	config "github.com/blacksmith-vish/sso/internal/store/config_test"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

type dgraphClient struct {
	Client *dgo.Dgraph
}

func NewClient(
	ctx context.Context,
) *dgraphClient {

	conf := loadConfig(config.NewFSConfig[dgraphConfig](""))

	return &dgraphClient{
		Client: newClient(ctx, conf),
	}
}

func newClient(
	ctx context.Context,
	conf dgraphConfig,
) *dgo.Dgraph {

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

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

	err = client.Login(ctx, conf.User, conf.Password)
	if err != nil {
		// TODO: handle error
	}

	return client
}
