#!/bin/bash

# Make sure the container name is the same as the one defined in the docker-compose.yaml file
garage="docker exec -ti shared_garage /garage" # You can make this as an alias in your own host for convenience, but here we define it as a variable for clarity

# Check if layout is already assigned
status_output=$($garage status)

# Since we are using docker exec, the sample output may look like this:
# 2026-03-14T05:45:47.412917Z  INFO garage_net::netapp: Connected to 127.0.0.1:3901, negotiating handshake...    
# 2026-03-14T05:45:47.455963Z  INFO garage_net::netapp: Connection established to <node_id>    
# ==== HEALTHY NODES ====
# ID                Hostname      Address         Tags  Zone  Capacity          DataAvail  Version
# <node_id>         <host_id>     127.0.0.1:3901              NO ROLE ASSIGNED             v2.1.0

# We want to get the line that contains the node status, which is the second line after the header. We can use awk to do this.
node_status=$(echo "$status_output" | awk 'NR==5 {print $0}') # You can adjust the NR value if the output format changes

if echo "$node_status" | grep -q "NO ROLE ASSIGNED"; then
    echo "No layout assigned, creating layout..."
    node_id=$(echo "$node_status" | awk '{print $1}')
    $garage layout assign -z jkt1 -c 1G $node_id # You could adjust the zone name and capacity as needed
    $garage layout apply --version 1
else
    echo "Layout already assigned, skipping..."
fi