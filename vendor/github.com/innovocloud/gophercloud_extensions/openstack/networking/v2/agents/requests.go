package agents

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOpts holds options for filtering agents
type ListOpts struct {
	Host      string `q:"host"`
	AgentType string `q:"agent_type"`
}

// List makes a request against the API to list all network agents
func List(client *gophercloud.ServiceClient, opts *ListOpts) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := gophercloud.BuildQueryString(opts)
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query.String()
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AgentPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get makes a request against the API to get a network agent
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
