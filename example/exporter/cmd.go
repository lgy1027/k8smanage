package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"relaper.com/kubemanage/example/exporter/collector"
	_ "relaper.com/kubemanage/example/exporter/kube"
)

var (
// Set during go build
// version   string
// gitCommit string

// 命令行参数
//listenAddr       = flag.String("web.listen-port", "9998", "An port to listen on for web interface and telemetry.")
//metricsPath      = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
//metricsNamespace = flag.String("metric.namespace", "demo", "Prometheus metrics namespace, as the prefix of metrics name")
)

func main() {
	flag.Parse()

	metrics := collector.NewMetrics()
	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>A Prometheus Exporter</title></head>
			<body>
			<h1>A Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Starting Server at http://localhost:%s%s", "9998", "/metrics")
	log.Fatal(http.ListenAndServe(":9998", nil))
}
