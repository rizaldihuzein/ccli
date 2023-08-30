package src

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	errUnexpectedCode = errors.New("unexpected response code")
)

//go:generate mockgen -destination=fetch_mock.go -package=src -source=fetch.go
type (
	apiFetcherIface interface {
		getSampleAPIResourceRedirect(ctx context.Context, link []string) (data []UserData, err error)
	}

	apiFetcher struct {
		// httpClient *http.Client
		httpClient httpIface
	}

	httpIface interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

func newFetcher(client *http.Client) (apiFetcherIface, error) {
	if client == nil {
		return nil, errors.New("missing required params")
	}
	return &apiFetcher{
		httpClient: client,
	}, nil
}

func (f *apiFetcher) getSampleAPIResourceRedirect(ctx context.Context, link []string) (data []UserData, err error) {
	var (
		validLinks = 0
		validResp  = 0
	)
	for _, v := range link {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		validLinks++

		resp, err := f.fetchHTTP(ctx, http.MethodGet, v)
		if err != nil && err != errUnexpectedCode {
			return data, err
		}
		if resp.code != http.StatusOK {
			continue
		}

		validResp++
		err = json.Unmarshal(resp.content, &data)
		if err != nil {
			return data, err
		}
		if err == nil {
			return data, err
		}
	}

	if validLinks == 0 {
		return data, errors.New("all links are invalid")
	}

	if validResp == 0 {
		return data, errors.New("all links are down or gives unexpected response")
	}

	return
}

func (f *apiFetcher) fetchHTTP(ctx context.Context, method, link string) (resp httpResponseGeneral, err error) {
	link, method = strings.TrimSpace(link), strings.TrimSpace(method)
	if link == "" || method == "" {
		return resp, errors.New("missing required params")
	}

	_, err = url.ParseRequestURI(link)
	if err != nil {
		return
	}

	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
	default:
		return resp, errors.New("invalid method")
	}

	req, err := http.NewRequestWithContext(ctx, method, link, nil)
	if err != nil {
		return
	}

	httpResp, err := f.httpClient.Do(req)
	if err != nil {
		return
	}

	resp.code = httpResp.StatusCode

	switch resp.code {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
	case http.StatusServiceUnavailable:
		return
	default:
		return resp, errUnexpectedCode
	}

	defer httpResp.Body.Close()
	resp.content, err = ioutil.ReadAll(httpResp.Body)

	return
}
