package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/innovocloud/gophercloud_extensions/openstack/networking/v2/agents"
	osTime "github.com/innovocloud/gophercloud_extensions/openstack/time"
)

const agentListBody = `
{
  "agents": [
    {
      "binary": "neutron-openvswitch-agent",
      "description": null,
      "availability_zone": null,
      "heartbeat_timestamp": "2017-10-09 13:24:53",
      "admin_state_up": true,
      "alive": true,
      "id": "02306fa8-d790-4272-8e2f-c3a52c8e9f15",
      "topic": "N/A",
      "host": "DE-IX-001-02-02-13-7",
      "agent_type": "Open vSwitch agent",
      "started_at": "2017-10-06 11:10:44",
      "created_at": "2017-08-02 08:55:58",
      "configurations": {
        "ovs_hybrid_plug": true,
        "in_distributed_mode": false,
        "datapath_type": "system",
        "vhostuser_socket_dir": "/var/run/openvswitch",
        "tunneling_ip": "10.90.248.23",
        "arp_responder_enabled": true,
        "devices": 0,
        "ovs_capabilities": {
          "datapath_types": [
            "netdev",
            "system"
          ],
          "iface_types": [
            "geneve",
            "gre",
            "internal",
            "ipsec_gre",
            "lisp",
            "patch",
            "stt",
            "system",
            "tap",
            "vxlan"
          ]
        },
        "log_agent_heartbeats": false,
        "l2_population": true,
        "tunnel_types": [
          "vxlan"
        ],
        "extensions": [
          "qos"
        ],
        "enable_distributed_routing": false,
        "bridge_mappings": {}
      }
    },
    {
      "binary": "neutron-metadata-agent",
      "description": null,
      "availability_zone": null,
      "heartbeat_timestamp": "2017-10-09 13:25:12",
      "admin_state_up": true,
      "alive": true,
      "id": "0b4d3a07-b680-40da-88a3-8de99fd34100",
      "topic": "N/A",
      "host": "DE-IX-001-02-02-01-6",
      "agent_type": "Metadata agent",
      "started_at": "2017-10-06 10:55:09",
      "created_at": "2017-08-02 08:57:04",
      "configurations": {
        "log_agent_heartbeats": false,
        "nova_metadata_port": 8775,
        "nova_metadata_ip": "10.90.64.0",
        "metadata_proxy_socket": "/var/lib/neutron/kolla/metadata_proxy"
      }
    }
  ]
}
`

var (
	ovsAgent = agents.Agent{
		ID:               "02306fa8-d790-4272-8e2f-c3a52c8e9f15",
		Binary:           "neutron-openvswitch-agent",
		Description:      "",
		AvailabilityZone: "",
		AdminStateUp:     true,
		Alive:            true,
		Topic:            "N/A",
		Host:             "DE-IX-001-02-02-13-7",
		AgentType:        "Open vSwitch agent",
		Heartbeat:        osTime.OpenStackTime{Time: time.Date(2017, 10, 9, 13, 24, 53, 0, time.UTC)},
	}
	metadataAgent = agents.Agent{
		ID:               "0b4d3a07-b680-40da-88a3-8de99fd34100",
		Binary:           "neutron-metadata-agent",
		Description:      "",
		AvailabilityZone: "",
		AdminStateUp:     true,
		Alive:            true,
		Topic:            "N/A",
		Host:             "DE-IX-001-02-02-01-6",
		AgentType:        "Metadata agent",
		Heartbeat:        osTime.OpenStackTime{Time: time.Date(2017, 10, 9, 13, 25, 12, 0, time.UTC)},
	}
)

func handleAgentListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, agentListBody)
	})
}
