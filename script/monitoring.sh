## Usage: ./monitoring.sh
## Make sure the .yml are placed at the below mentioned directories
## or modify the commands to use the correct location as per your changes

#!/bin/bash

## Node exporter Docker install
sudo docker run -d --net="host" --pid="host" -v "/:/host:ro,rslave" quay.io/prometheus/node-exporter:latest --path.rootfs=/host

## Process exporter Docker install

docker run -d --rm -p 9256:9256 --privileged -v /proc:/host/proc -v `pwd`:/config ncabatoff/process-exporter --procfs /host/proc -config.path /config/process-exporter.yml

## Prometheus Docker install

sudo docker run -d --name prometheus -p 9090:9090 -v ~/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

## Grafana Docker install

docker run -d --name=grafana -p 3456:3000 grafana/grafana


