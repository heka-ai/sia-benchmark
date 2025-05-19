#!/bin/bash

echo "Installing the API"

curl -L https://git.sia-partners.com/paul.planchon/bench-api/-/raw/master/build/main?inline=false -o /home/ubuntu/api
chmod +x /home/ubuntu/api

cp /home/ubuntu/ec2/api/bench.toml /home/ubuntu/config.toml

sudo mv /home/ubuntu/ec2/api/api.service /etc/systemd/system/api.service
sudo systemctl daemon-reload
sudo systemctl enable api.service
sudo systemctl start api.service

echo "Installation complete"
