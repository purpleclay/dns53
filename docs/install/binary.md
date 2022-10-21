---
icon: material/package-variant-closed
---

# Installing the Binary

You can use various package managers to install the `dns53` binary. Take your pick.

## Package Managers

### Homebrew

To use [Homebrew](https://brew.sh/):

```sh
brew install purpleclay/tap/dns53
```

### Scoop

To use [Scoop](https://scoop.sh/):

```sh
scoop bucket add purpleclay https://github.com/purpleclay/scoop-bucket.git
scoop install dns53
```

### Apt

To install using the [apt](https://ubuntu.com/server/docs/package-management) package manager:

```sh
echo 'deb [trusted=yes] https://fury.purpleclay.dev/apt/ /' | sudo tee /etc/apt/sources.list.d/purpleclay.list
sudo apt update
sudo apt install -y dns53
```

You may need to install the `ca-certificates` package if you encounter [trust issues](https://gemfury.com/help/could-not-verify-ssl-certificate/) with regard to the Gemfury certificate:

```sh
sudo apt update && sudo apt install -y ca-certificates
```

### Yum

To install using the yum package manager:

```sh
echo '[purpleclay]
name=purpleclay
baseurl=https://fury.purpleclay.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/purpleclay.repo
sudo yum install -y dns53
```

### Aur

To install from the [aur](https://archlinux.org/) using [yay](https://github.com/Jguer/yay):

```sh
yay -S dns53-bin
```

### Linux Packages

Download and manually install one of the `.deb`, `.rpm` or `.apk` packages from the [Releases](https://github.com/purpleclay/dns53/releases) page.

=== "Apt"

    ```sh
    sudo apt install dns53_*.deb
    ```

=== "Yum"

    ```sh
    sudo yum localinstall dns53_*.rpm
    ```

=== "Apk"

    ```sh
    sudo apk add --no-cache --allow-untrusted dns53_*.apk
    ```

### Go Install

```sh
go install github.com/purpleclay/dns53@latest
```

### Bash Script

To install the latest version using a bash script:

```sh
curl https://raw.githubusercontent.com/purpleclay/dns53/main/scripts/install | bash
```

Download a specific version using the `-v` flag. The script uses `sudo` by default but can be disabled through the `--no-sudo` flag.

```sh
curl https://raw.githubusercontent.com/purpleclay/dns53/main/scripts/install | bash -s -- -v v0.1.0 --no-sudo
```

## Manual Download of Binary

Head over to the [Releases](https://github.com/purpleclay/dns53/releases) page on GitHub and download any release artefact. Unpack the `dns53` binary and add it to your `PATH`.

## Verifying a Binary with Cosign

All binaries can be verified using [cosign](https://github.com/sigstore/cosign).

1. Download the checksum files that need to be verified:

   ```sh
   curl -sL https://github.com/purpleclay/dns53/releases/download/v0.8.0/checksums.txt -O
   curl -sL https://github.com/purpleclay/dns53/releases/download/v0.8.0/checksums.txt.sig -O
   curl -sL https://github.com/purpleclay/dns53/releases/download/v0.8.0/checksums.txt.pem -O
   ```

1. Verify the signature of the checksum file:

   ```sh
   cosign verify-blob --cert checksums.txt.pem --signature checksums.txt.sig checksums.txt
   ```

1. Download any release artefact and verify its SHA256 signature matches the entry within the checksum file:

   ```sh
   sha256sum --ignore-missing -c checksums.txt
   ```
