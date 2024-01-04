#!/bin/bash


# Delete routers and neighbours
curl -k -X DELETE -sS \
  --header 'Authorization: Basic YWRtaW46YWRtaW4=' \
  'http://172.20.0.3:8181/restconf/config/openconfig-network-instance:network-instances/network-instance/global-bgp/openconfig-network-instance:protocols/' 

# Delete topology
curl -k -X DELETE -sS \
  --header 'Authorization: Basic YWRtaW46YWRtaW4=' \
  'http://172.20.0.3:8181/restconf/config/network-topology:network-topology' 
