package network

import (
	"context"
	"github.com/djcass44/go-utils/pkg/httputils"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func Download(ctx context.Context, client *http.Client, target string, w http.ResponseWriter) error {
	fields := log.Fields{"site": target}
	log.WithFields(fields).Info("downloading content")
	ctx, _ = context.WithTimeout(ctx, time.Second*5)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		log.WithError(err).WithFields(fields).Error("failed to prepare request")
		return err
	}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).WithFields(fields).Error("failed to execute request")
		return err
	}
	defer resp.Body.Close()
	log.WithFields(fields).Debugf("upstream responded with %d in %s", resp.StatusCode, time.Since(start))
	w.WriteHeader(resp.StatusCode)
	if httputils.IsHTTPError(resp.StatusCode) {
		_, _ = w.Write([]byte("failed to download icon"))
		return nil
	}
	_, err = io.Copy(w, resp.Body)
	return err
}
