package main

import (
	"github.com/djcass44/go-tracer/tracer"
	"gitlab.com/autokubeops/serverless"
	"gitlab.dcas.dev/open-source/fav4/internal/api"
	"net/http"
)

func main() {
	route := api.NewIconAPI()
	serverless.NewBuilder(tracer.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route.ServeHTTP(w, r)
	}))).
		WithPrometheus().
		Run()
}
