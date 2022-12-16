package store

type StoreInterface interface {
	//Event repository
	Event() EventRepository
}
