package src

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang/mock/gomock"
)

var uc usecaseIface

//go:generate mockgen -destination=process_mock.go -package=src -source=process.go
type (
	usecaseIface interface {
		GetSampleAPIResourceRedirect(ctx context.Context, link []string) (data []UserData, err error)
		StoreAndReplaceUserDataToCSV(ctx context.Context, data []UserData, path string) (err error)
		SearchUserWithTags(ctx context.Context, tags []string, path string) (data []UserData, err error)
	}

	usecase struct {
		api     apiFetcherIface
		storage storageIface
	}
)

func newUsecase() usecaseIface {
	if uc != nil {
		return uc
	}

	api, err := newFetcher(&http.Client{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	storage := newStorage()

	u := &usecase{
		api:     api,
		storage: storage,
	}
	uc = u
	return u
}

func newMockUC(mockCtrl *gomock.Controller) *MockusecaseIface {
	mock := NewMockusecaseIface(mockCtrl)
	uc = mock
	return mock
}

func (u *usecase) GetSampleAPIResourceRedirect(ctx context.Context, link []string) (data []UserData, err error) {
	return u.api.getSampleAPIResourceRedirect(ctx, link)
}

func (u *usecase) StoreAndReplaceUserDataToCSV(ctx context.Context, data []UserData, path string) (err error) {
	return u.storage.storeAndReplaceUserDataToCSV(ctx, data, path)
}

func (u *usecase) SearchUserWithTags(ctx context.Context, tags []string, path string) (data []UserData, err error) {
	return u.storage.searchFromCSV(ctx, tags, path)
}
