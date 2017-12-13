package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/innovocloud/gophercloud_extensions/openstack/networking/v2/agents"
)

func TestListAgents(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	handleAgentListSuccessfully(t)

	pages := 0
	err := agents.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := agents.ExtractAgents(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, ovsAgent, actual[0])
		th.CheckDeepEquals(t, metadataAgent, actual[1])
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
