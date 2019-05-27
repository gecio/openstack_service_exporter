package collector

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNetworkUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	handleNetworkAgent(t)

	tests := []struct {
		desc string
		want prometheus.Metric
	}{
		{
			desc: "1. Up (of up and enabled)",
			want: prometheus.MustNewConstMetric(networkUpDesc, prometheus.GaugeValue, 1, "02306fa8-d790-4272-8e2f-c3a52c8e9f15", "neutron-openvswitch-agent", "DE-IX-001-02-02-13-1", "", "Open vSwitch agent", "N/A"),
		},
		{
			desc: "1. Enabled (of up and enabled)",
			want: prometheus.MustNewConstMetric(networkEnabledDesc, prometheus.GaugeValue, 1, "02306fa8-d790-4272-8e2f-c3a52c8e9f15", "neutron-openvswitch-agent", "DE-IX-001-02-02-13-1", "", "Open vSwitch agent", "N/A"),
		},
		{
			desc: "1. Last Seen (of up and enabled)",
			want: prometheus.MustNewConstMetric(networkLastSeenDesc, prometheus.CounterValue, 1507555493, "02306fa8-d790-4272-8e2f-c3a52c8e9f15", "neutron-openvswitch-agent", "DE-IX-001-02-02-13-1", "", "Open vSwitch agent", "N/A"),
		},
		{
			desc: "2. Down (of down and enabled)",
			want: prometheus.MustNewConstMetric(networkUpDesc, prometheus.GaugeValue, 0, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-metadata-agent", "DE-IX-001-02-02-13-2", "", "Metadata agent", "N/A"),
		},
		{
			desc: "2. Enabled (of down and enabled)",
			want: prometheus.MustNewConstMetric(networkEnabledDesc, prometheus.GaugeValue, 1, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-metadata-agent", "DE-IX-001-02-02-13-2", "", "Metadata agent", "N/A"),
		},
		{
			desc: "2. Last Seen (of down and enabled)",
			want: prometheus.MustNewConstMetric(networkLastSeenDesc, prometheus.CounterValue, 1507555512, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-metadata-agent", "DE-IX-001-02-02-13-2", "", "Metadata agent", "N/A"),
		},
		{
			desc: "3. Up (of up and disabled)",
			want: prometheus.MustNewConstMetric(networkUpDesc, prometheus.GaugeValue, 1, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-l3-agent", "DE-IX-001-02-02-13-3", "", "L3 agent", "N/A"),
		},
		{
			desc: "3. Disabled (of up and disabled)",
			want: prometheus.MustNewConstMetric(networkEnabledDesc, prometheus.GaugeValue, 0, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-l3-agent", "DE-IX-001-02-02-13-3", "", "L3 agent", "N/A"),
		},
		{
			desc: "3. Last Seen (of up and disabled)",
			want: prometheus.MustNewConstMetric(networkLastSeenDesc, prometheus.CounterValue, 1507555512, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-l3-agent", "DE-IX-001-02-02-13-3", "", "L3 agent", "N/A"),
		},
		{
			desc: "4. Down (of down and disabled)",
			want: prometheus.MustNewConstMetric(networkUpDesc, prometheus.GaugeValue, 0, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-dhcp-agent", "DE-IX-001-02-02-13-4", "", "DHCP agent", "N/A"),
		},
		{
			desc: "4. Disabled (of down and disabled)",
			want: prometheus.MustNewConstMetric(networkEnabledDesc, prometheus.GaugeValue, 0, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-dhcp-agent", "DE-IX-001-02-02-13-4", "", "DHCP agent", "N/A"),
		},
		{
			desc: "4. Last Seen (of down and disabled)",
			want: prometheus.MustNewConstMetric(networkLastSeenDesc, prometheus.CounterValue, 1507555512, "0b4d3a07-b680-40da-88a3-8de99fd34100", "neutron-dhcp-agent", "DE-IX-001-02-02-13-4", "", "DHCP agent", "N/A"),
		},
	}

	collector := networkCollector{
		client: client.ServiceClient(),
	}

	ch := make(chan prometheus.Metric)
	go func() {
		err := collector.Update(ch)
		if err != nil {
			t.Errorf("Update failed: %v", err)
		}
		close(ch)
	}()

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", idx, tt.desc), func(t *testing.T) {
			err := verifyMetric(ch, tt.want)
			if err != nil {
				t.Errorf("%s: %v", tt.desc, err)
			}
		})
	}

	// drain the channel
	res := 0
	for range ch {
		res++
	}
	if res > 0 {
		t.Errorf("Got %d unexpected metrics.", res)
	}
}

func handleNetworkAgent(t *testing.T) {
	th.Mux.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, networkAgentBody)
	})
}

const networkAgentBody = `
{
  "agents": [
    {
      "binary": "neutron-openvswitch-agent",
      "description": null,
      "availability_zone": null,
      "heartbeat_timestamp": "2017-10-09 13:24:53",
      "admin_state_up": true,
      "alive": true,
      "id": "02306fa8-d790-4272-8e2f-c3a52c8e9f15",
      "topic": "N/A",
      "host": "DE-IX-001-02-02-13-1",
      "agent_type": "Open vSwitch agent",
      "started_at": "2017-10-06 11:10:44",
      "created_at": "2017-08-02 08:55:58"
    },
    {
      "binary": "neutron-metadata-agent",
      "description": null,
      "availability_zone": null,
      "heartbeat_timestamp": "2017-10-09 13:25:12",
      "admin_state_up": true,
      "alive": false,
      "id": "0b4d3a07-b680-40da-88a3-8de99fd34100",
      "topic": "N/A",
      "host": "DE-IX-001-02-02-13-2",
      "agent_type": "Metadata agent",
      "started_at": "2017-10-06 10:55:09",
      "created_at": "2017-08-02 08:57:04"
    },
    {
      "binary": "neutron-l3-agent",
      "description": null,
      "availability_zone": null,
      "heartbeat_timestamp": "2017-10-09 13:25:12",
      "admin_state_up": false,
      "alive": true,
      "id": "0b4d3a07-b680-40da-88a3-8de99fd34100",
      "topic": "N/A",
      "host": "DE-IX-001-02-02-13-3",
      "agent_type": "L3 agent",
      "started_at": "2017-10-06 10:55:09",
      "created_at": "2017-08-02 08:57:04"
    },
    {
      "binary": "neutron-dhcp-agent",
      "description": null,
      "availability_zone": null,
      "heartbeat_timestamp": "2017-10-09 13:25:12",
      "admin_state_up": false,
      "alive": false,
      "id": "0b4d3a07-b680-40da-88a3-8de99fd34100",
      "topic": "N/A",
      "host": "DE-IX-001-02-02-13-4",
      "agent_type": "DHCP agent",
      "started_at": "2017-10-06 10:55:09",
      "created_at": "2017-08-02 08:57:04"
    }
  ]
}
`
