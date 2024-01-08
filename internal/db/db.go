package db

type db interface {
	GetLast() (string, error)
	Update(string) error
}
