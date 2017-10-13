package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/innovocloud/gophercloud_extensions/openstack/orchestration/v1/services"
)

const serviceListBody = `
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

var (
	heatEngine1 = services.Service{
		ID:             "06529f29-a258-4eb1-8966-61141c3f0c68",
		EngineID:       "1ab44a4a-c62b-476a-96a8-a3e4525eb1c0",
		Status:         "down",
		Binary:         "heat-engine",
		ReportInterval: 60,
		Host:           "DE-IX-001-02-02-01-2",
		Hostname:       "DE-IX-001-02-02-01-2",
		Topic:          "engine",
	}

	heatEngine2 = services.Service{
		ID:             "07040a0c-bb75-437f-a59c-7ff9c92ef833",
		EngineID:       "3ed4082d-fdc2-442c-b920-627566993626",
		Status:         "up",
		Binary:         "heat-engine",
		ReportInterval: 60,
		Host:           "DE-ES-001-03-09-01-1",
		Hostname:       "DE-ES-001-03-09-01-1",
		Topic:          "engine",
	}
)

func handleServiceListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, serviceListBody)
	})
}
