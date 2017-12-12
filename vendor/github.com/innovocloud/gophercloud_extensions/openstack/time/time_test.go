package time

import (
	"encoding/json"
	"fmt"
	"testing"
)

type testStruct struct {
	Heartbeat OpenStackTime `json:"heartbeat_timestamp"`
}

func TestWithinStruct(t *testing.T) {
	data := []byte(`{
		"heartbeat_timestamp": "2017-12-12 11:02:41"
	}`)

	var ts testStruct
	err := json.Unmarshal(data, &ts)
	if err != nil {
		t.Error(err)
	}

	if ts.Heartbeat.Unix() != 1513076561 {
		t.Errorf("unexpected result:\n- want: 1513076561\n-  got: %d", ts.Heartbeat.Unix())
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		desc string
		data string
		want int64
	}{
		{
			desc: "testing HeartbeatNoTZ",
			data: "2017-12-12 11:02:41",
			want: 1513076561,
		},
		{
			desc: "testing RFC3339NanoNoTZ",
			data: "2017-12-12 11:03:08.165067",
			want: 1513076588,
		},
		{
			desc: "testing original behavoir time.RFC3339Nano",
			data: "2017-12-12T12:07:30+01:00",
			want: 1513076850,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.desc), func(t *testing.T) {
			data := []byte(`"` + tc.data + `"`)
			var oTime OpenStackTime
			err := oTime.UnmarshalJSON(data)
			if err != nil {
				t.Errorf("Unable to parse %s: %v", tc.data, err)
			}

			if oTime.Unix() != tc.want {
				t.Errorf("unexpected unix timestamp:\n- want: %d\n-  got: %d", tc.want, oTime.Unix())
			}
		})
	}
}
