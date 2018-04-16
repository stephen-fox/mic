# mic (macOS ISO Creator)

## What is it?
An application that generates macOS installer .iso files from the installers
found in the macOS App Store.

## Warning for High Sierra and later users
Starting with High Sierra, Apple sometimes provides "minimal" installers
through the App Store. These installers are usually less than 1 GB in size.
Otherwise, they appear the same as a normal installer. Make sure to verify
installer before using this application.

For more information, refer to [this forum post](https://www.jamf.com/jamf-nation/discussions/25519/macos-high-sierra-10-13-0-where-to-download-full-installer.

## How do I use it?
First, download a macOS installer from the App Store. Once the download is
complete, execute:
```bash
./mic -i '/Applications/<installer-name>.app' -o ~/Desktop/macos.iso
```

## Where do I download it from?
Refer to the project's releases / tags.
