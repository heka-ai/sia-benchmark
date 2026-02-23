#!/bin/bash
#
# GPU image bootstrap: install uv-managed Python 3.12 + vLLM stack.

set -euo pipefail

SECTION="GPU INSTALL"
info() { echo "[$SECTION] $*"; }

info "--------------------------------------------------"
info "Installing system dependencies"
info "--------------------------------------------------"
sudo apt-get update -y
sudo apt-get install -y git curl jq openssh-client

info "--------------------------------------------------"
info "Installing uv for user 'ubuntu'"
info "--------------------------------------------------"
sudo -u ubuntu curl -LsSf https://astral.sh/uv/install.sh | sudo -u ubuntu sh

info "--------------------------------------------------"
info "Preparing uv environment and vLLM packages"
info "--------------------------------------------------"
sudo -u ubuntu bash -c '
  set -euo pipefail
  export PATH=/home/ubuntu/.local/bin:$PATH
  source /home/ubuntu/.local/bin/env
  cd /home/ubuntu

  uv python install 3.12
  uv python pin 3.12
  uv venv

  # vLLM 0.10.2 passes DisabledTqdm to snapshot_download; huggingface_hub>=1.1 calls tqdm_class(..., disable=...),
  # causing "multiple values for keyword argument 'disable'" when DisabledTqdm does super(..., disable=True).
  # Pin huggingface_hub to <1.1 to use the pre-bytes_progress code path that does not pass disable= to tqdm_class.
  # TODO : Upgrade vllm. HF<1.1 is too slow for snapshot_download.
  uv pip install "huggingface_hub>=0.20,<1.1" hf_transfer vllm==0.10.2 --torch-backend=auto --index-strategy unsafe-best-match

  echo "--------------------------------------------------"
  echo "vLLM and dependencies installed successfully"
  echo "Testing vLLM CLI via uv"
  echo "--------------------------------------------------"

  uv run vllm --version
'

info "--------------------------------------------------"
info "GPU environment ready. Use: uv run vllm serve --model ..."
info "--------------------------------------------------"
