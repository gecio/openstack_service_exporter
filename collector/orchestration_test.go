package collector

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/prometheus/client_golang/prometheus"
)

func TestOrchestrationUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	handleOrchestrationService(t)

	tests := []struct {
		desc string
		want prometheus.Metric
	}{
		{
			desc: "heat-engine down",
			want: prometheus.MustNewConstMetric(orchestrationUpDesc, prometheus.GaugeValue, 0, "06529f29-a258-4eb1-8966-61141c3f0c68", "heat-engine", "DE-IX-001-02-02-01-2", "1ab44a4a-c62b-476a-96a8-a3e4525eb1c0", "engine"),
		},
		{
			desc: "Last seen of heat-engine down",
			want: prometheus.MustNewConstMetric(orchestrationLastSeenDesc, prometheus.CounterValue, 1507287003, "06529f29-a258-4eb1-8966-61141c3f0c68", "heat-engine", "DE-IX-001-02-02-01-2", "1ab44a4a-c62b-476a-96a8-a3e4525eb1c0", "engine"),
		},
		{
			desc: "heat-engine up",
			want: prometheus.MustNewConstMetric(orchestrationUpDesc, prometheus.GaugeValue, 1, "07040a0c-bb75-437f-a59c-7ff9c92ef833", "heat-engine", "DE-ES-001-03-09-01-1", "3ed4082d-fdc2-442c-b920-627566993626", "engine"),
		},
		{
			desc: "Last seen of heat-engine up",
			want: prometheus.MustNewConstMetric(orchestrationLastSeenDesc, prometheus.CounterValue, 1507560592, "07040a0c-bb75-437f-a59c-7ff9c92ef833", "heat-engine", "DE-ES-001-03-09-01-1", "3ed4082d-fdc2-442c-b920-627566993626", "engine"),
		},
	}

	collector := orchestrationCollector{
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

func handleOrchestrationService(t *testing.T) {
	th.Mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, orchestrationServiceBody)
	})
}

const orchestrationServiceBody = `
{
  "services": [
    {
      "status": "down",
      "binary": "heat-engine",
      "report_interval": 60,
      "engine_id": "1ab44a4a-c62b-476a-96a8-a3e4525eb1c0",
      "created_at": "2017-09-29T11:31:18.000000",
      "hostname": "DE-IX-001-02-02-01-2",
      "updated_at": "2017-10-06T10:50:03.000000",
      "topic": "engine",
      "host": "DE-IX-001-02-02-01-2",
      "deleted_at": null,
      "id": "06529f29-a258-4eb1-8966-61141c3f0c68"
    },
    {
      "status": "up",
      "binary": "heat-engine",
      "report_interval": 60,
      "engine_id": "3ed4082d-fdc2-442c-b920-627566993626",
      "created_at": "2017-10-06T10:50:28.000000",
      "hostname": "DE-ES-001-03-09-01-1",
      "updated_at": "2017-10-09T14:49:52.000000",
      "topic": "engine",
      "host": "DE-ES-001-03-09-01-1",
      "deleted_at": null,
      "id": "07040a0c-bb75-437f-a59c-7ff9c92ef833"
    }
  ]
}
`
