package services

type RoutesStore interface{}

type RoutesService struct {
	store RoutesStore
}

func NewRoutesService(store RoutesStore) *RoutesService {
	return &RoutesService{
		store: store,
	}
}
