package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/innovocloud/gophercloud_extensions/openstack/blockstorage/v2/services"
)

const serviceListBody = `
{
  "services": [
    {
      "status": "enabled",
      "binary": "cinder-scheduler",
      "zone": "nova",
      "state": "up",
      "updated_at": "2017-10-09T13:53:33.000000",
      "host": "DE-IX-001-02-02-01-2",
      "disabled_reason": null
    },
    {
      "status": "enabled",
      "binary": "cinder-scheduler",
      "zone": "nova",
      "state": "up",
      "updated_at": "2017-10-09T13:53:28.000000",
      "host": "DE-ES-001-03-09-01-3",
      "disabled_reason": null
    }
  ]
}
`

var (
	cinderScheduler1 = services.Service{
		Status:         "enabled",
		Binary:         "cinder-scheduler",
		Host:           "DE-IX-001-02-02-01-2",
		Zone:           "nova",
		State:          "up",
		DisabledReason: "",
		Updated:        "2017-10-09T13:53:33.000000",
	}

	cinderScheduler2 = services.Service{
		Status:         "enabled",
		Binary:         "cinder-scheduler",
		Host:           "DE-ES-001-03-09-01-3",
		Zone:           "nova",
		State:          "up",
		DisabledReason: "",
		Updated:        "2017-10-09T13:53:28.000000",
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
