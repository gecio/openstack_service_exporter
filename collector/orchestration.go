package collector

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/innovocloud/gophercloud_extensions/openstack/orchestration/v1/services"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	orchestrationUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "orchestration", "up"),
		"Status of orchestration services",
		[]string{"id", "binary", "service_host", "engine_id", "topic"}, nil,
	)
	orchestrationLastSeenDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "orchestration", "last_seen"),
		"Last time the service was seen by OpenStack",
		[]string{"id", "binary", "service_host", "engine_id", "topic"}, nil,
	)
)

func init() {
	registerCollector("orchestration", newOrchestrationCollector)
}

func newOrchestrationCollector(provider *gophercloud.ProviderClient, opts gophercloud.EndpointOpts) (Collector, error) {
	client, err := openstack.NewOrchestrationV1(provider, opts)
	if err != nil {
		return nil, err
	}
	return orchestrationCollector{client}, nil
}

type orchestrationCollector struct {
	client *gophercloud.ServiceClient
}

func (c orchestrationCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- orchestrationUpDesc
	ch <- orchestrationLastSeenDesc
}

func (c orchestrationCollector) Update(ch chan<- prometheus.Metric) error {
	pager := services.List(c.client)
	if pager.Err != nil {
		return fmt.Errorf("Unable to get data: %v", pager.Err)
	}

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		services, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}

		for _, service := range services {
			var state float64
			if service.Status == "up" {
				state = 1
			}
			ch <- prometheus.MustNewConstMetric(orchestrationUpDesc, prometheus.GaugeValue, state, service.ID, service.Binary, service.Host, service.EngineID, service.Topic)
			ch <- prometheus.MustNewConstMetric(orchestrationLastSeenDesc, prometheus.CounterValue, float64(service.Updated.Unix()), service.ID, service.Binary, service.Host, service.EngineID, service.Topic)
		}

		return true, nil
	})
	if err != nil {
		return fmt.Errorf("During fetching services: %v", err)
	}

	return nil
}
