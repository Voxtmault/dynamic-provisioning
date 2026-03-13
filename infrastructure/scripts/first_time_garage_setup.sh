#!/bin/bash

# Make sure the container name is the same as the one defined in the docker-compose.yaml file
garage="docker exec -ti shared_garage /garage"

# Check if layout is already assigned
status_output=$($garage status)
node_status=$(echo "$status_output" | awk 'NR==2 {print $0}')

if echo "$node_status" | grep -q "NO ROLE ASSIGNED"; then
    echo "No layout assigned, creating layout..."
    node_id=$(echo "$node_status" | awk '{print $1}')
    $garage layout assign -z jkt1 -c 1G $node_id
    $garage layout apply --version 1
else
    echo "Layout already assigned, skipping..."
fi