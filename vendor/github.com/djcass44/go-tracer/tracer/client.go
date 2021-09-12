/*
 *    Copyright 2020 Django Cass
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package tracer

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

var DefaultRequestHeader = "X-Request-ID"

type traceableRoundTripper struct {
	Proxied http.RoundTripper
}

func (rt *traceableRoundTripper) RoundTrip(r *http.Request) (res *http.Response, err error) {
	// inject the request id header
	id := GetRequestId(r)
	if id != "" {
		r.Header.Add(DefaultRequestHeader, id)
	}
	res, err = rt.Proxied.RoundTrip(r)
	return
}

// Apply modifies an http client's transport to inject the X-Request-ID header
func Apply(c *http.Client) {
	log.Debug("injecting request tracing into http client")
	var transport http.RoundTripper
	if c.Transport == nil {
		transport = http.DefaultTransport
	} else {
		transport = c.Transport
	}
	rt := traceableRoundTripper{transport}
	c.Transport = &rt
}
