package src

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
)

func Test_newFetcher(t *testing.T) {
	mockClient := &http.Client{}

	type args struct {
		client *http.Client
	}
	tests := []struct {
		name    string
		args    args
		want    apiFetcherIface
		wantErr bool
	}{
		{
			name:    "test1_fail",
			args:    args{},
			wantErr: true,
		},
		{
			name: "test2_success",
			args: args{
				client: mockClient,
			},
			wantErr: false,
			want: &apiFetcher{
				httpClient: mockClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newFetcher(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFetcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFetcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiFetcher_fetchHTTP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type fields struct {
		httpClient httpIface
	}
	type args struct {
		ctx    context.Context
		method string
		link   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp httpResponseGeneral
		wantErr  bool
	}{
		{
			name: "test1_missing_required_params",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "test2_missing_required_params",
			args: args{
				ctx:  context.Background(),
				link: "http://localhost:8080",
			},
			wantErr: true,
		},
		{
			name: "test3_invalid_url",
			args: args{
				ctx:    context.Background(),
				link:   "a",
				method: http.MethodGet,
			},
			wantErr: true,
		},
		{
			name: "test4_fail_do_request",
			args: args{
				ctx:    context.Background(),
				link:   "http://localhost:8080",
				method: http.MethodGet,
			},
			wantErr: true,
			fields: fields{
				httpClient: func() httpIface {
					mock := NewMockhttpIface(mockCtrl)
					mock.EXPECT().Do(gomock.Any()).Return(nil, errors.New("err"))
					return mock
				}(),
			},
		},
		{
			name: "test5_service_unavailable",
			args: args{
				ctx:    context.Background(),
				link:   "http://localhost:8080",
				method: http.MethodGet,
			},
			wantErr: false,
			wantResp: httpResponseGeneral{
				code: http.StatusServiceUnavailable,
			},
			fields: fields{
				httpClient: func() httpIface {
					mock := NewMockhttpIface(mockCtrl)
					mock.EXPECT().Do(gomock.Any()).Return(&http.Response{
						StatusCode: http.StatusServiceUnavailable,
					}, nil)
					return mock
				}(),
			},
		},
		{
			name: "test6_success",
			args: args{
				ctx:    context.Background(),
				link:   "http://localhost:8080",
				method: http.MethodGet,
			},
			wantErr: false,
			wantResp: httpResponseGeneral{
				content: []byte("[]"),
				code:    http.StatusOK,
			},
			fields: fields{
				httpClient: func() httpIface {
					httpmock.RegisterResponder("GET", "http://localhost:8080", func(req *http.Request) (*http.Response, error) {
						resp := []UserData{}
						return httpmock.NewJsonResponse(http.StatusOK, resp)
					},
					)
					return &http.Client{}
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &apiFetcher{
				httpClient: tt.fields.httpClient,
			}
			gotResp, err := f.fetchHTTP(tt.args.ctx, tt.args.method, tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiFetcher.fetchHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("apiFetcher.fetchHTTP() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_apiFetcher_getSampleAPIResourceRedirect(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type fields struct {
		httpClient httpIface
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
		mock     func()
	}{
		{
			name: "test1_all_links_are_invalid",
			args: args{
				ctx: context.Background(),
				link: []string{
					"  ", "  ",
				},
			},
			wantErr: true,
		},
		{
			name: "test2_all_links_are_down",
			args: args{
				ctx: context.Background(),
				link: []string{
					"http://localhost:8080", "http://localhost:8081",
				},
			},
			wantErr: true,
			fields: fields{
				httpClient: func() httpIface {
					mock := NewMockhttpIface(mockCtrl)
					mock.EXPECT().Do(gomock.Any()).Return(&http.Response{
						StatusCode: http.StatusServiceUnavailable,
					}, nil).Times(2)
					return mock
				}(),
			},
		},
		{
			name: "test3_success_fetch_first_link",
			args: args{
				ctx: context.Background(),
				link: []string{
					"http://localhost:8080", "http://localhost:8081",
				},
			},
			wantErr: false,
			wantData: []UserData{
				{
					ID:      "12",
					Balance: "100",
					Tags:    []string{"tag"},
				},
			},
			fields: fields{
				httpClient: &http.Client{},
			},
			mock: func() {
				httpmock.RegisterResponder("GET", "http://localhost:8080", func(req *http.Request) (*http.Response, error) {
					resp := []UserData{
						{
							ID:      "12",
							Balance: "100",
							Tags:    []string{"tag"},
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				},
				)
			},
		},
		{
			name: "test4_success_fetch_second_link",
			args: args{
				ctx: context.Background(),
				link: []string{
					"http://localhost:8080", "http://localhost:8081",
				},
			},
			wantErr: false,
			wantData: []UserData{
				{
					ID:      "12",
					Balance: "100",
					Tags:    []string{"tag"},
				},
			},
			fields: fields{
				httpClient: &http.Client{},
			},
			mock: func() {
				httpmock.RegisterResponder("GET", "http://localhost:8080", func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusUnavailableForLegalReasons, nil)
				},
				)
				httpmock.RegisterResponder("GET", "http://localhost:8081", func(req *http.Request) (*http.Response, error) {
					resp := []UserData{
						{
							ID:      "12",
							Balance: "100",
							Tags:    []string{"tag"},
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				},
				)
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			f := &apiFetcher{
				httpClient: tt.fields.httpClient,
			}
			gotData, err := f.getSampleAPIResourceRedirect(tt.args.ctx, tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiFetcher.getSampleAPIResourceRedirect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("apiFetcher.getSampleAPIResourceRedirect() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
