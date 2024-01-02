#!/bin/bash

# change directory
cd /opt/opendaylight/bin

# start odl 
./start

# install required features
./client -r 7 "feature:install odl-restconf odl-bgpcep-bgp odl-bgpcep-pcep odl-dlux-core odl-dluxapps-nodes odl-dluxapps-topology odl-dluxapps-yangui odl-dluxapps-yangvisualizer odl-dluxapps-yangman"

# sleep, sometimes errors happen if configuring right after feature installation
sleep 10s

# run user scripts
find /requests/ -type f -regex '.*\.sh' | sort | xargs -I {} sh {}

echo "Started"

$SHELL
