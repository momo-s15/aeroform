#!/usr/bin/env bash
# Push Formula/aeroform.rb to momo-s15/homebrew-aeroform after a release build.
# Requires: TAP_TOKEN (PAT with contents:write on the tap repo), AEROFORM_TAG (e.g. v1.0.6),
#           CHECKSUMS_FILE (path to checksums.txt from dist/).
set -euo pipefail

TAP_TOKEN="${TAP_TOKEN:?TAP_TOKEN is required}"
AEROFORM_TAG="${AEROFORM_TAG:?AEROFORM_TAG is required}"
CHECKSUMS_FILE="${CHECKSUMS_FILE:?CHECKSUMS_FILE is required}"
TAP_SLUG="${TAP_SLUG:-momo-s15/homebrew-aeroform}"
MAIN_REPO="${MAIN_REPO:-momo-s15/aeroform}"

ver_num="${AEROFORM_TAG#v}"
if [[ "$ver_num" == "$AEROFORM_TAG" ]]; then
  echo "error: AEROFORM_TAG must start with v (got $AEROFORM_TAG)" >&2
  exit 1
fi

hash_for() {
  local name="$1"
  local h
  h="$(awk -v n="$name" '$2 == n { print $1; exit }' "$CHECKSUMS_FILE")"
  if [[ -z "$h" ]]; then
    echo "error: no sha256 line for $name in $CHECKSUMS_FILE" >&2
    exit 1
  fi
  echo "$h"
}

H_DARWIN_ARM64="$(hash_for aeroform-darwin-arm64)"
H_DARWIN_AMD64="$(hash_for aeroform-darwin-amd64)"
H_LINUX_ARM64="$(hash_for aeroform-linux-arm64)"
H_LINUX_AMD64="$(hash_for aeroform-linux-amd64)"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

git clone --depth 1 "https://x-access-token:${TAP_TOKEN}@github.com/${TAP_SLUG}.git" "$tmpdir/tap"
cd "$tmpdir/tap"

mkdir -p Formula

cat > Formula/aeroform.rb <<EOF
class Aeroform < Formula
  desc "CLI that turns plain English into cloud infrastructure (Terraform)"
  homepage "https://github.com/${MAIN_REPO}"
  version "${ver_num}"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/${MAIN_REPO}/releases/download/${AEROFORM_TAG}/aeroform-darwin-arm64"
      sha256 "${H_DARWIN_ARM64}"

      def install
        bin.install "aeroform-darwin-arm64" => "aeroform"
      end
    end
    on_intel do
      url "https://github.com/${MAIN_REPO}/releases/download/${AEROFORM_TAG}/aeroform-darwin-amd64"
      sha256 "${H_DARWIN_AMD64}"

      def install
        bin.install "aeroform-darwin-amd64" => "aeroform"
      end
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/${MAIN_REPO}/releases/download/${AEROFORM_TAG}/aeroform-linux-arm64"
      sha256 "${H_LINUX_ARM64}"

      def install
        bin.install "aeroform-linux-arm64" => "aeroform"
      end
    end
    on_intel do
      url "https://github.com/${MAIN_REPO}/releases/download/${AEROFORM_TAG}/aeroform-linux-amd64"
      sha256 "${H_LINUX_AMD64}"

      def install
        bin.install "aeroform-linux-amd64" => "aeroform"
      end
    end
  end

  test do
    system "#{bin}/aeroform", "version"
  end
end
EOF

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

git add Formula/aeroform.rb
if git diff --cached --quiet; then
  echo "Formula identical to tap main; skipping commit/push."
  exit 0
fi

git commit -m "aeroform ${AEROFORM_TAG}"
git push origin "HEAD:${TAP_BRANCH:-main}"

echo "Homebrew tap updated for ${AEROFORM_TAG}."
