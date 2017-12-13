package time

import (
	"errors"
	"time"
)

const (
	RFC3339NanoNoTZ = "2006-01-02T15:04:05.999999999"
	HeartbeatNoTZ   = "2006-01-02 15:04:05"
)

// OpenStackTime is just a wrap around time.Time to use a custom UnmarshalJSON
type OpenStackTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 RFC3339NanoNoTZ or HeartbeatNoTZ format.
func (o *OpenStackTime) UnmarshalJSON(data []byte) error {
	if o == nil {
		return errors.New("no such object")
	}
	var err error
	err = o.Time.UnmarshalJSON(data)
	// If there is no error, parsing was successfull
	if err == nil {
		return nil
	}

	o.Time, err = time.Parse(`"`+HeartbeatNoTZ+`"`, string(data))
	if err == nil {
		return nil
	}

	o.Time, err = time.Parse(`"`+RFC3339NanoNoTZ+`"`, string(data))
	return err
}
