#!/bin/bash

echo "Installing system dependencies"

apt-get update

apt-get install -y \
    software-properties-common \
    openssh-client \
    curl \
    git

add-apt-repository ppa:deadsnakes/ppa -y
apt update
apt install -y python3.12 python3.12-venv python3.12-dev

echo "Bootstrapping pip for Python 3.12"
# Ensure pip is installed for Python 3.12 (avoid system pip that requires distutils)
if ! python3.12 -m pip --version >/dev/null 2>&1; then
  if ! python3.12 -m ensurepip --upgrade >/dev/null 2>&1; then
    curl -fsSL https://bootstrap.pypa.io/get-pip.py | python3.12
  fi
fi

echo "Setting up vllm"

echo "Installing setuptools and chz as they are required for gptoss"
python3.12 -m pip install --upgrade pip
python3.12 -m pip install setuptools==68.2.2
python3.12 -m pip install chz

echo "Installing vllm (GPT-OSS)"
python3.12 -m pip install --pre "vllm==0.10.1+gptoss" \
  --extra-index-url https://wheels.vllm.ai/gpt-oss \
  --extra-index-url https://download.pytorch.org/whl/nightly/cu128


VLLM_BIN=$(command -v vllm) || { echo "Error: vllm not found in PATH" >&2; exit 1; }
VLLM_PY=$(head -n1 "$VLLM_BIN" | sed 's/^#!//')
if [ -z "$VLLM_PY" ]; then
  echo "Error: failed to resolve vllm interpreter" >&2
  exit 1
fi
"$VLLM_PY" -m pip install -U xformers

echo "Installation complete"
