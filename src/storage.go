package src

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

//go:generate mockgen -destination=storage_mock.go -package=src -source=storage.go
type (
	storageIface interface {
		storeAndReplaceUserDataToCSV(ctx context.Context, data []UserData, path string) error
		searchFromCSV(ctx context.Context, tags []string, path string) (data []UserData, err error)
	}

	storage struct {
		fileReader fReaderIface
		csvHandler csvHandlerIface
	}
)

var (
	ErrMissingFile = errors.New("missing file")
)

func newStorage() storageIface {
	return &storage{
		fileReader: &fileHandler{},
		csvHandler: &csvHandler{},
	}
}

func (s *storage) storeAndReplaceUserDataToCSV(ctx context.Context, data []UserData, path string) error {
	if path == "" {
		path = "data.csv"
	}
	// file, err := os.Create(path)
	file, err := s.fileReader.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := s.csvHandler.NewWriter(file)
	defer writer.Flush()
	// writer := csv.NewWriter(file)
	for _, v := range data {
		tagBytes, err := json.Marshal(&v.Tags)
		if err != nil {
			return err
		}
		err = writer.Write([]string{v.ID, strconv.FormatBool(v.ActiveStatus), v.Balance, string(tagBytes)})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *storage) searchFromCSV(ctx context.Context, tags []string, path string) (data []UserData, err error) {
	if path == "" {
		path = "data.csv"
	}

	file, err := s.fileReader.Open(path)
	if err != nil {
		return nil, ErrMissingFile
	}
	defer file.Close()

	csvReader := s.csvHandler.NewReader(bufio.NewReader(file))
	for {
		res, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(res) < 4 {
			return nil, errors.New("bad csv row format")
		}

		var rowTags []string
		err = json.Unmarshal([]byte(res[3]), &rowTags)
		if err != nil {
			return nil, err
		}

		tagMap := make(map[string]struct{})
		for _, v := range rowTags {
			tagMap[v] = struct{}{}
		}

		shouldAppend := true
		for _, v := range tags {
			if _, ok := tagMap[v]; !ok {
				shouldAppend = false
				break
			}
		}
		if !shouldAppend {
			continue
		}

		data = append(data, UserData{
			ID:      res[0],
			Balance: res[2],
		})
	}

	return
}
