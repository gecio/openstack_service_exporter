package collector

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/prometheus/client_golang/prometheus"
)

func TestBlockstorageUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	handleBlockstorageService(t)

	tests := []struct {
		desc string
		want prometheus.Metric
	}{
		{
			desc: "1. Up (of up and enabled)",
			want: prometheus.MustNewConstMetric(blockstorageUpDesc, prometheus.GaugeValue, 1, "cinder-scheduler", "DE-IX-001-02-02-01-1", "nova"),
		},
		{
			desc: "1. Enabled (of up and enabled)",
			want: prometheus.MustNewConstMetric(blockstorageEnabledDesc, prometheus.GaugeValue, 1, "cinder-scheduler", "DE-IX-001-02-02-01-1", "nova"),
		},
		{
			desc: "2. Down (of down and enabled)",
			want: prometheus.MustNewConstMetric(blockstorageUpDesc, prometheus.GaugeValue, 0, "cinder-scheduler", "DE-IX-001-02-02-01-2", "nova"),
		},
		{
			desc: "2. Enabled (of down and enabled)",
			want: prometheus.MustNewConstMetric(blockstorageEnabledDesc, prometheus.GaugeValue, 1, "cinder-scheduler", "DE-IX-001-02-02-01-2", "nova"),
		},
		{
			desc: "3. Up (of up and disabled)",
			want: prometheus.MustNewConstMetric(blockstorageUpDesc, prometheus.GaugeValue, 1, "cinder-scheduler", "DE-ES-001-02-02-01-3", "nova"),
		},
		{
			desc: "3. Disabled (of up and disabled)",
			want: prometheus.MustNewConstMetric(blockstorageEnabledDesc, prometheus.GaugeValue, 0, "cinder-scheduler", "DE-ES-001-02-02-01-3", "nova"),
		},
		{
			desc: "4. Down (of down and disabled)",
			want: prometheus.MustNewConstMetric(blockstorageUpDesc, prometheus.GaugeValue, 0, "cinder-scheduler", "DE-ES-001-02-02-01-4", "nova"),
		},
		{
			desc: "4. Disabled (of down and disabled)",
			want: prometheus.MustNewConstMetric(blockstorageEnabledDesc, prometheus.GaugeValue, 0, "cinder-scheduler", "DE-ES-001-02-02-01-4", "nova"),
		},
	}

	collector := blockstorageCollector{
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

func handleBlockstorageService(t *testing.T) {
	th.Mux.HandleFunc("/os-services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, blockstorageServiceBody)
	})
}

const blockstorageServiceBody = `
{
  "services": [
    {
      "status": "enabled",
      "binary": "cinder-scheduler",
      "zone": "nova",
      "state": "up",
      "updated_at": "2017-10-09T13:53:33.000000",
      "host": "DE-IX-001-02-02-01-1",
      "disabled_reason": null
    },
    {
      "status": "enabled",
      "binary": "cinder-scheduler",
      "zone": "nova",
      "state": "down",
      "updated_at": "2017-10-09T13:53:33.000000",
      "host": "DE-IX-001-02-02-01-2",
      "disabled_reason": null
    },
    {
      "status": "disabled",
      "binary": "cinder-scheduler",
      "zone": "nova",
      "state": "up",
      "updated_at": "2017-10-09T13:53:28.000000",
      "host": "DE-ES-001-02-02-01-3",
      "disabled_reason": null
    },
    {
      "status": "disabled",
      "binary": "cinder-scheduler",
      "zone": "nova",
      "state": "down",
      "updated_at": "2017-10-09T13:53:28.000000",
      "host": "DE-ES-001-02-02-01-4",
      "disabled_reason": null
    }
  ]
}
`
