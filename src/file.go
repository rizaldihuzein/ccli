package src

import "os"

//go:generate mockgen -destination=file_mock.go -package=src -source=file.go
type (
	fReaderIface interface {
		Create(name string) (*os.File, error)
		Open(name string) (*os.File, error)
	}

	fileHandler struct{}
)

func (f *fileHandler) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f *fileHandler) Open(name string) (*os.File, error) {
	return os.Open(name)
}
