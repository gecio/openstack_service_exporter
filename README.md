# Prometheus OpenStack Service Exporter

Exports service states of various OpenStack Services for Prometheus.

Currently implemented equivalents:

* `openstack volume service list`
* `openstack compute service list`
* `openstack network agent list`
* `openstack orchestration service list` (disabled by default, use `-collector.orchestration` to enable)


You need to pass OpenStack auth environment variables (see `openrc`).

```
OS_AUTH_URL
OS_USERNAME
OS_USERID
OS_PASSWORD
OS_TENANT_ID
OS_TENANT_NAME
OS_DOMAIN_ID
OS_DOMAIN_NAME
```

Note: There is no `OS_PROJECT_NAME` or `OS_PROJECT_ID` available in Gopherclouds env parsing. Use `OS_TENANT_NAME` or `OS_TENANT_ID`.


## Format

```
# HELP openstack_service_blockstorage_enabled Admin status of blockstorage services
# TYPE openstack_service_blockstorage_enabled gauge
openstack_service_blockstorage_enabled{binary="cinder-scheduler",service_host="DE-ES-001-03-09-01-1",zone="nova"} 1
# HELP openstack_service_blockstorage_up Status of blockstorage services
# TYPE openstack_service_blockstorage_up gauge
openstack_service_blockstorage_up{binary="cinder-scheduler",service_host="DE-ES-001-03-09-01-1",zone="nova"} 1

# HELP openstack_service_compute_enabled Admin status of compute services
# TYPE openstack_service_compute_enabled gauge
openstack_service_compute_enabled{binary="nova-compute",id="10",service_host="DE-IX-001-02-02-09-2",zone="ix1"} 1
# HELP openstack_service_compute_up Status of compute services
# TYPE openstack_service_compute_up gauge
openstack_service_compute_up{binary="nova-compute",id="10",service_host="DE-IX-001-02-02-09-2",zone="ix1"} 1

# HELP openstack_service_network_enabled State of network agents
# TYPE openstack_service_network_enabled gauge
openstack_service_network_enabled{agent_type="Open vSwitch agent",binary="neutron-openvswitch-agent",id="f218ea43-92db-4f2b-b8a6-5a9a161f264e",service_host="DE-IX-001-02-02-13-6",topic="N/A",zone=""} 1
# HELP openstack_service_network_up State of network agents
# TYPE openstack_service_network_up gauge
openstack_service_network_up{agent_type="Open vSwitch agent",binary="neutron-openvswitch-agent",id="f218ea43-92db-4f2b-b8a6-5a9a161f264e",service_host="DE-IX-001-02-02-13-6",topic="N/A",zone=""} 1
```
