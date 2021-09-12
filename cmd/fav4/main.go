package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/autokubeops/serverless"
	"gitlab.dcas.dev/open-source/fav3/pkg/routes"
	"net/http"
)

func main() {
	route := routes.NewIconRoute(true)
	serverless.Run(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s %s [%s]", r.Method, r.URL.Path, r.UserAgent(), r.RemoteAddr, r.URL.Query().Get("site"))
		route.ServeHTTP(w, r)
	}))
}
