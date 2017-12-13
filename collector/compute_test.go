package collector

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/prometheus/client_golang/prometheus"
)

func TestComputeUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	handleComputeService(t)

	tests := []struct {
		desc string
		want prometheus.Metric
	}{
		{
			desc: "1. Up (of up and enabled)",
			want: prometheus.MustNewConstMetric(computeUpDesc, prometheus.GaugeValue, 1, "10", "nova-compute", "DE-IX-001-02-02-09-2", "ix1"),
		},
		{
			desc: "1. Enabled (of up and enabled)",
			want: prometheus.MustNewConstMetric(computeEnabledDesc, prometheus.GaugeValue, 1, "10", "nova-compute", "DE-IX-001-02-02-09-2", "ix1"),
		},
		{
			desc: "1. Last Seen (of up and enabled)",
			want: prometheus.MustNewConstMetric(computeLastSeenDesc, prometheus.CounterValue, 1507554374, "10", "nova-compute", "DE-IX-001-02-02-09-2", "ix1"),
		},
		{
			desc: "2. Down (of down and enabled)",
			want: prometheus.MustNewConstMetric(computeUpDesc, prometheus.GaugeValue, 0, "185", "nova-conductor", "DE-IX-001-02-02-01-1", "internal"),
		},
		{
			desc: "2. Enabled (of down and enabled)",
			want: prometheus.MustNewConstMetric(computeEnabledDesc, prometheus.GaugeValue, 1, "185", "nova-conductor", "DE-IX-001-02-02-01-1", "internal"),
		},
		{
			desc: "2. Last Seen (of down and enabled)",
			want: prometheus.MustNewConstMetric(computeLastSeenDesc, prometheus.CounterValue, 1507554376, "185", "nova-conductor", "DE-IX-001-02-02-01-1", "internal"),
		},
		{
			desc: "3. Down (of down and disabled)",

			want: prometheus.MustNewConstMetric(computeUpDesc, prometheus.GaugeValue, 0, "210", "nova-consoleauth", "DE-ES-001-03-09-01-3", "internal"),
		},
		{
			desc: "3. Disabled (of down and disabled)",
			want: prometheus.MustNewConstMetric(computeEnabledDesc, prometheus.GaugeValue, 0, "210", "nova-consoleauth", "DE-ES-001-03-09-01-3", "internal"),
		},
		{
			desc: "3. Last Seen (of down and disabled)",
			want: prometheus.MustNewConstMetric(computeLastSeenDesc, prometheus.CounterValue, 1507554377, "210", "nova-consoleauth", "DE-ES-001-03-09-01-3", "internal"),
		},
	}

	collector := computeCollector{
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
				t.Errorf("%v", err)
			}
		})
	}

	// drain the channel
	res := 0
	for _ = range ch {
		res++
	}
	if res > 0 {
		t.Errorf("Got %d unexpected metrics.", res)
	}
}

func handleComputeService(t *testing.T) {
	th.Mux.HandleFunc("/os-services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, computeServiceBody)
	})
}

const computeServiceBody = `
{
  "services": [
    {
      "status": "enabled",
      "binary": "nova-compute",
      "host": "DE-IX-001-02-02-09-2",
      "zone": "ix1",
      "state": "up",
      "disabled_reason": null,
      "id": 10,
      "updated_at": "2017-10-09T13:06:14.000000"
    },
    {
      "status": "enabled",
      "binary": "nova-conductor",
      "host": "DE-IX-001-02-02-01-1",
      "zone": "internal",
      "state": "down",
      "disabled_reason": null,
      "id": 185,
      "updated_at": "2017-10-09T13:06:16.000000"
    },
    {
      "status": "disabled",
      "binary": "nova-consoleauth",
      "host": "DE-ES-001-03-09-01-3",
      "zone": "internal",
      "state": "down",
      "disabled_reason": null,
      "id": 210,
      "updated_at": "2017-10-09T13:06:17.000000"
    }
  ]
}
`
