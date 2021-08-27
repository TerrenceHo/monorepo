package http

import (
	"net/http"

	"github.com/TerrenceHo/monorepo/fastlinks/services"
	"github.com/labstack/echo/v4"
)

type STATUS string

const (
	OK        STATUS = "OK"
	UNHEALTHY STATUS = "UNHEALTHY"
)

type RootController struct {
	healthService services.IHealthService
}

func NewRootController(healthService services.IHealthService) *RootController {
	return &RootController{
		healthService: healthService,
	}
}

func (rc *RootController) Mount(g *echo.Group) {
	g.GET("/", rc.hello)
	g.GET("/hello", rc.hello)
	g.GET("/health", rc.health)
}

func (rc *RootController) hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type HealthResponse struct {
	Status            STATUS                      `json:"status"`
	HealthCheckErrors []services.HealthCheckError `json:"healthCheckErrors"`
}

func (rc *RootController) health(c echo.Context) error {
	healthCheckErrors := rc.healthService.Check()
	var httpStatus int
	var hr *HealthResponse
	if len(healthCheckErrors) == 0 {
		hr = &HealthResponse{
			Status:            OK,
			HealthCheckErrors: healthCheckErrors,
		}
		httpStatus = http.StatusOK
	} else {
		hr = &HealthResponse{
			Status:            UNHEALTHY,
			HealthCheckErrors: healthCheckErrors,
		}
		httpStatus = http.StatusInternalServerError
	}
	return c.JSON(httpStatus, hr)
}
