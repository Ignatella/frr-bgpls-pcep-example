#!/bin/bash

set -e

curl -k -X POST -sS \
  -f \
  --header 'Authorization: Basic YWRtaW46YWRtaW4=' \
  --header 'Content-Type: application/xml' \
  --data '<topology xmlns="urn:TBD:params:xml:ns:yang:network-topology"> \
    <topology-id>test-linkstate-topology</topology-id> \
    <topology-types> \
        <bgp-linkstate-topology xmlns="urn:opendaylight:params:xml:ns:yang:odl-bgp-topology-types"></bgp-linkstate-topology> \
    </topology-types> \
    <rib-id xmlns="urn:opendaylight:params:xml:ns:yang:odl-bgp-topology-config">bgp-test</rib-id> \
</topology>' \
   'http://172.20.0.3:8181/restconf/config/network-topology:network-topology' \
|| exit 1

