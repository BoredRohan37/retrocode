#!/bin/bash
set -euo pipefail
sudo apt-get update
sudo apt-get install -y cmake g++ pkg-config ripgrep zstd \
  libx11-dev libxft-dev libxext-dev libxinerama-dev libxcursor-dev \
  libxrender-dev libxfixes-dev libglu1-mesa-dev
rustup component add clippy rustfmt rust-analyzer
go install golang.org/x/tools/gopls@latest
echo "[setup] toolchain ready"
