package src

import (
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_newStorage(t *testing.T) {
	tests := []struct {
		name string
		want storageIface
	}{
		{
			name: "test1_success",
			want: &storage{
				fileReader: &fileHandler{},
				csvHandler: &csvHandler{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_storeAndReplaceUserDataToCSV(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		fileReader fReaderIface
		csvHandler csvHandlerIface
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
						ID:           "1",
						ActiveStatus: true,
						Balance:      "1000",
						Tags:         []string{"a", "b"},
					},
				},
				path: "dd.csv",
			},
			fields: fields{
				fileReader: func() fReaderIface {
					mock := NewMockfReaderIface(mockCtrl)
					mock.EXPECT().Create("dd.csv").Return(&os.File{}, nil).Times(1)
					return mock
				}(),
				csvHandler: func() csvHandlerIface {
					mockWriter := NewMockcsvWriterIface(mockCtrl)
					mock := NewMockcsvHandlerIface(mockCtrl)
					mock.EXPECT().NewWriter(gomock.Any()).Return(mockWriter).Times(1)
					mockWriter.EXPECT().Write([]string{"1", "true", "1000", "[\"a\",\"b\"]"}).Return(nil).Times(1)
					mockWriter.EXPECT().Flush().Times(1)
					return mock
				}(),
			},
		},
		{
			name: "test2_fail_open",
			args: args{
				ctx: context.Background(),
				data: []UserData{
					{
						ID:           "1",
						ActiveStatus: true,
						Balance:      "1000",
						Tags:         []string{"a", "b"},
					},
				},
				path: "dd.csv",
			},
			wantErr: true,
			fields: fields{
				fileReader: func() fReaderIface {
					mock := NewMockfReaderIface(mockCtrl)
					mock.EXPECT().Create("dd.csv").Return(nil, errors.New("err")).Times(1)
					return mock
				}(),
				csvHandler: func() csvHandlerIface {
					mock := NewMockcsvHandlerIface(mockCtrl)
					return mock
				}(),
			},
		},
		{
			name: "test3_fail_write",
			args: args{
				ctx: context.Background(),
				data: []UserData{
					{
						ID:           "1",
						ActiveStatus: true,
						Balance:      "1000",
						Tags:         []string{"a", "b"},
					},
				},
				path: "dd.csv",
			},
			wantErr: true,
			fields: fields{
				fileReader: func() fReaderIface {
					mock := NewMockfReaderIface(mockCtrl)
					mock.EXPECT().Create("dd.csv").Return(&os.File{}, nil).Times(1)
					return mock
				}(),
				csvHandler: func() csvHandlerIface {
					mockWriter := NewMockcsvWriterIface(mockCtrl)
					mock := NewMockcsvHandlerIface(mockCtrl)
					mock.EXPECT().NewWriter(gomock.Any()).Return(mockWriter).Times(1)
					mockWriter.EXPECT().Write([]string{"1", "true", "1000", "[\"a\",\"b\"]"}).Return(errors.New("err")).Times(1)
					mockWriter.EXPECT().Flush().Times(1)
					return mock
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				fileReader: tt.fields.fileReader,
				csvHandler: tt.fields.csvHandler,
			}
			if err := s.storeAndReplaceUserDataToCSV(tt.args.ctx, tt.args.data, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("storage.storeAndReplaceUserDataToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_storage_searchFromCSV(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		fileReader fReaderIface
		csvHandler csvHandlerIface
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
				path: "a.csv",
			},
			wantData: []UserData{
				{
					ID:      "1",
					Balance: "1000",
				},
			},
			wantErr: false,
			fields: fields{
				fileReader: func() fReaderIface {
					mock := NewMockfReaderIface(mockCtrl)
					mock.EXPECT().Open("a.csv").Return(&os.File{}, nil).Times(1)
					return mock
				}(),
				csvHandler: func() csvHandlerIface {
					mock := NewMockcsvHandlerIface(mockCtrl)
					mockReader := NewMockcsvReaderIface(mockCtrl)
					mock.EXPECT().NewReader(gomock.Any()).Return(mockReader).Times(1)
					mockReader.EXPECT().Read().Return([]string{
						"1", "true", "1000", "[\"a\",\"b\"]",
					}, nil).Times(1)
					mockReader.EXPECT().Read().Return(nil, io.EOF).Times(1)
					return mock
				}(),
			},
		},
		{
			name: "test2_success",
			args: args{
				ctx:  context.Background(),
				path: "a.csv",
			},
			wantData: []UserData{
				{
					ID:      "1",
					Balance: "1000",
				},
			},
			wantErr: false,
			fields: fields{
				fileReader: func() fReaderIface {
					mock := NewMockfReaderIface(mockCtrl)
					mock.EXPECT().Open("a.csv").Return(&os.File{}, nil).Times(1)
					return mock
				}(),
				csvHandler: func() csvHandlerIface {
					mock := NewMockcsvHandlerIface(mockCtrl)
					mockReader := NewMockcsvReaderIface(mockCtrl)
					mock.EXPECT().NewReader(gomock.Any()).Return(mockReader).Times(1)
					mockReader.EXPECT().Read().Return([]string{
						"1", "true", "1000", "[\"a\",\"b\"]",
					}, nil).Times(1)
					mockReader.EXPECT().Read().Return(nil, io.EOF).Times(1)
					return mock
				}(),
			},
		},
		{
			name: "test3_fail_open",
			args: args{
				ctx:  context.Background(),
				tags: []string{"a", "b"},
				path: "a.csv",
			},
			wantErr: true,
			fields: fields{
				fileReader: func() fReaderIface {
					mock := NewMockfReaderIface(mockCtrl)
					mock.EXPECT().Open("a.csv").Return(nil, errors.New("err")).Times(1)
					return mock
				}(),
				csvHandler: func() csvHandlerIface {
					mock := NewMockcsvHandlerIface(mockCtrl)
					return mock
				}(),
			},
		},
		{
			name: "test4_fail_read",
			args: args{
				ctx:  context.Background(),
				tags: []string{"a", "b"},
				path: "a.csv",
			},
			wantErr: true,
			fields: fields{
				fileReader: func() fReaderIface {
					mock := NewMockfReaderIface(mockCtrl)
					mock.EXPECT().Open("a.csv").Return(&os.File{}, nil).Times(1)
					return mock
				}(),
				csvHandler: func() csvHandlerIface {
					mock := NewMockcsvHandlerIface(mockCtrl)
					mockReader := NewMockcsvReaderIface(mockCtrl)
					mock.EXPECT().NewReader(gomock.Any()).Return(mockReader).Times(1)
					mockReader.EXPECT().Read().Return(nil, errors.New("err")).Times(1)
					return mock
				}(),
			},
		},
		{
			name: "test5_fail_len",
			args: args{
				ctx:  context.Background(),
				tags: []string{"a", "b"},
				path: "a.csv",
			},
			wantErr: true,
			fields: fields{
				fileReader: func() fReaderIface {
					mock := NewMockfReaderIface(mockCtrl)
					mock.EXPECT().Open("a.csv").Return(&os.File{}, nil).Times(1)
					return mock
				}(),
				csvHandler: func() csvHandlerIface {
					mock := NewMockcsvHandlerIface(mockCtrl)
					mockReader := NewMockcsvReaderIface(mockCtrl)
					mock.EXPECT().NewReader(gomock.Any()).Return(mockReader).Times(1)
					mockReader.EXPECT().Read().Return([]string{
						"1", "true", "1000",
					}, nil).Times(1)
					return mock
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				fileReader: tt.fields.fileReader,
				csvHandler: tt.fields.csvHandler,
			}
			gotData, err := s.searchFromCSV(tt.args.ctx, tt.args.tags, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.searchFromCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("storage.searchFromCSV() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
