package testutils

import (
	"io"
	"net/http"
	"strings"
)

type StubRT struct {
	Resp       string
	StatusCode int
	Opts       []func(resp *http.Response)
}

func (rt StubRT) RoundTrip(_ *http.Request) (*http.Response, error) {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	resp := &http.Response{
		StatusCode: rt.StatusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(rt.Resp)),
	}

	for _, opt := range rt.Opts {
		opt(resp)
	}

	return resp, nil
}
