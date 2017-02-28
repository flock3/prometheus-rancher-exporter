package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	rancherURL string
	accessKey  string
	secretKey  string
	hideSys    bool
	mutex      sync.RWMutex
	gaugeVecs  map[string]*prometheus.GaugeVec
}

// newExporter creates the metrics we wish to monitor
func newExporter(config config.Config) *Exporter {

	gaugeVecs := addMetrics()
	return &Exporter{
		gaugeVecs:  gaugeVecs,
		rancherURL: config.RancherURL(),
		accessKey:  config.AccessKey(),
		secretKey:  config.SecretKey(),
		hideSys:    config.HideSys(),
	}
}
