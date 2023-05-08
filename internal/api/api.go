package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-logr/logr"
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
	ctx := r.Context()
	log := logr.FromContextOrDiscard(ctx)
	log.Info("processing request", "site", r.URL.Query().Get("site"))
	target, err := api.parse(ctx, r.URL.Query().Get("site"))
	if target == nil || err != nil {
		log.V(1).Info("missing site parameter")
		if err != nil {
			log.Error(err, "failed to parse 'site' parameter")
		}
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

func (*IconAPI) parse(ctx context.Context, s string) (*url.URL, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("str", s)
	if s == "" {
		return nil, nil
	}
	if !strings.HasPrefix(s, "http") {
		s = fmt.Sprintf("https://%s", s)
	}
	uri, err := url.Parse(s)
	if err != nil {
		log.Error(err, "failed to parse URI")
		return nil, err
	}
	// only https
	if uri.Scheme == "http" {
		uri.Scheme = "https"
	}
	return uri, nil
}
