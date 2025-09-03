#!/bin/bash

echo "Installing system dependencies"

apt-get update

apt-get install -y \
    openssh-client \
    curl \
    git

echo "Installing vllm"

pip install --pre "vllm==0.10.1+gptoss" \
  --extra-index-url https://wheels.vllm.ai/gpt-oss/ \
  --extra-index-url https://download.pytorch.org/whl/nightly/cu128

echo "Installation complete"