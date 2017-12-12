package routers

import (
	"github.com/gophercloud/gophercloud"
)

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}

// AddResult contains the response body and error from a Create request.
type AddResult struct {
	gophercloud.Result
}
