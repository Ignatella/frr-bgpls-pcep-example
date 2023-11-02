#!/bin/bash

# change directory
cd /opt/opendaylight/bin

# start odl 
./start

# install required features
./client -r 7 "feature:install odl-restconf odl-bgpcep-bgp odl-bgpcep-pcep odl-dlux-core odl-dluxapps-nodes odl-dluxapps-topology odl-dluxapps-yangui odl-dluxapps-yangvisualizer odl-dluxapps-yangman"

$SHELL
