package fastlinks

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	server "github.com/TerrenceHo/monorepo/fastlinks/adapters/http"
	"github.com/TerrenceHo/monorepo/fastlinks/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Hidebanner bool
	Port       int
}

func Start(conf Config) {

	healthService := services.NewHealthService()

	app := initServer(conf)

	rootController := server.NewRootController(healthService)
	rootController.Mount(app.Group(""))

	go func() {
		if err := app.Start(":" + strconv.Itoa(conf.Port)); err != http.ErrServerClosed {
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
