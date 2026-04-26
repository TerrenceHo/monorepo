package services

type HealthStore interface {
	Name() string
	Health() error
}

type IHealthService interface {
	Check() []HealthCheckError
}

type HealthService struct {
	stores []HealthStore
}

type HealthCheckError struct {
	Name  string
	Error error
}

func NewHealthService(stores ...HealthStore) *HealthService {
	return &HealthService{
		stores: stores,
	}
}

func (hs *HealthService) Check() []HealthCheckError {
	healthErrors := []HealthCheckError{}
	for _, store := range hs.stores {
		if err := store.Health(); err != nil {
			name := store.Name()
			healthErrors = append(healthErrors, HealthCheckError{
				Name:  name,
				Error: err,
			})
		}
	}
	return healthErrors
}
