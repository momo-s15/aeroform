# Installing Aeroform

## 1. Install script (macOS and Linux — recommended)

No Go toolchain required. Detects OS and CPU, downloads the matching binary from [GitHub Releases](https://github.com/momo-s15/aeroform/releases), and installs to `/usr/local/bin` (uses `sudo` if the directory is not writable).

```bash
curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | sh
```

Pin a specific version (tag must exist on GitHub, e.g. `v1.0.5`):

```bash
curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | AEROFORM_VERSION=v1.0.5 sh
```

Custom location:

```bash
curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | INSTALL_DIR="$HOME/.local/bin" sh
```

Ensure `INSTALL_DIR` is on your `PATH`.

Supported platforms: **linux/amd64**, **linux/arm64**, **darwin/amd64**, **darwin/arm64**. Asset names match the [release workflow](https://github.com/momo-s15/aeroform/blob/main/.github/workflows/release.yml) (`aeroform-{os}-{arch}`).

---

## 2. Homebrew (macOS / Linux)

Uses a separate tap repository (**[momo-s15/homebrew-aeroform](https://github.com/momo-s15/homebrew-aeroform)**) with **`Formula/aeroform.rb`**. Maintainer steps and a ready-to-paste formula are in **[homebrew-tap.md](homebrew-tap.md)**.

```bash
brew tap momo-s15/aeroform
brew install aeroform
brew upgrade aeroform
```

One-liner equivalent:

```bash
brew install momo-s15/aeroform/aeroform
```

Maintainers: with repository secret **`HOMEBREW_TAP_TOKEN`** (PAT scoped to **`homebrew-aeroform`**), [`.github/workflows/release.yml`](../.github/workflows/release.yml) updates **`Formula/aeroform.rb`** on every **`v*.*.*`** tag. Without it, update the tap manually (see **[homebrew-tap.md](homebrew-tap.md)**). Optional: [GoReleaser `brews`](https://goreleaser.com/) — not used in CI; see [`.goreleaser.yaml`](../.goreleaser.yaml).

---

## 3. `go install` (Go developers)

If you already have **Go 1.26.1+**:

```bash
go install github.com/momo-s15/aeroform@latest
```

Pin a version:

```bash
go install github.com/momo-s15/aeroform@v1.0.5
```

Ensure `$(go env GOPATH)/bin` is on your `PATH`. Best for **contributors**; less ideal for students who do not use Go yet.

---

## Windows

There is no one-line installer yet. From [Releases](https://github.com/momo-s15/aeroform/releases), download **`aeroform-windows-amd64.exe`** or **`aeroform-windows-arm64.exe`**, put it in a folder that is on your **PATH** (or add that folder to PATH). You can rename the file to `aeroform.exe` if you prefer.

Interactive menus use **plain line input** on Windows (arrow-key `promptui` prompts are unreliable in some terminals). For **GCP** Simple Mode, you can set **`AEROFORM_GCP_PROJECT_ID`** so you are not blocked if project entry misbehaves.

Future options: **Scoop** or **winget** manifests (similar to Homebrew, separate repo or PR to community buckets).

---

## Verify

```bash
aeroform version
```

---

## What releases contain

Tags matching **`v*.*.*`** trigger [`.github/workflows/release.yml`](../.github/workflows/release.yml), which:

1. Runs `go test ./...`
2. Cross-compiles six binaries into `dist/`
3. Uploads them (and `checksums.txt`) to the GitHub Release

The repository also contains a [`.goreleaser.yaml`](../.goreleaser.yaml) for future use (archives, Homebrew, etc.); the live workflow uses explicit `go build` steps today so asset names stay predictable for `install.sh`.
