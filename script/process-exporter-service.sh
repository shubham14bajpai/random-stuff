#!/bin/bash

sudo apt-get update

wget https://github.com/ncabatoff/process-exporter/releases/download/v0.7.10/process-exporter-0.7.10.linux-amd64.tar.gz

sudo groupadd -f process_exporter
sudo useradd -g process_exporter --no-create-home --shell /bin/false process_exporter
sudo mkdir /etc/process_exporter
sudo chown process_exporter:process_exporter /etc/process_exporter

tar -xvf process-exporter-0.7.10.linux-amd64.tar.gz
mv process-exporter-0.7.10.linux-amd64 process_exporter-files

sudo cp process_exporter-files/process-exporter /usr/bin/
sudo chown process_exporter:process_exporter /usr/bin/process-exporter

sudo wget -c https://raw.githubusercontent.com/shubham14bajpai/random-stuff/main/script/process-exporter.yml -O /etc/process_exporter/process-exporter.yaml

sudo chown process_exporter:process_exporter /etc/process_exporter/process-exporter.yaml

sudo echo "[Unit]
Description=Process Exporter for Prometheus
Documentation=https://github.com/ncabatoff/process-exporter
Wants=network-online.target
After=network-online.target

[Service]
User=process_exporter
Group=process_exporter
Type=simple
Restart=on-failure
ExecStart=/usr/bin/process-exporter \
  --config.path /etc/process_exporter/process-exporter.yaml \
  --web.listen-address=:9256

[Install]
WantedBy=multi-user.target" | sudo tee /usr/lib/systemd/system/process_exporter.service

sudo chmod 664 /usr/lib/systemd/system/process_exporter.service

sudo systemctl daemon-reload
sudo systemctl start process_exporter

sudo systemctl enable process_exporter.service

