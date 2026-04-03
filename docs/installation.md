# Installing Aeroform

## 1. Install script (macOS and Linux — recommended)

No Go toolchain required. Detects OS and CPU, downloads the matching binary from [GitHub Releases](https://github.com/momo-s15/aeroform/releases), and installs to `/usr/local/bin` (uses `sudo` if the directory is not writable).

```bash
curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | sh
```

Pin a specific version (tag must exist on GitHub, e.g. `v1.0.2`):

```bash
curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | AEROFORM_VERSION=v1.0.2 sh
```

Custom location:

```bash
curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | INSTALL_DIR="$HOME/.local/bin" sh
```

Ensure `INSTALL_DIR` is on your `PATH`.

Supported platforms: **linux/amd64**, **linux/arm64**, **darwin/amd64**, **darwin/arm64**. Asset names match the [release workflow](https://github.com/momo-s15/aeroform/blob/main/.github/workflows/release.yml) (`aeroform-{os}-{arch}`).

---

## 2. Homebrew (macOS / Linux — planned)

Target experience (not wired yet — track as a **v1.1.0** milestone):

```bash
brew install momo-s15/aeroform/aeroform
brew upgrade aeroform
```

That requires a separate **Homebrew tap** repository (e.g. `homebrew-aeroform`) and a formula. [GoReleaser](https://goreleaser.com/) can publish the formula on each tag once the tap exists; see comments in [`.goreleaser.yaml`](../.goreleaser.yaml) in this repo for a starting point.

Until then, use the **install script** or **`go install`**.

---

## 3. `go install` (Go developers)

If you already have **Go 1.26.1+**:

```bash
go install github.com/momo-s15/aeroform@latest
```

Pin a version:

```bash
go install github.com/momo-s15/aeroform@v1.0.2
```

Ensure `$(go env GOPATH)/bin` is on your `PATH`. Best for **contributors**; less ideal for students who do not use Go yet.

---

## Windows

There is no one-line installer yet. From [Releases](https://github.com/momo-s15/aeroform/releases), download **`aeroform-windows-amd64.exe`** or **`aeroform-windows-arm64.exe`**, put it in a folder that is on your **PATH** (or add that folder to PATH). You can rename the file to `aeroform.exe` if you prefer.

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
