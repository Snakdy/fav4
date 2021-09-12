package network

import (
	"context"
	"fmt"
	"github.com/djcass44/go-utils/pkg/httputils"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func (dl *DirectLoader) Get(ctx context.Context, target string) string {
	res := make(chan *string)
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	for i := range imageMimes {
		go dl.fetch(ctx, fmt.Sprintf("%s/favicon.%s", target, imageMimes[i]), res)
	}
	for i := 0; i < len(imageMimes); i++ {
		response := <-res
		if response == nil {
			continue
		}
		return *response
	}
	return ""
}

func (dl *DirectLoader) fetch(ctx context.Context, target string, res chan *string) {
	log.WithField("site", target).Info("downloading favicon")
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, target, nil)
	if err != nil {
		log.WithError(err).Error("failed to prepare request")
		res <- nil
		return
	}
	resp, err := dl.client.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to execute request")
		res <- nil
		return
	}
	defer resp.Body.Close()
	if !httputils.IsHTTPSuccess(resp.StatusCode) {
		log.Infof("unexpected response code: %d", resp.StatusCode)
		res <- nil
		return
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") || contentType == "application/octet-stream" {
		res <- &target
		return
	}
	log.Debugf("unexpected content-type: %s", contentType)
	res <- nil
}
