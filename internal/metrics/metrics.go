package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Status = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "http_response_status",
	Help: "Status of HTTP requests peer endpoint",
}, []string{"code", "endpoint"})

func StartServer(port ...string) {
	if len(port) == 0 {
		port[0] = "8081"
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Serving metrics...")
		if err := http.ListenAndServe(":"+port[0], nil); err != nil {
			log.Fatalf("err creating metrics server: %+v\n", err)
		}
	}()
}
