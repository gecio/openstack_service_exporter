package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/innovocloud/gophercloud_extensions/openstack/compute/v2/services"
)

func TestListServices(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	handleServiceListSuccessfully(t)

	pages := 0
	err := services.List(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, novaCompute, actual[0])
		th.CheckDeepEquals(t, novaConductor, actual[1])
		th.CheckDeepEquals(t, novaConsoleAuth, actual[2])
		th.CheckDeepEquals(t, novaScheduler, actual[3])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}
