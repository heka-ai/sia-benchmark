#!/bin/bash

echo "Installing the API"

curl -L https://github.com/heka-ai/sia-benchmark/releases/download/v0.1.0/api-v0.1.0-linux-amd64.tar.gz -o /tmp/api-v0.1.0-linux-amd64.tar.gz
tar -xzf /tmp/api-v0.1.0-linux-amd64.tar.gz -C /home/ubuntu/

chmod +x /home/ubuntu/api

cp /home/ubuntu/ec2/api/bench.toml /home/ubuntu/config.toml

sudo mv /home/ubuntu/ec2/api/api.service /etc/systemd/system/api.service
sudo systemctl daemon-reload
sudo systemctl enable api.service
sudo systemctl start api.service

echo "Installation complete"
