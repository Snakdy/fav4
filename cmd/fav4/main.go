package main

import (
	"github.com/djcass44/go-tracer/tracer"
	"github.com/rs/cors"
	"gitlab.com/autokubeops/serverless"
	"gitlab.dcas.dev/open-source/fav4/internal/api"
)

func main() {
	route := api.NewIconAPI()
	serverless.NewBuilder(tracer.NewHandler(cors.AllowAll().Handler(route))).
		WithPrometheus().
		Run()
}
