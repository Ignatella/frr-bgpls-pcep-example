#!/bin/bash

set -e

curl -k -X POST -sS \
  -f \
  --header 'Authorization: Basic YWRtaW46YWRtaW4=' \
  --header 'Content-Type: application/json' \
  --data '{
    "protocol": [
        {
            "identifier": "openconfig-policy-types:BGP",
            "name": "bgp-test",
            "bgp-openconfig-extensions:bgp": {
                "global": {
                    "afi-safis": {
                        "afi-safi": [
                            {"afi-safi-name": "openconfig-bgp-types:IPV4-UNICAST"},
                            {"afi-safi-name": "bgp-openconfig-extensions:IPV4-FLOW"},
                            {"afi-safi-name": "bgp-openconfig-extensions:LINKSTATE"},
                            {"afi-safi-name": "openconfig-bgp-types:IPV4-LABELLED-UNICAST"},
                            {"afi-safi-name": "bgp-openconfig-extensions:IPV6-L3VPN-FLOW"},
                            {"afi-safi-name": "openconfig-bgp-types:IPV6-LABELLED-UNICAST"},
                            {"afi-safi-name": "openconfig-bgp-types:L3VPN-IPV4-UNICAST"},
                            {"afi-safi-name": "openconfig-bgp-types:L3VPN-IPV6-UNICAST"},
                            {"afi-safi-name": "bgp-openconfig-extensions:IPV6-FLOW"},
                            {"afi-safi-name": "openconfig-bgp-types:L2VPN-EVPN"},
                            {"afi-safi-name": "bgp-openconfig-extensions:IPV4-L3VPN-FLOW"},
                            {"afi-safi-name": "openconfig-bgp-types:IPV6-UNICAST"}
                        ]
                    },
                    "config": {
                        "router-id": "172.20.0.3",
                        "as": 1
                    }
                }
            }
        }
    ]
}' \
   'http://172.20.0.3:8181/restconf/config/openconfig-network-instance:network-instances/network-instance/global-bgp/openconfig-network-instance:protocols/' \
|| exit 1

