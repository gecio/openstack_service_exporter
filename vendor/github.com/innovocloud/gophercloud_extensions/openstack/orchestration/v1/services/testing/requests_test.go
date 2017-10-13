package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/innovocloud/gophercloud_extensions/openstack/orchestration/v1/services"
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

		th.CheckDeepEquals(t, heatEngine1, actual[0])
		th.CheckDeepEquals(t, heatEngine2, actual[1])
		if len(actual) != 2 {
			t.Errorf("Expected 2 results, saw %d", len(actual))
		}

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}
