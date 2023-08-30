package src

import (
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestGetFromSource(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name     string
		wantData []UserData
		wantErr  bool
		mock     func()
	}{
		{
			name: "test1_success",
			wantData: []UserData{
				{
					ID: "1",
				},
			},
			wantErr: false,
			mock: func() {
				mock := newMockUC(mockCtrl)
				mock.EXPECT().GetSampleAPIResourceRedirect(gomock.Any(), []string{
					APILink1,
					APILink2,
				}).Return([]UserData{
					{
						ID: "1",
					},
				}, nil).Times(1)
			},
		},
		{
			name:    "test2_fail",
			wantErr: true,
			mock: func() {
				mock := newMockUC(mockCtrl)
				mock.EXPECT().GetSampleAPIResourceRedirect(gomock.Any(), []string{
					APILink1,
					APILink2,
				}).Return(nil, errors.New("err")).Times(1)
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := GetFromSource()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("GetFromSource() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestSetAndReplaceToCSV(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		data []UserData
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "test1_success",
			args: args{
				data: []UserData{
					{
						ID: "1",
					},
				},
				path: "data.csv",
			},
			wantErr: false,
			mock: func() {
				mock := newMockUC(mockCtrl)
				mock.EXPECT().StoreAndReplaceUserDataToCSV(gomock.Any(), []UserData{
					{
						ID: "1",
					},
				}, "data.csv").Return(nil).Times(1)
			},
		},
		{
			name: "test2_fail",
			args: args{
				data: []UserData{
					{
						ID: "1",
					},
				},
				path: "data.csv",
			},
			wantErr: true,
			mock: func() {
				mock := newMockUC(mockCtrl)
				mock.EXPECT().StoreAndReplaceUserDataToCSV(gomock.Any(), []UserData{
					{
						ID: "1",
					},
				}, "data.csv").Return(errors.New("err")).Times(1)
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := SetAndReplaceToCSV(tt.args.data, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("SetAndReplaceToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchFromCSV(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		tags []string
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantData []UserData
		wantErr  bool
		mock     func()
	}{
		{
			name: "test1_success",
			args: args{
				tags: []string{"a", "b"},
				path: "data.csv",
			},
			wantErr: false,
			wantData: []UserData{
				{
					ID: "1",
				},
			},
			mock: func() {
				mock := newMockUC(mockCtrl)
				mock.EXPECT().SearchUserWithTags(gomock.Any(), []string{"a", "b"}, "data.csv").Return([]UserData{
					{
						ID: "1",
					},
				}, nil).Times(1)
			},
		},
		{
			name: "test2_fail",
			args: args{
				tags: []string{"a", "b"},
				path: "data.csv",
			},
			wantErr: true,
			mock: func() {
				mock := newMockUC(mockCtrl)
				mock.EXPECT().SearchUserWithTags(gomock.Any(), []string{"a", "b"}, "data.csv").Return(nil, errors.New("err")).Times(1)
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := SearchFromCSV(tt.args.tags, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchFromCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("SearchFromCSV() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
