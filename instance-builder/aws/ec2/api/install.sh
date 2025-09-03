#!/bin/bash

echo "Installing the API requirements"

chmod +x /home/ubuntu/api || exit 1

cp /home/ubuntu/ec2/api/bench.toml /home/ubuntu/config.toml

# Ensure EnvironmentFile exists for systemd unit (optional but present)
sudo mkdir -p /etc/default
sudo touch /etc/default/benchmark-api

sudo mv /home/ubuntu/ec2/api/api.service /etc/systemd/system/api.service
sudo systemctl daemon-reload
sudo systemctl enable api.service
sudo systemctl start api.service
sudo systemctl status api.service

echo "Api setup complete"
