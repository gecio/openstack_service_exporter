package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/innovocloud/gophercloud_extensions/openstack/compute/v2/services"
	osTime "github.com/innovocloud/gophercloud_extensions/openstack/time"
)

const serviceListBody = `
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
      "state": "up",
      "disabled_reason": null,
      "id": 185,
      "updated_at": "2017-10-09T13:06:16.000000"
    },
    {
      "status": "enabled",
      "binary": "nova-consoleauth",
      "host": "DE-ES-001-03-09-01-3",
      "zone": "internal",
      "state": "up",
      "disabled_reason": null,
      "id": 210,
      "updated_at": "2017-10-09T13:06:17.000000"
    },
    {
      "status": "enabled",
      "binary": "nova-scheduler",
      "host": "DE-ES-001-03-09-01-3",
      "zone": "internal",
      "state": "up",
      "disabled_reason": null,
      "id": 235,
      "updated_at": "2017-10-09T13:06:16.000000"
    }
  ]
}
`

var (
	novaCompute = services.Service{
		ID:             10,
		Status:         "enabled",
		Binary:         "nova-compute",
		Host:           "DE-IX-001-02-02-09-2",
		Zone:           "ix1",
		State:          "up",
		DisabledReason: "",
		Updated:        osTime.OpenStackTime{Time: time.Date(2017, 10, 9, 13, 06, 14, 0, time.UTC)},
	}

	novaConductor = services.Service{
		ID:             185,
		Status:         "enabled",
		Binary:         "nova-conductor",
		Host:           "DE-IX-001-02-02-01-1",
		Zone:           "internal",
		State:          "up",
		DisabledReason: "",
		Updated:        osTime.OpenStackTime{Time: time.Date(2017, 10, 9, 13, 06, 16, 0, time.UTC)},
	}

	novaConsoleAuth = services.Service{
		ID:             210,
		Status:         "enabled",
		Binary:         "nova-consoleauth",
		Host:           "DE-ES-001-03-09-01-3",
		Zone:           "internal",
		State:          "up",
		DisabledReason: "",
		Updated:        osTime.OpenStackTime{Time: time.Date(2017, 10, 9, 13, 06, 17, 0, time.UTC)},
	}

	novaScheduler = services.Service{
		ID:             235,
		Status:         "enabled",
		Binary:         "nova-scheduler",
		Host:           "DE-ES-001-03-09-01-3",
		Zone:           "internal",
		State:          "up",
		DisabledReason: "",
		Updated:        osTime.OpenStackTime{Time: time.Date(2017, 10, 9, 13, 06, 16, 0, time.UTC)},
	}
)

func handleServiceListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, serviceListBody)
	})
}
