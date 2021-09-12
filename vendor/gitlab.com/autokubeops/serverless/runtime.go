package serverless

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	MetricNamespace = "serverless"
	MetricSubsystem = "function"
)

var (
	metricOpsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricNamespace,
		Subsystem: MetricSubsystem,
		Name:      "ops_total",
	}, []string{"method", "path"})
	metricOpsDurationTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricNamespace,
		Subsystem: MetricSubsystem,
		Name:      "ops_duration_total",
	}, []string{"method", "path"})
)

type Builder struct {
	handler http.Handler
	port    int
}

func NewBuilder(handler http.Handler) *Builder {
	return &Builder{
		handler: handler,
		port:    8080,
	}
}

// WithPrometheus enables Prometheus metric collection
// from function invocation.
func (b *Builder) WithPrometheus() *Builder {
	h := b.handler
	b.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// preprocessing
		metricOpsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		start := time.Now()
		// call the handler
		h.ServeHTTP(w, r)
		// postprocessing
		metricOpsDurationTotal.WithLabelValues(r.Method, r.URL.Path).Add(float64(time.Since(start).Milliseconds()))
	})
	return b
}

func (b *Builder) WithPort(port int) *Builder {
	b.port = port
	return b
}

func (b *Builder) Run() {
	router := http.NewServeMux()
	router.Handle("/", b.handler)
	go func() {
		log.Printf("listening on :%d (h2c)", b.port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", b.port), h2c.NewHandler(router, &http2.Server{})))
	}()

	// wait for a signal
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigC
	log.Printf("received SIGTERM/SIGINT (%s), shutting down...", sig)
}
