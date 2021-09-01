package http

import (
	"net/http"

	"github.com/TerrenceHo/monorepo/fastlinks/services"
	"github.com/labstack/echo/v4"
)

type RoutesService interface{}

type RoutesController struct {
	routesService RoutesService
	healthService services.IHealthService
}

func NewRoutesController(routesService RoutesService, healthService services.IHealthService) *RoutesController {
	return &RoutesController{
		routesService: routesService,
		healthService: healthService,
	}
}

func (rc *RoutesController) Mount(g *echo.Group) {
	g.GET("/", rc.hello)
}

func (rc *RoutesController) hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello Routes!")
}
