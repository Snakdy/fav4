package main

import (
	"github.com/djcass44/go-tracer/tracer"
	"gitlab.com/autokubeops/serverless"
	"gitlab.dcas.dev/open-source/fav4/internal/api"
)

func main() {
	route := api.NewIconAPI()
	serverless.NewBuilder(tracer.NewHandler(route)).
		WithPrometheus().
		Run()
}
