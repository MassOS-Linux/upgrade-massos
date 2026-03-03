#!/usr/bin/env bash
set -e

# Don't use CGO until we can confirm it doesn't break things...
export CGO_ENABLED=0

# Use Go in path if available.
if go version &>/dev/null; then
  echo "Going to use the installed copy of Go."
  gotmp=""
  GO_BINARY=go
else
  echo "Going to download a temporary copy of Go."
  gotmp="$(mktemp -d /tmp/.ugmh.XXXXXX)"
  curl -fL https://dl.google.com/go/go1.26.0.linux-amd64.tar.gz -o "$gotmp/go.tar.gz"
  echo "aac1b08a0fb0c4e0a7c1555beb7b59180b05dfc5a3d62e40e9de90cd42f88235 $gotmp/go.tar.gz" | sha256sum -c
  tar -xf "$gotmp/go.tar.gz" -C "$gotmp"
  GO_BINARY="$gotmp/go/bin/go"
fi

# Delete old binary first.
echo "Deleting old ugm-install-helper binary..."
rm -f ugm-install-helper

# Compile new binary.
echo "Compiling ugm-install-helper..."
"$GO_BINARY" build -trimpath ugm-install-helper.go

echo "Stripping ugm-install-helper..."
strip --strip-unneeded ugm-install-helper

# Finish and clean up working directory if needed.
if test -d "$gotmp"; then
  rm -rf "$gotmp"
fi
