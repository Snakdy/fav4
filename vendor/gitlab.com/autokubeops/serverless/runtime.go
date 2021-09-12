package serverless

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(handler http.Handler) {
	router := http.NewServeMux()
	router.Handle("/", handler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	// wait for a signal
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigC
	log.Printf("received SIGTERM/SIGINT (%s), shutting down...", sig)
}
