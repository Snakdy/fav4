package network

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func must(s string) *url.URL {
	uri, _ := url.Parse(s)
	return uri
}

func TestHeadLoader_Get(t *testing.T) {
	ctx := logr.NewContext(context.TODO(), testr.NewWithOptions(t, testr.Options{Verbosity: 10}))
	var cases = []struct {
		in  *url.URL
		out string
	}{
		{
			must("https://rancher.com"),
			"https://www.rancher.com/assets/img/favicon.png",
		},
		{
			must("https://rancher.com/docs"),
			"https://rancher.com/docs/img/favicon.png",
		},
		{
			must("https://kubernetes.io/docs"),
			"https://kubernetes.io/icons/favicon-64.png",
		},
		{
			must("https://console.io.dcas.dev"),
			"https://console.io.dcas.dev/static/assets/okd-favicon.png",
		},
	}

	l := &HeadLoader{
		client: http.DefaultClient,
	}
	for _, tt := range cases {
		t.Run(tt.in.String(), func(t *testing.T) {
			path, err := l.Get(ctx, tt.in)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, path)
		})
	}
}
