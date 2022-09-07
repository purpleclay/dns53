# Install

There are many different ways to install dns53. You can install the binary using either a supported package manager, manually, or by compiling the source yourself.

## Installing the binary

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

You may need to install the `ca-certificates` package if you encounter [trust issues](https://gemfury.com/help/could-not-verify-ssl-certificate/) with regards to the gemfury certificate:

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

A specific version can be downloaded by using the `-v` flag. By default the script uses `sudo`, which can be turned off by using the `--no-sudo` flag.

```sh
curl https://raw.githubusercontent.com/purpleclay/dns53/main/scripts/install | bash -s -- -v v0.1.0 --no-sudo
```

### Manually

Binary downloads of dns53 can be found on the [Releases](https://github.com/purpleclay/dns53/releases) page. Unpack the dns53 binary and add it to your `PATH`.

## Compiling from source

dns53 is written using [Go 1.18+](https://go.dev/doc/install) and should be installed along with [go-task](https://taskfile.dev/#/installation), as it is preferred over using make.

Then clone the code from GitHub:

```sh
git clone https://github.com/purpleclay/dns53.git
cd dns53
```

Build dns53:

```sh
task
```

And check that everything works:

```sh
./dns53 version
```

## Verifying Artefacts

All verification is carried out using cosign and it must be [installed](https://docs.sigstore.dev/cosign/installation) before proceeding.

### Binaries

All binaries can be verified using the checksum file, which has been signed using cosign.

1. Download the checksum files that need to be verified:

   ```sh
   curl -sL https://github.com/purpleclay/dns53/releases/download/v0.1.0/checksums.txt -O
   curl -sL https://github.com/purpleclay/dns53/releases/download/v0.1.0/checksums.txt.sig -O
   curl -sL https://github.com/purpleclay/dns53/releases/download/v0.1.0/checksums.txt.pem -O
   ```

1. Verify the signature of the checksum file:

   ```sh
   cosign verify-blob --cert checksums.txt.pem --signature checksums.txt.sig checksums.txt
   ```

1. Download any release artefact and verify its SHA256 signature matches the entry within the checksum file:

   ```sh
   sha256sum --ignore-missing -c checksums.txt
   ```

!!!warning "Don't mix versions"

    For checksum verification to work, all artefacts must be downloaded from the same release
