package services

import (
	"github.com/TerrenceHo/monorepo/fastlinks/models"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
)

type RoutesStore interface {
	Name() string
	Health() error
	GetByKey(key string) (*models.Route, error)
	GetAll() ([]*models.Route, error)
	Put(route *models.Route) error
	Delete(key string) error
	Migrate() error
}

type RoutesService struct {
	store RoutesStore
}

func NewRoutesService(store RoutesStore) *RoutesService {
	return &RoutesService{
		store: store,
	}
}

func (rs *RoutesService) GetAllRoutes() ([]*models.Route, error) {
	routes, err := rs.store.GetAll()
	if err != nil {
		return nil, stackerrors.Wrap(err, "failed to get all routes")
	}
	return routes, nil
}

func (rs *RoutesService) GetRoute(key string) (*models.Route, error) {
	route, err := rs.store.GetByKey(key)
	if err != nil {
		return nil, stackerrors.Wrapf(err, "failed to get route with key %s", key)
	}
	return route, nil
}

func (rs *RoutesService) CreateRoute(key string, RedirectURL string, ExtendedURL string) error {
	// validations?
	route := &models.Route{
		Key:         key,
		RedirectURL: RedirectURL,
		ExtendedURL: ExtendedURL,
	}
	err := rs.store.Put(route)
	if err != nil {
		return stackerrors.Wrapf(
			err, "failed to create route with key(%s), redirect URL(%s), extended url(%s)",
			key, RedirectURL, ExtendedURL,
		)
	}

	return nil
}

func (rs *RoutesService) DeleteRoute(key string) error {
	err := rs.store.Delete(key)
	if err != nil {
		return stackerrors.Wrapf(err, "failed to delete route with key %s", key)
	}
	return nil
}
