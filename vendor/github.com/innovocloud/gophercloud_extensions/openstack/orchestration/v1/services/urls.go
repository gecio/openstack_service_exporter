package services

import (
	"github.com/gophercloud/gophercloud"
)

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("services")
}
