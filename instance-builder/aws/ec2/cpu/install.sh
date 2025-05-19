#!/bin/bash

echo "Installing system dependencies"

# Install system dependencies
apt-get update && apt-get install -y \
    git \
    curl \
    jq \
    python3 \
    python3-pip \
    python3-venv \
    openssh-client

echo "Installing Python packages"

# Install Python packages
pip install -r requirements.txt  

echo "Installation complete"