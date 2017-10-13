package collector

import (
	"errors"
	"sync"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	namespace = "openstack_service"
)

var (
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"openstack_service_exporter: Whether a collector succeeded.",
		[]string{"collector"}, nil,
	)
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"openstack_service_exporter: Whether a collector succeeded.",
		[]string{"collector"}, nil,
	)

	factories = make(map[string]func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (Collector, error))
)

// Collector must be implemented by collectors
type Collector interface {
	Describe(chan<- *prometheus.Desc)
	Update(chan<- prometheus.Metric) error
}

func registerCollector(name string, fn func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (Collector, error)) {
	factories[name] = fn
}

// NewCollector creates a new prometheus.Collector which scapes all enabledCollectors
func NewCollector(provider *gophercloud.ProviderClient, opts gophercloud.EndpointOpts, enabledCollectors ...string) (*ServiceCollector, error) {
	collectors := make(map[string]Collector)
	for _, collector := range enabledCollectors {
		factory, ok := factories[collector]
		if !ok {
			log.Errorf("No such collector: %s.", collector)
		} else {
			newCollector, err := factory(provider, opts)
			if err != nil {
				log.Errorf("Unable to use collector %s: %v", collector, err)
				continue
			}
			collectors[collector] = newCollector
		}
	}
	if len(collectors) == 0 {
		return nil, errors.New("no collectors available")
	}
	return &ServiceCollector{collectors: collectors}, nil
}

// ServiceCollector combines all collectors and implements prometheus.Collector
type ServiceCollector struct {
	collectors map[string]Collector
}

// Describe sends metrics descriptions to ch
func (c ServiceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc

	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
}

// Collect all the data from all collectors
func (c ServiceCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.collectors))

	for name, collector := range c.collectors {
		go func(name string, collector Collector) {
			begin := time.Now()
			err := collector.Update(ch)
			duration := time.Since(begin)
			var success float64

			if err != nil {
				log.Errorf("ERROR: %s failed after %fs: %s", name, duration.Seconds(), err)
				success = 0
			} else {
				log.Debugf("OK: %s succeeded after %fs.", name, duration.Seconds())
				success = 1
			}

			ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
			ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)

			wg.Done()
		}(name, collector)
	}

	wg.Wait()
}
