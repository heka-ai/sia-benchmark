#!/bin/bash

echo "Installing the API"

# Ensure Go toolchain is available
if ! command -v go >/dev/null 2>&1; then
	sudo apt-get update -y
	sudo apt-get install -y golang-go
fi

# Build API from local repo
cd /home/ubuntu/benchmark-cli/api || exit 1

go mod download

go build -o /home/ubuntu/api ./cmd || exit 1

chmod +x /home/ubuntu/api

cp /home/ubuntu/ec2/api/bench.toml /home/ubuntu/config.toml

sudo mv /home/ubuntu/ec2/api/api.service /etc/systemd/system/api.service
sudo systemctl daemon-reload
sudo systemctl enable api.service
sudo systemctl start api.service

echo "Installation complete"
