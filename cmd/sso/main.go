package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/blacksmith-vish/sso/internal/app"
	"github.com/blacksmith-vish/sso/internal/lib/config"
	config_yaml "github.com/blacksmith-vish/sso/internal/store/filesystem/config/yaml"
	"github.com/blacksmith-vish/sso/pkg/logger"
)

// @title           sso
// @version         0.0.1
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8080
// @BasePath  /
//
// @securityDefinitions.basic  BasicAuth
//
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	runServer()
}

func runServer() {
	// Инициализация конфига
	yaml := config_yaml.MustLoad()

	conf := config.NewConfig(yaml)

	log := logger.SetupLogger(conf.Env)

	log.Info("start app")

	// Инициализация приложения
	application := app.NewApp(log, conf)

	// Инициализация gRPC-сервер
	go application.GRPCServer.MustRun()

	go application.RESTServer.MustRun()

	// Graceful shut down

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	log.Info("app stopping", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), conf.RestConfig.Timeout)
	defer func() {
		// extra handling here
		cancel()
	}()

	application.GRPCServer.Stop()

	application.RESTServer.Stop(ctx)

	log.Info("app stopped")
}
