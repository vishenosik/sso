package dgraph

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"path"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/vishenosik/sso/internal/store/dgraph/components/users"
	"github.com/vishenosik/sso/internal/store/models"
	"github.com/vishenosik/sso/pkg/helpers/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

type Client struct {
	client *dgo.Dgraph
	users  *users.Users
}

type Config struct {
	Credentials config.Credentials
	GrpcServer  config.Server
}

func NewClientCtx(ctx context.Context, config Config) (*Client, error) {

	client, err := connect(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		users:  users.NewUsersStore(client),
	}, nil
}

func (cli *Client) Migrate(log *slog.Logger, migrations fs.FS) error {

	schemaFile, err := migrations.Open(path.Join("migrations", "dgraph", "schema.gql"))
	if err != nil {
		return err
	}
	defer schemaFile.Close()

	schema, err := io.ReadAll(schemaFile)
	if err != nil {
		return err
	}

	log.Debug("Migrating schema", slog.String("schema", string(schema)))

	op := &api.Operation{
		Schema: string(schema),
	}

	if err := cli.client.Alter(context.Background(), op); err != nil {
		return err
	}

	return nil
}

func connect(
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

func (cli *Client) SaveUser(ctx context.Context, id string, nickname string, email string, passHash []byte) error {
	return cli.users.SaveUser(ctx, id, nickname, email, passHash)
}

func (cli *Client) UserByEmail(ctx context.Context, email string) (models.User, error) {
	return cli.users.UserByEmail(ctx, email)
}

func (cli *Client) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return false, nil
}
