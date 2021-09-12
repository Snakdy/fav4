package main

import (
	"gitlab.com/autokubeops/serverless"
	"net/http"
)

func main() {
	serverless.Run(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello " + r.UserAgent()))
	}))
}
