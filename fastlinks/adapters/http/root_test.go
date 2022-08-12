package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TerrenceHo/monorepo/fastlinks/services"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// func TestHello(t *testing.T) {
// 	assert := assert.New(t)

// 	healthService := services.NewHealthService()
// 	rootController := NewRootController(healthService)

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetPath("/")

// 	if assert.NoError(rootController.hello(c)) {
// 		assert.Equal(http.StatusOK, rec.Code)
// 		assert.Equal("Hello, World!", rec.Body.String())
// 	}
// }

type MockHealthService struct {
	HealthCheckErrors []services.HealthCheckError
}

func (m *MockHealthService) Check() []services.HealthCheckError {
	return m.HealthCheckErrors
}

func TestHealth(t *testing.T) {
	assert := assert.New(t)

	type testcase struct {
		healthService services.IHealthService
		want          string
		status        int
	}

	testcases := []testcase{
		{
			healthService: services.NewHealthService(),
			want:          "{\"status\":\"OK\",\"healthCheckErrors\":[]}\n",
			status:        http.StatusOK,
		},
		{
			healthService: &MockHealthService{
				HealthCheckErrors: []services.HealthCheckError{
					{
						Name:  "Fake Health Check",
						Error: stackerrors.New("fake health check"),
					},
				},
			},
			want:   "{\"status\":\"UNHEALTHY\",\"healthCheckErrors\":[{\"Name\":\"Fake Health Check\",\"Error\":{}}]}\n",
			status: http.StatusInternalServerError,
		},
	}

	for _, testcase := range testcases {
		rootController := NewRootController(nil, testcase.healthService)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/health")

		if assert.NoError(rootController.health(c)) {
			assert.Equal(testcase.status, rec.Code)
			assert.Equal(testcase.want, rec.Body.String())
		}
	}

}
