## Installing Onboard HCE AWS from Release Binaries

Onboard HCE AWS offers pre-compiled binaries available for download on our [releases page](https://github.com/uditgaurav/onboard_hce_aws/releases).

To install, follow these steps:
1. Download the appropriate binary for your platform from the "Assets" section.
2. Rename the downloaded file to `onboard_hce_aws` (or `onboard_hce_aws.exe` for Windows).
3. Move this file to your `$PATH` at your preferred binary installation directory.

### For Linux:

For AMD64 / x86_64:

```bash
[ $(uname -m) = x86_64 ] && curl -Lo ./onboard_hce_aws https://github.com/uditgaurav/onboard_hce_aws/releases/download/0.1.0/cli-linux-amd64
```

For ARM64:

```bash
[ $(uname -m) = aarch64 ] && curl -Lo ./onboard_hce_aws https://github.com/uditgaurav/onboard_hce_aws/releases/download/0.1.0/cli-linux-386
```

After downloading, add execution permissions and move it to your binary installation directory:

```bash
chmod +x ./onboard_hce_aws
sudo mv ./onboard_hce_aws /usr/local/bin/onboard_hce_aws
```
### For MacOS:

For Intel Macs:

```bash
[ $(uname -m) = x86_64 ] && curl -Lo ./onboard_hce_aws https://github.com/uditgaurav/onboard_hce_aws/releases/download/0.1.0/cli-darwin-amd64
```

For M1 / ARM Macs:

```bash
[ $(uname -m) = arm64 ] && curl -Lo ./onboard_hce_aws https://github.com/uditgaurav/onboard_hce_aws/releases/download/0.1.0/cli-darwin-arm64
```

After downloading, add execution permissions and move it to your binary installation directory:

```bash
chmod +x ./onboard_hce_aws
mv ./onboard_hce_aws /some-dir-in-your-PATH/onboard_hce_aws
```

### For Windows:

In PowerShell:

```bash
curl.exe -Lo onboard_hce_aws-windows-amd64.exe https://kind.sigs.k8s.io/dl/v0.20.0/kind-windows-amd64
```

After downloading, move it to your binary installation directory:

```bash
Move-Item .\onboard_hce_aws-windows-amd64.exe c:\some-dir-in-your-PATH\onboard_hce_aws.exe
```