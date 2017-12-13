package collector

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/innovocloud/gophercloud_extensions/openstack/networking/v2/agents"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	networkUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "network", "up"),
		"State of network agents",
		[]string{"id", "binary", "service_host", "zone", "agent_type", "topic"}, nil,
	)
	networkEnabledDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "network", "enabled"),
		"State of network agents",
		[]string{"id", "binary", "service_host", "zone", "agent_type", "topic"}, nil,
	)
	networkLastSeenDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "network", "last_seen"),
		"Last time the service was seen by OpenStack",
		[]string{"id", "binary", "service_host", "zone", "agent_type", "topic"}, nil,
	)
)

func init() {
	registerCollector("network", newNetworkCollector)
}

func newNetworkCollector(provider *gophercloud.ProviderClient, opts gophercloud.EndpointOpts) (Collector, error) {
	client, err := openstack.NewNetworkV2(provider, opts)
	if err != nil {
		return nil, err
	}
	return networkCollector{client}, nil
}

type networkCollector struct {
	client *gophercloud.ServiceClient
}

func (c networkCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- networkUpDesc
}

func (c networkCollector) Update(ch chan<- prometheus.Metric) error {
	pager := agents.List(c.client, nil)
	if pager.Err != nil {
		return fmt.Errorf("Unable to get data: %v", pager.Err)
	}

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		agents, err := agents.ExtractAgents(page)
		if err != nil {
			return false, err
		}

		for _, agent := range agents {
			var state float64
			var enabled float64
			if agent.Alive {
				state = 1
			}
			// Disabled by admin. So it is neither up (1) nor down (0).
			if agent.AdminStateUp {
				enabled = 1
			}

			ch <- prometheus.MustNewConstMetric(networkUpDesc, prometheus.GaugeValue, state, agent.ID, agent.Binary, agent.Host, agent.AvailabilityZone, agent.AgentType, agent.Topic)
			ch <- prometheus.MustNewConstMetric(networkEnabledDesc, prometheus.GaugeValue, enabled, agent.ID, agent.Binary, agent.Host, agent.AvailabilityZone, agent.AgentType, agent.Topic)
			ch <- prometheus.MustNewConstMetric(networkLastSeenDesc, prometheus.CounterValue, float64(agent.Heartbeat.Unix()), agent.ID, agent.Binary, agent.Host, agent.AvailabilityZone, agent.AgentType, agent.Topic)
		}

		return true, nil
	})
	if err != nil {
		return fmt.Errorf("During fetching services: %v", err)
	}

	return nil
}
