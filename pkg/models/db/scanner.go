package db

type scanner interface {
	Scan(dest ...interface{}) error
}
