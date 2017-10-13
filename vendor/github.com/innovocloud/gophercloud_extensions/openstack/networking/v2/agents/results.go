package agents

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// AgentPage abstracts the raw results of making a List() request against
// the API. As OpenStack extensions may freely alter the response bodies of
// structures returned to the client, you may only safely access the data
// provided through the ExtractAgent call.
type AgentPage struct {
	pagination.LinkedPageBase
}

// Agent represents a blockstorage service in the OpenStack cloud.
type Agent struct {
	ID               string `json:"id"`
	Binary           string `json:"binary"`
	Description      string `json:"description"`
	AvailabilityZone string `json:"availability_zone"`
	AdminStateUp     bool   `json:"admin_state_up"`
	Alive            bool   `json:"alive"`
	Topic            string `json:"topic"`
	Host             string `json:"host"`
	AgentType        string `json:"agent_type"`
}

// IsEmpty returns true if a page contains no Agents results.
func (r AgentPage) IsEmpty() (bool, error) {
	services, err := ExtractAgents(r)
	return len(services) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r AgentPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"agents_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractAgents interprets the results of a single page from a List() call,
// producing a slice of Agents entities.
func ExtractAgents(r pagination.Page) ([]Agent, error) {
	var s []Agent
	err := r.(AgentPage).Result.ExtractIntoSlicePtr(&s, "agents")
	return s, err
}
