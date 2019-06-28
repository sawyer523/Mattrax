# Contributing

Welcome to Mattrax! If you encounter any troubles please feel free to contact a maintainer.

## Code of Conduct

Mattrax uses the [Contributor Covenant](https://www.contributor-covenant.org/version/1/4/code-of-conduct) if you believe someone is in violation of that please contact a maintainer.

## Specifications

Mattrax uses the following specifications. Please understand and use them when they are relevant, that includes pull requests.

- [Contributor Covenant 1.4](https://www.contributor-covenant.org/version/1/4/code-of-conduct)
- [Conventional Commits 1.0.0-beta.4](https://www.conventionalcommits.org/en/v1.0.0-beta.4/)
- [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)
- [Keep A Changelog 1.0.0](https://keepachangelog.com/en/1.0.0/)

Mattrax uses the official MDM protocols built into the respective device that it is managing. Please understand and abide by them at all times.

- [Apple MDM Protocol Reference](https://developer.apple.com/library/content/documentation/Miscellaneous/Reference/MobileDeviceManagementProtocolRef/3-MDM_Protocol/MDM_Protocol.html)
- [MS-MDE2](https://winprotocoldoc.blob.core.windows.net/productionwindowsarchives/MS-MDE2/%5bMS-MDE2%5d.pdf), [MS-MDM](https://winprotocoldoc.blob.core.windows.net/productionwindowsarchives/MS-MDM/%5bMS-MDM%5d.pdf) and other related protocols

## Running the project locally

To run Mattrax from the sources, you will need the latest version of [Go Lang](https://golang.org/dl/) installed.

```bash
git clone https://github.com/mattrax/Mattrax.git && cd Mattrax
go run ./cmd/mattrax
```
