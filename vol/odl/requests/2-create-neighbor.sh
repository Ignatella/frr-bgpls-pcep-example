#!/bin/bash

curl -k -X POST -sS \
  --header 'Authorization: Basic YWRtaW46YWRtaW4=' \
  --header 'Content-Type: application/json' \
  --data '{
    "neighbor": [
        {
            "neighbor-address": "172.20.0.5",
            "timers": {
                "config": {
                    "keepalive-interval": 30,
                    "hold-time": 180,
                    "connect-retry": 100
                }
            },
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
            "transport": {
                "config": {
                    "mtu-discovery": false,
                    "remote-port": 1790,
                    "passive-mode": false
                }
            },
            "config": {
                "send-community": "NONE",
                "peer-as": 1,
                "peer-type": "INTERNAL",
                "route-flap-damping": false
            }
        }
    ]
}' \
   'http://172.20.0.3:8181/restconf/config/openconfig-network-instance:network-instances/network-instance/global-bgp/openconfig-network-instance:protocols/protocol/openconfig-policy-types:BGP/bgp-test/bgp/neighbors'



