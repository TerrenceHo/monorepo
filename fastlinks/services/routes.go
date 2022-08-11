package services

type RoutesStore interface {
	Name() string
	Health() error
}

type RoutesService struct {
	store RoutesStore
}

func NewRoutesService(store RoutesStore) *RoutesService {
	return &RoutesService{
		store: store,
	}
}
