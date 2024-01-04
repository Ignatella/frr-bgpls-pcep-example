#!/bin/bash

# change directory
cd /opt/opendaylight/bin

# start odl 
./start

# install required features
./client -r 100 "feature:install odl-restconf odl-bgpcep-bgp odl-bgpcep-pcep odl-dlux-core odl-dluxapps-nodes odl-dluxapps-topology odl-dluxapps-yangui odl-dluxapps-yangvisualizer odl-dluxapps-yangman"

# sleep, sometimes errors happen if configuring right after feature installation
# sleep 10s

# run user scripts until success
find /requests/ -type f -regex '.*\.sh' | sort | while read -r script; do
    success=false

    while [ "$success" == false ]; do
        # Execute the script and check the exit status
        sh "$script"

        # Check the exit status
        if [ $? -eq 0 ]; then
            success=true  # Set success to true if the script succeeds
        else
            echo "Error: Failed to execute $script. Retrying..."
            sleep 1  # Add a small delay before retrying
        fi
    done
done



# run user scripts
# find /requests/ -type f -regex '.*\.sh' | sort | xargs -I {} sh {}

echo "Started"

$SHELL
