package collector

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/innovocloud/gophercloud_extensions/openstack/blockstorage/v2/services"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	blockstorageUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "blockstorage", "up"),
		"Status of blockstorage services",
		[]string{"binary", "service_host", "zone"}, nil,
	)
	blockstorageEnabledDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "blockstorage", "enabled"),
		"Admin status of blockstorage services",
		[]string{"binary", "service_host", "zone"}, nil,
	)
)

func init() {
	registerCollector("blockstorage", newBlockstorageCollector)
}

func newBlockstorageCollector(provider *gophercloud.ProviderClient, opts gophercloud.EndpointOpts) (Collector, error) {
	client, err := openstack.NewBlockStorageV2(provider, opts)
	if err != nil {
		return nil, err
	}
	return blockstorageCollector{client}, nil
}

type blockstorageCollector struct {
	client *gophercloud.ServiceClient
}

func (c blockstorageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- blockstorageUpDesc
	ch <- blockstorageEnabledDesc
}

func (c blockstorageCollector) Update(ch chan<- prometheus.Metric) error {
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
			var enabled float64
			if service.State == "up" {
				state = 1
			}
			// Disabled by admin. So it is neither up (1) nor down (0).
			if service.Status == "enabled" {
				enabled = 1
			}
			ch <- prometheus.MustNewConstMetric(blockstorageUpDesc, prometheus.GaugeValue, state, service.Binary, service.Host, service.Zone)
			ch <- prometheus.MustNewConstMetric(blockstorageEnabledDesc, prometheus.GaugeValue, enabled, service.Binary, service.Host, service.Zone)
		}

		return true, nil
	})
	if err != nil {
		return fmt.Errorf("During fetching services: %v", err)
	}

	return nil
}
