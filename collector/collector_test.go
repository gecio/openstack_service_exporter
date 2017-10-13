package collector

import (
	"fmt"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func verifyMetric(ch <-chan prometheus.Metric, want prometheus.Metric) error {
	select {
	case got := <-ch:
		if !reflect.DeepEqual(want, got) {
			return fmt.Errorf("Unexpected metric:\n- want: %+v\n-  got: %+v", want, got)
		}
	case <-time.After(2 * time.Second):
		return fmt.Errorf("Expected to see a metric")
	}
	return nil
}
