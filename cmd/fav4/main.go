package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/autokubeops/serverless"
	"net/http"
)

func main() {
	serverless.Run(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s %s", r.Method, r.URL.Path, r.UserAgent(), r.RemoteAddr)
		_, _ = w.Write([]byte("Hello " + r.UserAgent()))
	}))
}
