/*
 *   This Source Code Form is subject to the terms of the Mozilla Public
 *   License, v. 2.0. If a copy of the MPL was not distributed with this
 *   file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"
	"fmt"
	"github.com/djcass44/go-tracer/tracer"
	"github.com/djcass44/go-utils/pkg/httputils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	MetricDownloadSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_download_success_total",
		Help: "Total number of successful favicon downloads",
	})
	MetricDownloadErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_download_errs_total",
		Help: "Total number of failed favicon downloads",
	})
	MetricDownloadBytes = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_download_bytes_total",
		Help: "Total number of downloaded bytes",
	})
)

func GetData(ctx context.Context, client *http.Client, target string) (*CacheEntry, error) {
	id := tracer.GetContextId(ctx)
	log.WithFields(log.Fields{
		"id":     id,
		"target": target,
	}).Info("downloading content")
	response, err := client.Get(target)
	if err != nil {
		MetricDownloadErrors.Inc()
		log.WithError(err).WithFields(log.Fields{
			"id":     id,
			"target": target,
		}).Error("failed to execute download request")
		return nil, err
	}
	// no point going further if we didn't get a 200
	if httputils.IsHTTPError(response.StatusCode) {
		MetricDownloadErrors.Inc()
		log.WithFields(log.Fields{
			"id":     id,
			"target": target,
			"code":   response.StatusCode,
		}).Error("received invalid status code in response")
		return nil, fmt.Errorf("got invalid status code for request: %d >= 400 or < 200", response.StatusCode)
	}
	// close the body
	defer response.Body.Close()
	// read data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		MetricDownloadErrors.Inc()
		return nil, err
	}
	MetricDownloadBytes.Add(float64(len(body)))
	MetricDownloadSuccess.Inc()
	return &CacheEntry{
		Data:        string(body),
		ContentType: response.Header.Get("Content-Type"),
	}, nil
}

func ProxyContent(entry *CacheEntry, w *http.ResponseWriter) {
	writer := *w
	// write the contentType header
	writer.Header().Set("Content-Type", entry.ContentType)
	writer.Header().Set("Cache-Control", "max-age=604800")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(entry.Data))
}
