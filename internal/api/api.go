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
	client  *http.Client
	loaders []network.Loader
}

func NewIconAPI() *IconAPI {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig.MinVersion = tls.VersionTLS13
	transport.ForceAttemptHTTP2 = true
	api := new(IconAPI)
	api.client = &http.Client{
		Transport: transport,
	}
	api.loaders = []network.Loader{
		network.NewDirectLoader(api.client),
		network.NewHeadLoader(api.client),
	}

	return api
}

func (api *IconAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("%s %s %s %s [%s]", r.Method, r.URL.Path, r.UserAgent(), r.RemoteAddr, r.URL.Query().Get("site"))
	target, err := api.parse(r.URL.Query().Get("site"))
	if target == nil || err != nil {
		http.Error(w, "failed to parse 'site' parameter", http.StatusBadRequest)
		return
	}
	var val string
	for _, l := range api.loaders {
		val, err = l.Get(r.Context(), target)
		if err != nil {
			continue
		}
		break
	}
	if val == "" {
		http.NotFound(w, r)
		return
	}
	_ = network.Download(r.Context(), api.client, val, w)
}

func (*IconAPI) parse(s string) (*url.URL, error) {
	if s == "" {
		return nil, nil
	}
	if !strings.HasPrefix(s, "http") {
		s = fmt.Sprintf("https://%s", s)
	}
	uri, err := url.Parse(s)
	if err != nil {
		log.WithError(err).Error("failed to parse URI")
		return nil, err
	}
	// only https
	if uri.Scheme == "http" {
		uri.Scheme = "https"
	}
	return uri, nil
}
