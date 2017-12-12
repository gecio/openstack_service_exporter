package routers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/pagination"
)

// List makes a request against the API to list all network agents
func List(client *gophercloud.ServiceClient, agentID string) pagination.Pager {
	url := listURL(client, agentID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return routers.RouterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Remove a router from an agent
func Remove(client *gophercloud.ServiceClient, agentID string, routerID string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, agentID, routerID), nil)
	return
}

// Add a router to an agent
func Add(client *gophercloud.ServiceClient, agentID string, routerID string) (r AddResult) {
	data := struct {
		RouterID string `json:"router_id"`
	}{RouterID: routerID}
	_, r.Err = client.Post(addURL(client, agentID), data, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200, 201}})
	return
}
