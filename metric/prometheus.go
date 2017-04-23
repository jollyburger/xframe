package metric

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func ServePrometheus(l net.Listener) {
	http.Serve(l, nil)
}

func init() {
	http.Handle("/metrics", prometheus.Handler())
}
