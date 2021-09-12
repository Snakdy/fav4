/*
 *   This Source Code Form is subject to the terms of the Mozilla Public
 *   License, v. 2.0. If a copy of the MPL was not distributed with this
 *   file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package util

import (
	"crypto/tls"
	"net/http"
)

// NewHTTPClient create an *http.Client with toggled hardening
func NewHTTPClient(secure bool) *http.Client {
	var tlsVersion uint16 = tls.VersionTLS12
	if secure {
		tlsVersion = tls.VersionTLS13
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				MinVersion: tlsVersion,
			},
		},
	}
}
