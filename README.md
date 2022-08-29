# mic (macOS ISO Creator)

## What is it?
An application that generates macOS installer .iso files from the installers
found in the macOS App Store.

## Issues with High Sierra and later

#### Minimal App Store installers
Starting with High Sierra, Apple sometimes provides "minimal" macOS installers
through the App Store. These installers are usually less than 1 GB in size.
Otherwise, they appear the same as a normal installer. Make sure to verify the
installer before using this tool. The tool does not support minimal installers.

For more information, refer to [this forum post](https://www.jamf.com/jamf-nation/discussions/25519/macos-high-sierra-10-13-0-where-to-download-full-installer).

#### VirtualBox EFI errors
If you plan to use a High Sierra .iso in VirtualBox, you may encounter EFI
errors after the macOS installation process completes. If this occurs,
VirtualBox will automatically boot you into the EFI shell. You can work around
this problem by executing the following commands in the EFI shell:

```bash
fs1:
cd "macOS Install Data"
cd "Locked Files"
cd "Boot Files"
boot.efi
```

## How do I use it?
First, download a macOS installer from the App Store. Once the download is
complete, execute:

```bash
./mic -i '/Applications/<installer-name>.app' -o ~/Desktop/macos.iso
```

## Where do I download the tool?
Since this is a Go (Golang) application, the preferred method of installation
is using `go install`. This automates downloading and building Go applications
from source in a secure manner. By default, this copies applications
into `~/go/bin/`.

You must first [install Go](https://golang.org/doc/install). After
installing Go, simply run the following command to install the application:

```sh
go install github.com/stephen-fox/mic/cmd/mic@latest
# If successful, the exectuable should be in "~/go/bin/".
```
