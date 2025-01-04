package http

import (
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	Rate   int
	Window time.Duration
}

type RateLimitedClient struct {
	client   *http.Client
	limiters []*rate.Limiter
}

func NewRateLimitedClient(limiters []RateLimit) *RateLimitedClient {
	ls := make([]*rate.Limiter, len(limiters))
	for idx, l := range limiters {
		r := l.Window / time.Duration(l.Rate)
		ls[idx] = rate.NewLimiter(rate.Every(r), l.Rate)
	}

	return &RateLimitedClient{
		client:   &http.Client{},
		limiters: ls,
	}
}

func (c *RateLimitedClient) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	eg := errgroup.Group{}
	for idx := 0; idx < len(c.limiters); idx++ {
		eg.Go(func() error {
			return c.limiters[idx].Wait(ctx)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return c.client.Do(req)
}
