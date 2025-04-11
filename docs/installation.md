# Installation

The **kr8+** binaries can be installed via two methods:

* [ice-bergtech/icetech Homebrew tap](https://github.com/ice-bergtech/homebrew-icetech)
* [Github releases page](https://github.com/ice-bergtech/kr8/releases)

## Homebrew Tap

This is the preferred way to install **kr8+** as it ensures that all dependencies are installed correctly.

```bash
brew tap "ice-bergtech/icetech"
brew install "kr8"
```

## Manually via Releases

kr8+ is a Go binary, which means you can simply download it from the [Github releases page](https://github.com/ice-bergtech/kr8/releases)

Build are produced for:

* **Linux** - amd64 deb, rpm, and apk
* **Darwin** - amd64


```sh
apk add kr8_$VERSION_linux_amd64.apk
```

```sh
deb install kr8_$VERSION_linux_amd64.deb
```

```sh
rpm install kr8_$VERSION_linux_amd64.rpm
```