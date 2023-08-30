package src

import (
	"encoding/csv"
	"io"
)

//go:generate mockgen -destination=csv_mock.go -package=src -source=csv.go
type (
	csvWriterIface interface {
		Write(record []string) error
		Flush()
	}

	csvReaderIface interface {
		Read() (record []string, err error)
	}

	csvHandlerIface interface {
		NewWriter(w io.Writer) csvWriterIface
		NewReader(r io.Reader) csvReaderIface
	}

	csvHandler struct{}
)

func (c *csvHandler) NewWriter(w io.Writer) csvWriterIface {
	return csv.NewWriter(w)
}

func (c *csvHandler) NewReader(r io.Reader) csvReaderIface {
	return csv.NewReader(r)
}
