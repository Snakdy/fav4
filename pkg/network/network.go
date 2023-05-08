package network

import (
	"context"
	"github.com/djcass44/go-utils/pkg/httputils"
	"github.com/go-logr/logr"
	"io"
	"net/http"
	"time"
)

func Download(ctx context.Context, client *http.Client, target string, w http.ResponseWriter) error {
	log := logr.FromContextOrDiscard(ctx).WithValues("url", target)
	log.Info("downloading content")
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	// send the request off
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		log.Error(err, "failed to prepare request")
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "failed to execute request")
		return err
	}
	defer resp.Body.Close()
	log.V(1).Info("upstream responded", "code", resp.StatusCode, "duration", time.Since(start))
	// handle the response
	w.WriteHeader(resp.StatusCode)
	if httputils.IsHTTPError(resp.StatusCode) {
		_, _ = w.Write([]byte("failed to download icon"))
		return nil
	}
	// copy the response
	_, err = io.Copy(w, resp.Body)
	return err
}
