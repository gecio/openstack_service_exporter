package routers

import (
	"github.com/gophercloud/gophercloud"
)

func listURL(client *gophercloud.ServiceClient, agentID string) string {
	return client.ServiceURL("agents", agentID, "l3-routers")
}

func deleteURL(client *gophercloud.ServiceClient, agentID string, routerID string) string {
	return client.ServiceURL("agents", agentID, "l3-routers", routerID)
}

func addURL(client *gophercloud.ServiceClient, agentID string) string {
	return client.ServiceURL("agents", agentID, "l3-routers")
}
