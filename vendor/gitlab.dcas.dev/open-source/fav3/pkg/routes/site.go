/*
 *   This Source Code Form is subject to the terms of the Mozilla Public
 *   License, v. 2.0. If a copy of the MPL was not distributed with this
 *   file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package routes

import (
	"errors"
	"fmt"
	"github.com/djcass44/go-tracer/tracer"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"gitlab.dcas.dev/open-source/fav3/pkg/network"
	"gitlab.dcas.dev/open-source/fav3/pkg/util"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	MetricParseSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_query_parse_success_total",
		Help: "Total number of successfully parsed requests",
	})
	MetricParseErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_query_parse_errs_total",
		Help: "Total number of unsuccessfully parsed requests",
	})
	MetricCacheHit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_query_cache_hit_total",
		Help: "Total number of cache hits",
	})
	MetricCacheMiss = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_query_cache_miss_total",
		Help: "Total number of cache misses",
	})
)

type IconRoute struct {
	dl     *network.DirectLoader
	sl     *network.ScrapeLoader
	c      *cache.Cache
	client *http.Client
}

func NewIconRoute(secure bool) *IconRoute {
	ir := new(IconRoute)
	ir.dl = network.NewDirectLoader(secure)
	ir.sl = network.NewScrapeLoader(secure)
	ir.c = cache.New(time.Hour*12, time.Hour*2)
	ir.client = util.NewHTTPClient(secure)

	return ir
}

func (ir *IconRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := tracer.GetRequestId(r)
	targetURL := r.URL.Query().Get("site")
	target, err := ir.getTargetURL(targetURL)
	// throw 400 if we get an error
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"id":     id,
			"target": targetURL,
		}).Error("failed to generate target url")
		http.Error(w, fmt.Sprintf("failed for url (%s): %s", targetURL, err), http.StatusBadRequest)
		return
	}
	// check if there's a cache hit first
	hit, found := ir.c.Get(target.String())
	if found {
		MetricCacheHit.Inc()
		// cache was hit, return that instead
		log.Infof("Found cache hit for key: %s", target)
		entry := hit.(*network.CacheEntry)
		network.ProxyContent(entry, &w)
		return
	}
	MetricCacheMiss.Inc()
	// check the directloader
	value, err := ir.dl.GetIcon(r.Context(), target.String())
	if err == nil {
		res, err := network.GetData(r.Context(), ir.client, value)
		if err == nil {
			ir.c.Set(target.String(), res, cache.DefaultExpiration)
			network.ProxyContent(res, &w)
			return
		}
	}
	log.WithError(err).Warning("DirectLoader failed to find a result")
	// check the network loader
	value, err = ir.sl.GetIcon(r.Context(), target.String())
	if err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	res, err := network.GetData(r.Context(), ir.client, value)
	if err != nil {
		log.WithError(err).Error("ScrapeLoader failed to find a result")
		return
	}
	ir.c.Set(target.String(), res, cache.DefaultExpiration)
	network.ProxyContent(res, &w)
}

func (ir *IconRoute) getTargetURL(s string) (*url.URL, error) {
	if s == "" {
		MetricParseErrors.Inc()
		return nil, errors.New("empty string is not a valid URI")
	}
	// add a scheme if we haven't been given one
	if !strings.HasPrefix(s, "http") {
		s = fmt.Sprintf("https://%s", s)
		log.Warnf("encountered url without a scheme, new url %s", s)
	}

	// parse the url
	u, err := url.Parse(s)
	if err != nil {
		MetricParseErrors.Inc()
		log.WithError(err).Error("failed to parse URL")
		return nil, err
	}
	log.Infof("extracted host from url: %s", u.Host)
	if u.Scheme == "http" {
		MetricParseErrors.Inc()
		log.Warnf("rejecting request for insecure target: %s", s)
		return nil, errors.New("http urls are not accepted")
	}

	// remove parts of the url we dont care about
	u2, err := url.Parse(fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	if err != nil || u2.Host == "" {
		MetricParseErrors.Inc()
		return nil, err
	}
	MetricParseSuccess.Inc()
	log.Infof("extracted minified url: %s", u2.String())
	return u2, nil
}
