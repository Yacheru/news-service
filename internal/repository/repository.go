package repository

type News interface {
}

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}
