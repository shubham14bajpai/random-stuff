#!/bin/bash

# Loop through all running containers
for container in $(crictl ps -q); do
    # Get the container name correctly
    container_name=$(crictl inspect -o json "$container" | jq -r '.status.metadata.name')

    echo "Checking capabilities for container: $container_name (ID: $container)"

    # Use crictl exec to check capabilities inside the container
    crictl exec -it "$container" bash -c "
        echo 'Capabilities for container $container_name (ID: $container):'

        if command -v capsh &> /dev/null; then
            capsh --print
        else
            echo 'capsh not found, using /proc/self/status'
            grep CapBnd /proc/self/status
        fi
    "
    echo "----------------------------------------------------"
done
