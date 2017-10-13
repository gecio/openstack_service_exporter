// +build acceptance

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
)

func TestDomainsList(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	if err != nil {
		t.Fatalf("Unable to obtain an identity client: %v", err)
	}

	var iTrue bool = true
	listOpts := domains.ListOpts{
		Enabled: &iTrue,
	}

	allPages, err := domains.List(client, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to list domains: %v", err)
	}

	allDomains, err := domains.ExtractDomains(allPages)
	if err != nil {
		t.Fatalf("Unable to extract domains: %v", err)
	}

	for _, domain := range allDomains {
		tools.PrintResource(t, domain)
	}
}

func TestDomainsGet(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	if err != nil {
		t.Fatalf("Unable to obtain an identity client: %v", err)
	}

	allPages, err := domains.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list domains: %v", err)
	}

	allDomains, err := domains.ExtractDomains(allPages)
	if err != nil {
		t.Fatalf("Unable to extract domains: %v", err)
	}

	domain := allDomains[0]
	p, err := domains.Get(client, domain.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get domain: %v", err)
	}

	tools.PrintResource(t, p)
}
