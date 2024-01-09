package storage

type DB interface {
	GetLast() (string, error)
	Update(string) error
	Connect()
	PathId(string) (string, error)
}
