package main

import "net/http"

func Run(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
