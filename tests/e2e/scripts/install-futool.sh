#!/bin/bash

if hash futool 2>/dev/null; then
  echo "[install-futool.sh] futool is already installed. skipping installation."
  exit 0
fi

echo "[install-futool.sh] installing futool."
cd futool || exit 1
make install
