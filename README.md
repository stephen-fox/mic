# mic (macOS ISO Creator)

## What is it?
An application that generates macOS installer .iso files from the installers
found in the macOS App Store.

## Issues with High Sierra and later

#### Minimal App Store installers
Starting with High Sierra, Apple sometimes provides "minimal" macOS installers
through the App Store. These installers are usually less than 1 GB in size.
Otherwise, they appear the same as a normal installer. Make sure to verify
installer before using this application. The application cannot use minimal
installers to create a .iso.

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
Refer to the project's releases / tags.
