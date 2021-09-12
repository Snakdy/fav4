/*
 *   This Source Code Form is subject to the terms of the Mozilla Public
 *   License, v. 2.0. If a copy of the MPL was not distributed with this
 *   file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package network

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djcass44/go-tracer/tracer"
	"github.com/djcass44/go-utils/pkg/sliceutils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"gitlab.dcas.dev/open-source/fav3/pkg/util"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strings"
)

var definitions = []string{"favicon", "shortcut icon", "apple-touch-icon"}

var (
	MetricScrapeSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_load_scrape_success_total",
		Help: "Total number of successful scraped favicon lookups",
	})
	MetricScrapeErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fav_load_scrape_errs_total",
		Help: "Total number of failed scraped favicon lookups",
	})
)

type ScrapeLoader struct {
	client *http.Client
}

func NewScrapeLoader(secure bool) *ScrapeLoader {
	sl := new(ScrapeLoader)
	sl.client = util.NewHTTPClient(secure)

	return sl
}

func (l *ScrapeLoader) GetIcon(ctx context.Context, target string) (string, error) {
	value, err := l.getIcon(ctx, target)
	if err == nil {
		log.Infof("Got acceptable icon at url: %s", value)
		return value, nil
	}
	return "", errors.New("failed to locate valid image")
}

func (l *ScrapeLoader) getIcon(ctx context.Context, target string) (string, error) {
	id := tracer.GetContextId(ctx)
	log.WithFields(log.Fields{
		"id":     id,
		"target": target,
	}).Infof("targeting host")
	response, err := l.client.Get(target)
	if err != nil {
		MetricScrapeErrors.Inc()
		log.WithError(err).WithFields(log.Fields{
			"id":     id,
			"target": target,
		}).Error("failed to scrape site content")
		return "", err
	}
	// close the body
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		MetricScrapeErrors.Inc()
		log.WithError(err).WithField("id", id).Error("failed to read response body")
		return "", err
	}

	data, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		MetricScrapeErrors.Inc()
		log.WithError(err).WithField("id", id).Error("failed to parse html content")
		return "", err
	}
	match := ""

	log.WithField("id", id).Debug("reading document")
	// read the document
	document := goquery.NewDocumentFromNode(data)
	// find relevant tags e.g. <link rel="shortcut icon" href="/img/favicon.png"/>
	document.Find("head").Find("link[rel]").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		log.WithFields(log.Fields{
			"id":    id,
			"nodes": len(selection.Nodes),
		}).Debug("checking nodes")
		// loop over all matching nodes
		for j := range selection.Nodes {
			rel := selection.Nodes[j]
			log.WithFields(log.Fields{
				"id":  id,
				"rel": rel.Attr,
			}).Debug("checking rel")

			// see if this node matches our requirements
			idx := l.isMatchingAttr(rel.Attr)
			if idx >= 0 {
				m := rel.Attr[idx].Val
				if strings.HasPrefix(m, "https://") {
					match = m
				} else {
					if strings.HasPrefix(m, "/") {
						match = target + m
					} else {
						match = fmt.Sprintf("%s/%s", target, m)
					}
				}
				log.WithField("id", id).Infof("found match: %s", match)
				// bail out of the loop
				return false
			}
		}
		return true
	})
	if match != "" {
		MetricScrapeSuccess.Inc()
		return match, nil
	}
	MetricScrapeErrors.Inc()
	return "", errors.New("failed to locate valid image")
}

// isMatchingAttr checks if this node has the attributes we want
func (*ScrapeLoader) isMatchingAttr(attrs []html.Attribute) int {
	// dont initialise to 0 because 0 is a valid array index
	rel, href := -1, -1
	for i, attr := range attrs {
		// check if we're looking at a 'rel' and its value is that of an icon (see definitions)
		if attr.Key == "rel" && sliceutils.Includes(definitions, attr.Val) {
			rel = i
		}
		// we must have an href
		if attr.Key == "href" {
			href = i
		}
	}
	// only return the index if we've found both a valid-rel and href
	if rel >= 0 && href >= 0 {
		return href
	}
	return -1
}
