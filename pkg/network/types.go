package network

import (
	"context"
	"net/url"
)

type Loader interface {
	Get(ctx context.Context, target *url.URL) (string, error)
}
