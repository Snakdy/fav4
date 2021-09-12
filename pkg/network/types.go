package network

import "context"

type Loader interface {
	Get(ctx context.Context, target string) string
}
