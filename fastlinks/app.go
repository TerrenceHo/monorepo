package fastlinks

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	server "github.com/TerrenceHo/monorepo/fastlinks/adapters/http"
	"github.com/TerrenceHo/monorepo/fastlinks/services"
	"github.com/TerrenceHo/monorepo/fastlinks/stores/postgresql"
	"github.com/TerrenceHo/monorepo/utils-go/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Config struct {
	Hidebanner bool
	Env        string
	Port       string
	Host       string
	DB         DBConfig
}

type DBConfig struct {
	User     string
	Password string
	DBName   string
	Port     string
	Host     string
	SSLMode  string
}

func Start(conf Config) {
	// Instantiate logger
	logger, err := logging.ConfigureLogger(loggerType(conf.Env))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate logger: %v", err)
		os.Exit(1)
	}
	logging.SetGlobalLogger(logger)

	// Instantiate database connections
	db, err := postgresql.NewConnection(
		conf.DB.User,
		conf.DB.Password,
		conf.DB.DBName,
		conf.DB.Port,
		conf.DB.Host,
		conf.DB.SSLMode,
	)
	defer db.Close()
	if err != nil {
		logging.Fatal(
			"failed to connect to database",
			zap.Error(err),
		)
	}

	// create stores
	routesStore := postgresql.NewRoutesStore(db)

	// create services
	healthService := services.NewHealthService(routesStore)
	routesService := services.NewRoutesService(routesStore)

	// create http controllers
	rootController := server.NewRootController(healthService)
	routesController := server.NewRoutesController(routesService, healthService)

	// start servers and attach http routes to server
	app := initServer(conf)
	rootController.Mount(app.Group(""))
	routesController.Mount(app.Group("routes"))

	go func() {
		if err := app.Start(":" + conf.Port); err != http.ErrServerClosed {
			app.Logger.Fatal("Shutting down the server.")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}
	fmt.Println("Shutdown fastlinks server. Goodbye.")
}

func initServer(c Config) *echo.Echo {
	app := echo.New()
	app.HideBanner = true

	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())
	app.Use(middleware.RemoveTrailingSlash())
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${method} {id":"${id}","remote_ip":"${remote_ip}",` +
			`"uri":"${uri}","status":${status},"latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\n",
	}))

	return app
}

func loggerType(env string) logging.LoggerType {
	switch env {
	case "dev":
		return logging.DevLogger
	case "prod":
		return logging.ProdLogger
	case "test":
		return logging.TestLogger
	default:
		return logging.TestLogger
	}
}
