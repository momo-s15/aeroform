# Homebrew tap (setup guide)

This tap distributes **pre-built binaries** from [GitHub Releases](https://github.com/momo-s15/aeroform/releases) (same files as `install.sh`). You maintain a **separate** repository; Homebrew users never clone the main `aeroform` repo to install.

## You already created the tap repo — what now?

Think of it as **one recipe file** Homebrew reads. That file lives in your **tap** repo (not in `momo-s15/aeroform`).

1. Open **`https://github.com/momo-s15/homebrew-aeroform`** in the browser.
2. Click **Add file** → **Create new file**.
3. In the name box type: **`Formula/aeroform.rb`**  
   (GitHub creates the `Formula` folder for you.)
4. Paste the **entire** Ruby block from the **“Ready-to-paste formula (v1.0.4)”** section below (or use the generic template under **“Manual template (any version)”** when you release a newer version).
5. Click **Commit changes** on **`main`**.

That is the whole setup. After that, anyone can run:

```bash
brew tap momo-s15/aeroform
brew install aeroform
```

**When you tag a new Aeroform version** (e.g. `v1.0.5`): edit `Formula/aeroform.rb` in the tap repo — change **`version`**, every **`v1.0.4`** in the **`url`** lines, and all four **`sha256`** strings. Use the new release’s **`checksums.txt`** on GitHub (same page as the binaries).

---

## Ready-to-paste formula (v1.0.4)

Below matches [release v1.0.4](https://github.com/momo-s15/aeroform/releases/tag/v1.0.4). Replace this whole file when you ship a newer tag.

```ruby
class Aeroform < Formula
  desc "CLI that turns plain English into cloud infrastructure (Terraform)"
  homepage "https://github.com/momo-s15/aeroform"
  version "1.0.4"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/momo-s15/aeroform/releases/download/v1.0.4/aeroform-darwin-arm64"
      sha256 "c6652913c8ea1f91da7f9d6c18308e6c7aa7a3d90e9219fe4fee40fff4e59288"

      def install
        bin.install "aeroform-darwin-arm64" => "aeroform"
      end
    end
    on_intel do
      url "https://github.com/momo-s15/aeroform/releases/download/v1.0.4/aeroform-darwin-amd64"
      sha256 "7ea5a774c63dc4cb6c76b1d8e8103bf9dd652e29fdb260a7420834e00f70f77b"

      def install
        bin.install "aeroform-darwin-amd64" => "aeroform"
      end
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/momo-s15/aeroform/releases/download/v1.0.4/aeroform-linux-arm64"
      sha256 "f3d289dd12dcc1f44e34bc1e85ccc6beec133eb6fe0df6d04fb7b3ae0daaf95e"

      def install
        bin.install "aeroform-linux-arm64" => "aeroform"
      end
    end
    on_intel do
      url "https://github.com/momo-s15/aeroform/releases/download/v1.0.4/aeroform-linux-amd64"
      sha256 "f833fad0b0d71a6bdea1bb8d02eb0c29ff8cfcd0c25f3908a71b0ba1aed1fe5e"

      def install
        bin.install "aeroform-linux-amd64" => "aeroform"
      end
    end
  end

  test do
    system "#{bin}/aeroform", "version"
  end
end
```

---

## Reference: tap repo naming

- Repo URL: **`https://github.com/momo-s15/homebrew-aeroform`**
- Homebrew short name: **`momo-s15/aeroform`** (the `homebrew-` prefix is dropped).

## Manual template (any version)

In that repo, create:

`Formula/aeroform.rb`

Use this template. Replace **`1.0.2`** with your real version (the part **after** `v` in the tag) and replace every **`REPLACE_SHA256_...`** with the matching line from the release’s **`checksums.txt`**.

```ruby
class Aeroform < Formula
  desc "CLI that turns plain English into cloud infrastructure (Terraform)"
  homepage "https://github.com/momo-s15/aeroform"
  version "1.0.2"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/momo-s15/aeroform/releases/download/v1.0.2/aeroform-darwin-arm64"
      sha256 "REPLACE_SHA256_DARWIN_ARM64"

      def install
        bin.install "aeroform-darwin-arm64" => "aeroform"
      end
    end
    on_intel do
      url "https://github.com/momo-s15/aeroform/releases/download/v1.0.2/aeroform-darwin-amd64"
      sha256 "REPLACE_SHA256_DARWIN_AMD64"

      def install
        bin.install "aeroform-darwin-amd64" => "aeroform"
      end
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/momo-s15/aeroform/releases/download/v1.0.2/aeroform-linux-arm64"
      sha256 "REPLACE_SHA256_LINUX_ARM64"

      def install
        bin.install "aeroform-linux-arm64" => "aeroform"
      end
    end
    on_intel do
      url "https://github.com/momo-s15/aeroform/releases/download/v1.0.2/aeroform-linux-amd64"
      sha256 "REPLACE_SHA256_LINUX_AMD64"

      def install
        bin.install "aeroform-linux-amd64" => "aeroform"
      end
    end
  end

  test do
    system "#{bin}/aeroform", "version"
  end
end
```

**Where to get SHA256 values**

After each release, open the tagged release on GitHub and download **`checksums.txt`**, or run:

```bash
curl -fsSL https://github.com/momo-s15/aeroform/releases/download/v1.0.2/checksums.txt
```

Each line looks like:

`<hash>  aeroform-darwin-arm64`

Copy the hash into the matching `sha256` field.

**URLs must use the same tag** as the release (`v1.0.2` in the path matches Git tag `v1.0.2`). Your [release workflow](https://github.com/momo-s15/aeroform/blob/main/.github/workflows/release.yml) publishes:

- `aeroform-darwin-amd64`, `aeroform-darwin-arm64`
- `aeroform-linux-amd64`, `aeroform-linux-arm64`
- `aeroform-windows-amd64.exe`, `aeroform-windows-arm64.exe` (not used by this formula; Windows users use Releases or `install.sh` is N/A)

Commit and push `Formula/aeroform.rb` to the **default branch** (usually `main`).

## What users run

```bash
brew tap momo-s15/aeroform
brew install aeroform
```

Upgrade later:

```bash
brew update
brew upgrade aeroform
```

Equivalent:

```bash
brew install momo-s15/aeroform/aeroform
```

## Each new Aeroform release

1. Tag **`vX.Y.Z`** on `momo-s15/aeroform` and wait for the release workflow to finish.
2. Download **`checksums.txt`** for that tag.
3. In **`homebrew-aeroform`**, edit **`Formula/aeroform.rb`**: bump **`version`**, all **`url`** paths (`.../download/vX.Y.Z/...`), and all four **`sha256`** lines.
4. Commit and push.

Optional: automate step 3 with [GoReleaser `brews`](https://goreleaser.com/customization/homebrew/) and a PAT that can push to `homebrew-aeroform` (see commented example in [`.goreleaser.yaml`](../.goreleaser.yaml)).

## Quick local check (maintainers)

```bash
brew install --build-from-source ./Formula/aeroform.rb   # wrong for binary formula
# Better: push to GitHub, then:
brew tap momo-s15/aeroform https://github.com/momo-s15/homebrew-aeroform
brew reinstall aeroform
aeroform version
```

Or after editing the tap repo:

```bash
brew update
brew upgrade --fetch-HEAD aeroform
```
