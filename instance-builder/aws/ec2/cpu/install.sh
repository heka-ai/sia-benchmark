#!/bin/bash

echo "Installing system dependencies"

# Install system dependencies
apt-get update && apt-get install -y \
    git \
    curl \
    jq \
    openssh-client

echo "Installing Python packages"

# Install uv for ubuntu user
sudo -u ubuntu curl -LsSf https://astral.sh/uv/install.sh | sudo -u ubuntu sh
sudo -u ubuntu bash -c "cd /home/ubuntu && source /home/ubuntu/.local/bin/env && \
    uv python install 3.12 && \
    uv python pin 3.12 && \
    uv venv && \
    uv pip install -r /home/ubuntu/ec2/cpu/requirements.txt"

echo "Installation complete"
