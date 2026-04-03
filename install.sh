#!/bin/sh
# Aeroform installer for macOS and Linux.
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | sh
# Optional:
#   AEROFORM_VERSION=v1.0.2  — pin a release tag (default: latest)
#   INSTALL_DIR=/usr/local/bin  — install location (default: /usr/local/bin)

set -eu

REPO="momo-s15/aeroform"
BINARY_BASE="aeroform"

# --- detect OS ---
uname_s=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$uname_s" in
  linux)  os="linux" ;;
  darwin) os="darwin" ;;
  *)
    echo "aeroform install: unsupported OS '$uname_s' (only Linux and macOS)." >&2
    echo "Windows: download aeroform-windows-*.zip from https://github.com/${REPO}/releases" >&2
    exit 1
    ;;
esac

# --- detect arch ---
uname_m=$(uname -m)
case "$uname_m" in
  x86_64|amd64)     arch="amd64" ;;
  arm64|aarch64)    arch="arm64" ;;
  *)
    echo "aeroform install: unsupported CPU '$uname_m' (need amd64 or arm64)." >&2
    exit 1
    ;;
esac

asset="${BINARY_BASE}-${os}-${arch}"
version="${AEROFORM_VERSION:-latest}"

if [ "$version" = "latest" ]; then
  url="https://github.com/${REPO}/releases/latest/download/${asset}"
else
  case "$version" in
    v*) ;;
    *) version="v${version}" ;;
  esac
  url="https://github.com/${REPO}/releases/download/${version}/${asset}"
fi

install_dir="${INSTALL_DIR:-/usr/local/bin}"
target="${install_dir}/${BINARY_BASE}"

tmp=""
cleanup() {
  [ -n "${tmp}" ] && rm -f "${tmp}"
}
trap cleanup EXIT INT TERM

echo "Installing aeroform (${os}/${arch}) from ${url}" >&2
tmp=$(mktemp)
if ! curl -fsSL -o "${tmp}" "${url}"; then
  echo "aeroform install: download failed. Check the release exists: https://github.com/${REPO}/releases" >&2
  exit 1
fi

chmod 0755 "${tmp}"

if [ -w "${install_dir}" ] 2>/dev/null; then
  mv "${tmp}" "${target}"
  tmp=""
else
  echo "aeroform install: using sudo to write ${target}" >&2
  sudo mv "${tmp}" "${target}"
  tmp=""
fi

echo "Installed ${target}" >&2
case ":${PATH}:" in
  *:"${install_dir}":*) ;;
  *) echo "Note: add ${install_dir} to your PATH if aeroform is not found." >&2 ;;
esac
command -v "${BINARY_BASE}" >/dev/null 2>&1 && "${BINARY_BASE}" version || true
