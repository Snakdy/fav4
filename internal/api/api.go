package api

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.dcas.dev/open-source/fav4/pkg/network"
	"net/http"
	"net/url"
	"strings"
)

type IconAPI struct {
	client *http.Client
	direct *network.DirectLoader
}

func NewIconAPI() *IconAPI {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig.MinVersion = tls.VersionTLS13
	api := new(IconAPI)
	api.client = &http.Client{
		Transport: transport,
	}
	api.direct = network.NewDirectLoader(api.client)

	return api
}

func (api *IconAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("%s %s %s %s [%s]", r.Method, r.URL.Path, r.UserAgent(), r.RemoteAddr, r.URL.Query().Get("site"))
	target, err := api.parse(r.URL.Query().Get("site"))
	if target == "" || err != nil {
		http.Error(w, "failed to parse 'site' parameter", http.StatusBadRequest)
		return
	}
	val := api.direct.Get(r.Context(), target)
	if val == "" {
		http.NotFound(w, r)
		return
	}
	_ = network.Download(r.Context(), api.client, val, w)
}

func (*IconAPI) parse(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if !strings.HasPrefix(s, "http") {
		s = fmt.Sprintf("https://%s", s)
	}
	uri, err := url.Parse(s)
	if err != nil {
		log.WithError(err).Error("failed to parse URI")
		return "", err
	}
	return fmt.Sprintf("https://%s", uri.Host), nil
}
