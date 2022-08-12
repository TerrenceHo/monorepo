package fastlinks

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	server "github.com/TerrenceHo/monorepo/fastlinks/adapters/http"
	"github.com/TerrenceHo/monorepo/fastlinks/services"
	"github.com/TerrenceHo/monorepo/fastlinks/stores/bbolt"
	// "github.com/TerrenceHo/monorepo/fastlinks/stores/postgresql"
	"github.com/TerrenceHo/monorepo/fastlinks/views"
	"github.com/TerrenceHo/monorepo/utils-go/file"
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
	Storage    string
	DB         DBConfig
	Local      LocalConfig
}

type DBConfig struct {
	User     string
	Password string
	DBName   string
	Port     string
	Host     string
	SSLMode  string
}

type LocalConfig struct {
	File string
}

//go:embed static
var staticContent embed.FS

func Start(conf Config) {
	// Instantiate logger
	logger, err := logging.ConfigureLogger(loggerType(conf.Env))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to instantiate logger: %v", err)
		os.Exit(1)
	}
	logging.SetGlobalLogger(logger)

	// Instantiate database connections and stores
	var routesStore services.RoutesStore
	switch conf.Storage {
	case "db":
		// db, err := postgresql.NewConnection(
		// 	conf.DB.User,
		// 	conf.DB.Password,
		// 	conf.DB.DBName,
		// 	conf.DB.Port,
		// 	conf.DB.Host,
		// 	conf.DB.SSLMode,
		// )
		// defer db.Close()
		// if err != nil {
		// 	logging.Fatal(
		// 		"failed to connect to database",
		// 		zap.Error(err),
		// 	)
		// }
		// routesStore = postgresql.NewRoutesStore(db)
	default:
		err := os.MkdirAll(filepath.Dir(conf.Local.File), 0750)
		if err != nil {
			logging.Fatal(
				"failed to create local bbolt database directory",
				zap.Error(err),
			)
		}
		err = file.CreateFileIfNotExists(conf.Local.File)
		if err != nil {
			logging.Fatal(
				"failed to create local bbolt database file",
				zap.Error(err),
			)
		}
		db, err := bbolt.NewConnection(conf.Local.File, os.FileMode(0666))
		defer db.Close()
		if err != nil {
			logging.Fatal(
				"failed to open local bbolt database file",
				zap.Error(err),
			)
		}
		routesStore = bbolt.NewRoutesStore(db, "fastlinks-bucket")
		// create bucket if it doesn't exist
		if err = routesStore.Migrate(); err != nil {
			logging.Fatal(
				"failed to migrate/create local bbolt database bucket fastlinks-bucket",
				zap.Error(err),
			)
		}
	}

	// create services
	healthService := services.NewHealthService(routesStore)
	routesService := services.NewRoutesService(routesStore)

	// create http controllers
	rootController := server.NewRootController(routesService, healthService)

	// create views
	rootView, err := views.NewView(staticContent, "static/*.html", "static/root/*.html")
	if err != nil {
		logging.Fatal("failed to initiate root templates", zap.Error(err))
	}
	newView, err := views.NewView(staticContent, "static/*.html", "static/new/*.html")
	if err != nil {
		logging.Fatal("failed to initiate new templates", zap.Error(err))
	}

	renderer := views.NewRenderer(rootView, newView)

	// start servers and attach http routes to server
	app := initServer(conf)
	app.Renderer = renderer
	rootController.Mount(app.Group(""))

	go func() {
		if err := app.Start(conf.Host + ":" + conf.Port); err != http.ErrServerClosed {
			logging.Fatal("failed to start server", zap.Error(err))
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
