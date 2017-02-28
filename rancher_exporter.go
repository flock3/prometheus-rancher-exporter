package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/infinityworksltd/go-common/logger"
	cfg "github.com/infinityworksltd/prometheus-rancher-exporter/config"
	"github.com/infinityworksltd/prometheus-rancher-exporter/measure"
)

const (
	namespace = "rancher" // Used to prepand Prometheus metrics created by this exporter.
)

// Runtime variables, user controllable for targeting, authentication and filtering.
var (
	config = cfg.Init()
	log    = logger.Start(config)
)

// Predefined variables that are used throughout the exporter
var (
	agentStates   = []string{"activating", "active", "reconnecting", "disconnected", "disconnecting", "finishing-reconnect", "reconnected"}
	hostStates    = []string{"activating", "active", "deactivating", "error", "erroring", "inactive", "provisioned", "purged", "purging", "registering", "removed", "removing", "requested", "restoring", "updating_active", "updating_inactive"}
	stackStates   = []string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "error", "erroring", "finishing_upgrade", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "upgraded", "upgrading"}
	serviceStates = []string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "deactivating", "finishing_upgrade", "inactive", "registering", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "updating_inactive", "upgraded", "upgrading"}
	healthStates  = []string{"healthy", "unhealthy"}
	endpoints     = []string{"stacks", "services", "hosts"} // EndPoints the exporter will trawl
	stackRef      = make(map[string]string)                 // Stores the StackID and StackName as a map, used to provide label dimensions to service metrics

)

func main() {
	flag.Parse()

	// check the rancherURL ($CATTLE_URL) has been provided correctly
	if config.RancherURL() == "" {
		log.Fatal("CATTLE_URL must be set and non-empty")
	}

	log.Info("Starting Prometheus Exporter for Rancher")
	log.Info("Runtime Configuration in-use: URL of Rancher Server: ", config.RancherURL(), " AccessKey: ", config.AccessKey(), "System Services hidden: ", config.HideSys())

	// Register internal metrics used for tracking the exporter performance
	measure.Init()

	// Register a new Exporter
	Exporter := newExporter(config)

	// Register Metrics from each of the endpoints
	// This invokes the Collect method through the prometheus client libraries.
	prometheus.MustRegister(Exporter)

	// Setup HTTP handler
	http.Handle(config.MetricsPath(), prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		                <head><title>Rancher exporter</title></head>
		                <body>
		                   <h1>rancher exporter</h1>
		                   <p><a href='` + config.MetricsPath() + `'>Metrics</a></p>
		                   </body>
		                </html>
		              `))
	})
	log.Printf("Starting Server on port %s and path %s", config.ListenPort(), config.MetricsPath())
	log.Fatal(http.ListenAndServe(config.ListenPort(), nil))
}
