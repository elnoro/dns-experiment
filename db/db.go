package db

type HostDb interface {
	Save(host string) error
	Delete(host string) error
	Get(host string) (bool, error)
}
