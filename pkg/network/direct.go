package network

import (
	"context"
	"fmt"
	"github.com/djcass44/go-utils/pkg/httputils"
	"github.com/go-logr/logr"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var imageMimes = []string{"png", "ico"}

type DirectLoader struct {
	client *http.Client
}

func NewDirectLoader(client *http.Client) *DirectLoader {
	return &DirectLoader{
		client: client,
	}
}

func (dl *DirectLoader) Get(ctx context.Context, target *url.URL) (string, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("url", target.String())
	log.V(1).Info("directly fetching favicon")
	res := make(chan *string)
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	for i := range imageMimes {
		log.V(2).Info("fetching favicon by mime", "mime", imageMimes[i])
		go dl.fetch(ctx, fmt.Sprintf("https://%s/favicon.%s", target.Host, imageMimes[i]), res)
	}
	for i := 0; i < len(imageMimes); i++ {
		response := <-res
		if response == nil {
			continue
		}
		return *response, nil
	}
	return "", ErrNotFound
}

func (dl *DirectLoader) fetch(ctx context.Context, target string, res chan *string) {
	log := logr.FromContextOrDiscard(ctx).WithValues("url", target)
	log.Info("downloading favicon")
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, target, nil)
	if err != nil {
		log.Error(err, "failed to prepare request")
		res <- nil
		return
	}
	resp, err := dl.client.Do(req)
	if err != nil {
		log.Error(err, "failed to execute request")
		res <- nil
		return
	}
	defer resp.Body.Close()
	if !httputils.IsHTTPSuccess(resp.StatusCode) {
		log.Info("unexpected response code", "code", resp.StatusCode)
		res <- nil
		return
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") || contentType == "application/octet-stream" {
		res <- &target
		return
	}
	log.V(1).Info("unexpected content-type", "content-type", contentType)
	res <- nil
}
