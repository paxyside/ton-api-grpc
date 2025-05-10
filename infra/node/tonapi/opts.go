package tonapi

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type rateLimiterRoundTripper struct {
	limiter  *rate.Limiter
	delegate http.RoundTripper
}

func (r *rateLimiterRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	err := r.limiter.Wait(req.Context())
	if err != nil {
		return nil, err
	}
	return r.delegate.RoundTrip(req)
}

type NodeOptFn func(opts *NodeOpts) error

type NodeOpts struct {
	Client *http.Client
}

func WithHTTPClient(client *http.Client) NodeOptFn {
	return func(cfg *NodeOpts) error {
		cfg.Client = client
		return nil
	}
}

func WithTimeout(timeout time.Duration) NodeOptFn {
	return func(opts *NodeOpts) error {
		if opts.Client == nil {
			opts.Client = &http.Client{}
		}

		opts.Client.Timeout = timeout
		return nil
	}
}
func WithRateLimit(rps, burst int) NodeOptFn {
	return func(opts *NodeOpts) error {
		baseTransport := http.DefaultTransport
		if opts.Client.Transport != nil {
			baseTransport = opts.Client.Transport
		}

		opts.Client.Transport = &rateLimiterRoundTripper{
			limiter:  rate.NewLimiter(rate.Limit(rps), burst),
			delegate: baseTransport,
		}

		return nil
	}
}
