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
  gover=1.27.0
  case "$(uname -m)" in
    x86_64) goarch=amd64 gosum=675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685 ;;
    aarch64) goarch=arm64 gosum=51798d2c42d0e1c6ed7fd9f48728b4193abac9e8aad6dbac2fe96a81f5909bda ;;
    i?86) goarch=386 gosum=eac4abaca4113170a1cf261b8bf1d38480e61e99deecbc6a14767deb8b19e8ad ;;
    armv*) goarch=armv6l gosum=e337ecd9c321377c0d8832690c2cb10463447c0bd0e65e2e3413dfff63a3435b ;;
    *) echo "Error: System-installed Go needed for $(uname -m)." >&2; exit 1 ;;
  esac
  gotmp="$(mktemp -d /tmp/.ugmh.XXXXXX)"
  curl -fL "https://dl.google.com/go/go$gover.linux-$goarch.tar.gz" -o "$gotmp/go.tar.gz"
  echo "$gosum $gotmp/go.tar.gz" | sha256sum -c
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
