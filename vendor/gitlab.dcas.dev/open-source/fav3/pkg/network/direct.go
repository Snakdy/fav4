/*
 *   This Source Code Form is subject to the terms of the Mozilla Public
 *   License, v. 2.0. If a copy of the MPL was not distributed with this
 *   file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"
	"errors"
	"fmt"
	"github.com/djcass44/go-tracer/tracer"
	"github.com/djcass44/go-utils/pkg/httputils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"gitlab.dcas.dev/open-source/fav3/pkg/util"
	"net/http"
	"strings"
)

// global values
var imagesMimes = []string{"png", "ico"}

var (
	MetricDirectSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_load_direct_success_total",
		Help: "Total number of successful direct favicon lookups",
	})
	MetricDirectErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_load_direct_errs_total",
		Help: "Total number of failed direct favicon lookups",
	})
	MetricDirectEmpty = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_load_direct_empty_total",
		Help: "Total number of direct favicon lookups which returned no results",
	})
)

type DirectLoader struct {
	client *http.Client
}

func NewDirectLoader(secure bool) *DirectLoader {
	dl := new(DirectLoader)
	dl.client = util.NewHTTPClient(secure)

	return dl
}

func (l *DirectLoader) GetIcon(ctx context.Context, target string) (string, error) {
	for i := range imagesMimes {
		value, err := l.getIcon(ctx, fmt.Sprintf("%s/favicon.%s", target, imagesMimes[i]))
		if err != nil {
			log.WithError(err).Error("failed to locate favicon")
		} else {
			return value, nil
		}
	}
	MetricDirectEmpty.Inc()
	return "", errors.New("failed to locate valid image")
}

func (l *DirectLoader) getIcon(ctx context.Context, target string) (string, error) {
	id := tracer.GetContextId(ctx)
	log.WithFields(log.Fields{
		"id":     id,
		"target": target,
	}).Infof("targeting host")
	// do the http request
	response, err := l.client.Head(target)
	if err != nil {
		MetricDirectErrors.Inc()
		log.WithError(err).WithFields(log.Fields{
			"id":     id,
			"target": target,
		}).Errorf("failed to execute request")
		return "", err
	}
	if !httputils.IsHTTPSuccess(response.StatusCode) {
		MetricDirectErrors.Inc()
		log.WithFields(log.Fields{
			"id":   id,
			"code": response.StatusCode,
		}).Errorf("got invalid status code for response")
		return "", fmt.Errorf("got non-200 response code for request: %d (%s)", response.StatusCode, response.Status)
	}
	// get the content type
	contentType := response.Header.Get("Content-Type")
	log.WithFields(log.Fields{
		"id":          id,
		"target":      target,
		"contentType": contentType,
	}).Info("reading response metadata")
	// check we've got a compatible contentType (image or octet stream)
	if strings.HasPrefix(contentType, "image/") || contentType == "application/octet-stream" {
		MetricDirectSuccess.Inc()
		return target, nil
	}
	MetricDirectErrors.Inc()
	return "", fmt.Errorf("illegal content type: %s", contentType)
}
