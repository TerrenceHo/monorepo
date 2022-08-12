package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TerrenceHo/monorepo/fastlinks/models"
	"github.com/TerrenceHo/monorepo/fastlinks/services"
	"github.com/TerrenceHo/monorepo/fastlinks/views"
	"github.com/labstack/echo/v4"
)

type RoutesService interface {
	GetAllRoutes() ([]*models.Route, error)
	GetRoute(key string) (*models.Route, error)
	CreateRoute(key string, RedirectURL string, ExtendedURL string) error
	DeleteRoute(key string) error
}

type STATUS string

const (
	OK        STATUS = "OK"
	UNHEALTHY STATUS = "UNHEALTHY"
)

type RootController struct {
	routesService RoutesService
	healthService services.IHealthService
}

func NewRootController(service RoutesService, healthService services.IHealthService) *RootController {
	return &RootController{
		routesService: service,
		healthService: healthService,
	}
}

func (rc *RootController) Mount(g *echo.Group) {
	g.GET("/", rc.getAll)
	g.GET("/__new", rc.newPage)
	g.GET("/:key", rc.redirectRoute)
	g.GET("/:key/:extended", rc.redirectExtendedRoute)
	g.POST("/", rc.createRoute)
	g.DELETE("/:key", rc.deleteRoute)
	g.GET("/health", rc.health)
}

func (rc *RootController) getAll(c echo.Context) error {
	routes, err := rc.routesService.GetAllRoutes()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	td := views.TemplateData{
		Type: views.ROOT_VIEW,
		Data: routes,
	}
	return c.Render(http.StatusOK, "page", td)
}

func (rc *RootController) newPage(c echo.Context) error {
	td := views.TemplateData{
		Type: views.NEW_VIEW,
	}
	return c.Render(http.StatusOK, "page", td)
}

func (rc *RootController) redirectRoute(c echo.Context) error {
	key := c.Param("key")
	fmt.Println(key)
	route, err := rc.routesService.GetRoute(key)
	if err != nil {
		return c.Redirect(http.StatusFound, "/__new")
	}
	return c.Redirect(http.StatusFound, route.RedirectURL)
}

func (rc *RootController) redirectExtendedRoute(c echo.Context) error {
	key := c.Param("key")
	extendedKey := c.Param("extended")
	fmt.Println(key, extendedKey)
	route, err := rc.routesService.GetRoute(key)
	if err != nil {
		return c.Redirect(http.StatusFound, "/__new")
	}
	return c.Redirect(http.StatusFound, strings.Replace(route.ExtendedURL, "{}", extendedKey, 1))
}

func (rc *RootController) createRoute(c echo.Context) error {
	key := c.FormValue("Key")
	redirectURL := c.FormValue("RedirectURL")
	extendedURL := c.FormValue("ExtendedURL")
	fmt.Printf("key: %s, redirect: %s, extended: %s\n", key, redirectURL, extendedURL)
	err := rc.routesService.CreateRoute(key, redirectURL, extendedURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.Redirect(http.StatusFound, redirectURL)
}

func (rc *RootController) deleteRoute(c echo.Context) error {
	key := c.Param("key")
	if err := rc.routesService.DeleteRoute(key); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
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
