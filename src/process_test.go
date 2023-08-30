package src

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_newUsecase(t *testing.T) {
	uc = nil
	mockAPI, _ := newFetcher(&http.Client{
		Timeout: 10 * time.Second,
	})
	mockStorage := newStorage()
	tests := []struct {
		name string
		want usecaseIface
	}{
		{
			name: "test1_success",
			want: &usecase{
				api:     mockAPI,
				storage: mockStorage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUsecase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_usecase_GetSampleAPIResourceRedirect(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		api     apiFetcherIface
		storage storageIface
	}
	type args struct {
		ctx  context.Context
		link []string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []UserData
		wantErr  bool
	}{
		{
			name: "test1_success",
			args: args{
				ctx:  context.Background(),
				link: []string{"a", "b"},
			},
			wantData: []UserData{
				{
					ID:      "12",
					Balance: "100",
					Tags:    []string{"tag"},
				},
			},
			wantErr: false,
			fields: fields{
				api: func() apiFetcherIface {
					mock := NewMockapiFetcherIface(mockCtrl)
					mock.EXPECT().getSampleAPIResourceRedirect(gomock.Any(), []string{"a", "b"}).Return([]UserData{
						{
							ID:      "12",
							Balance: "100",
							Tags:    []string{"tag"},
						},
					}, nil).Times(1)
					return mock
				}(),
			},
		},
		{
			name: "test2_fail",
			args: args{
				ctx:  context.Background(),
				link: []string{"a", "b"},
			},
			wantErr: true,
			fields: fields{
				api: func() apiFetcherIface {
					mock := NewMockapiFetcherIface(mockCtrl)
					mock.EXPECT().getSampleAPIResourceRedirect(gomock.Any(), []string{"a", "b"}).Return(nil, errors.New("err")).Times(1)
					return mock
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &usecase{
				api:     tt.fields.api,
				storage: tt.fields.storage,
			}
			gotData, err := u.GetSampleAPIResourceRedirect(tt.args.ctx, tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetSampleAPIResourceRedirect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("usecase.GetSampleAPIResourceRedirect() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func Test_usecase_StoreAndReplaceUserDataToCSV(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		api     apiFetcherIface
		storage storageIface
	}
	type args struct {
		ctx  context.Context
		data []UserData
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1_success",
			args: args{
				ctx: context.Background(),
				data: []UserData{
					{
						ID: "12",
					},
				},
				path: "a",
			},
			wantErr: false,
			fields: fields{
				storage: func() storageIface {
					mock := NewMockstorageIface(mockCtrl)
					mock.EXPECT().storeAndReplaceUserDataToCSV(gomock.Any(), []UserData{
						{
							ID: "12",
						},
					}, "a").Return(nil).Times(1)
					return mock
				}(),
			},
		},
		{
			name: "test2_fail",
			args: args{
				ctx: context.Background(),
				data: []UserData{
					{
						ID: "12",
					},
				},
				path: "a",
			},
			wantErr: true,
			fields: fields{
				storage: func() storageIface {
					mock := NewMockstorageIface(mockCtrl)
					mock.EXPECT().storeAndReplaceUserDataToCSV(gomock.Any(), []UserData{
						{
							ID: "12",
						},
					}, "a").Return(errors.New("err")).Times(1)
					return mock
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &usecase{
				api:     tt.fields.api,
				storage: tt.fields.storage,
			}
			if err := u.StoreAndReplaceUserDataToCSV(tt.args.ctx, tt.args.data, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("usecase.StoreAndReplaceUserDataToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_usecase_SearchUserWithTags(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		api     apiFetcherIface
		storage storageIface
	}
	type args struct {
		ctx  context.Context
		tags []string
		path string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []UserData
		wantErr  bool
	}{
		{
			name: "test1_success",
			args: args{
				ctx:  context.Background(),
				tags: []string{"a", "b"},
				path: "a",
			},
			wantData: []UserData{
				{
					ID: "12",
				},
			},
			wantErr: false,
			fields: fields{
				storage: func() storageIface {
					mock := NewMockstorageIface(mockCtrl)
					mock.EXPECT().searchFromCSV(gomock.Any(), []string{"a", "b"}, "a").Return([]UserData{
						{
							ID: "12",
						},
					}, nil).Times(1)
					return mock
				}(),
			},
		},
		{
			name: "test2_fail",
			args: args{
				ctx:  context.Background(),
				tags: []string{"a", "b"},
				path: "a",
			},
			wantErr: true,
			fields: fields{
				storage: func() storageIface {
					mock := NewMockstorageIface(mockCtrl)
					mock.EXPECT().searchFromCSV(gomock.Any(), []string{"a", "b"}, "a").Return(nil, errors.New("err")).Times(1)
					return mock
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &usecase{
				api:     tt.fields.api,
				storage: tt.fields.storage,
			}
			gotData, err := u.SearchUserWithTags(tt.args.ctx, tt.args.tags, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SearchUserWithTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("usecase.SearchUserWithTags() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
