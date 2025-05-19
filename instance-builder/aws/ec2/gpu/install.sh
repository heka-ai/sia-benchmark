#!/bin/bash

echo "Installing system dependencies"

apt-get update

apt-get install -y \
    openssh-client \
    curl \
    git

echo "Installing vllm"

pip install vllm==0.8.5.post1

echo "Installation complete"