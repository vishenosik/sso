package main

import (
	// std
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	// internal pkg
	"github.com/vishenosik/gocherry"
	"github.com/vishenosik/gocherry/pkg/cache"
	_ctx "github.com/vishenosik/gocherry/pkg/context"
	"github.com/vishenosik/gocherry/pkg/grpc"
	_http "github.com/vishenosik/gocherry/pkg/http"
	"github.com/vishenosik/gocherry/pkg/sql"
	"github.com/vishenosik/sso-sdk/api"

	// internal
	embed "github.com/vishenosik/sso"
	"github.com/vishenosik/sso/internal/dto"
	"github.com/vishenosik/sso/internal/services"
	"github.com/vishenosik/sso/internal/store/sql/sqlite"
)

var (
	BuildDate string
	GitBranch string
	GitCommit string
	GoVersion string
	GitTag    string
)

func main() {

	flag.BoolFunc("v", "Show build info", func(s string) error {
		defer os.Exit(0)
		printBuildInfo(os.Stdout)
		return nil
	})

	gocherry.ConfigFlags(
		services.AuthenticationConfigEnv{},
	)

	flag.Parse()
	ctx := context.Background()

	// App init
	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	err = app.Start(ctx)
	if err != nil {
		panic(err)
	}

	// Graceful shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopctx, cancel := context.WithTimeout(
		_ctx.WithStopCtx(ctx, <-stop),
		time.Second*5,
	)
	defer cancel()

	if err := app.Stop(stopctx); err != nil {
		panic(err)
	}
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type App struct {
	Server
}

func NewApp() (*App, error) {

	// Cache init

	cacheStore, err := cache.NewRedisCache()
	if err != nil {
		cacheStore = cache.NewNoopCache()
	}

	// Stores init

	sqlStore, err := sql.NewSqliteStore(
		sql.WithMigration(embed.Migrations, path.Join(embed.MigrationsPath, "sqlite")),
	)
	if err != nil {
		return nil, err
	}

	db, err := sqlStore.Open(context.TODO())
	if err != nil {
		return nil, err
	}

	// client, err := graph.NewClientCtx(context.TODO(), graph.DgraphConfig{})
	// if err != nil {
	// 	return nil, err
	// }

	// if err := client.Migrate(embed.Migrations); err != nil {
	// 	return nil, err
	// }

	// graphUsers := dgraph.NewUsersStore(client.Cli)

	usersStore := sqlite.NewUsersStore(db)
	appsStore := sqlite.NewAppsStore(db)

	// Usecases init
	authService, err := services.NewAuthenticationService(usersStore, usersStore, appsStore)
	if err != nil {
		return nil, err
	}

	authDTO := dto.NewAuthenticationDTO(authService)
	authApi := api.NewAuthenticationApi(authDTO)

	sys := services.NewSystem(100, false, 322)
	systemDTO := dto.NewSystemDTO(sys)
	systemApi := api.NewSystemApi(systemDTO)

	// Services init
	httpServer, err := _http.NewHttpServer(api.NewHttpHandler(
		authApi,
		systemApi,
	))
	if err != nil {
		return nil, err
	}

	grpcServer, err := grpc.NewGrpcServer(grpc.GrpcServices{
		authApi,
	})
	if err != nil {
		return nil, err
	}

	app, err := gocherry.NewApp()
	if err != nil {
		return nil, err
	}

	app.AddServices(
		httpServer,
		cacheStore,
		sqlStore,
		grpcServer,
		// client,
	)

	return &App{
		Server: app,
	}, nil
}

func printBuildInfo(writer io.Writer) {
	writer.Write([]byte("Build Info:\n"))
	writer.Write([]byte(fmt.Sprintf("Build Date: %s\n", BuildDate)))
	writer.Write([]byte(fmt.Sprintf("Git Branch: %s\n", GitBranch)))
	writer.Write([]byte(fmt.Sprintf("Git Commit: %s\n", GitCommit)))
	writer.Write([]byte(fmt.Sprintf("Go Version: %s\n", GoVersion)))
	writer.Write([]byte(fmt.Sprintf("Git Tag: %s\n", GitTag)))
}
