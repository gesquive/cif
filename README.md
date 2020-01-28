# cig
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/cig/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/gesquive/cig)
[![Build Status](https://img.shields.io/circleci/build/github/gesquive/cig?style=flat-square)](https://circleci.com/gh/gesquive/cig)
[![Coverage Report](https://img.shields.io/codecov/c/gh/gesquive/cig?style=flat-square)](https://codecov.io/gh/gesquive/cig)

Formats PEM certificates in a nice human readable (mkcert.org) format

## Installing

### Compile
This project has only been tested with go1.11+. To compile just run `go get -u github.com/gesquive/cig` and the executable should be built for you automatically in your `$GOPATH`. This project uses go mods, so you might need to set `GO111MODULE=on` in order for `go get` to complete properly.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/cig/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/cig`.

### Homebrew
This app is also avalable from this [homebrew tap](https://github.com/gesquive/homebrew-tap). Just install the tap and then the app will be available.
```shell
$ brew tap gesquive/tap
$ brew install cig
```

## Usage

```console
Generate certificate summary information for PEM certificates and output (in mkcert.org format)

Usage:
  cig [flags] <cert_path> [<cert_path>...]

Flags:
      --version   Display the version info and exit
```
Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug     Include debug statements in output
```

You can also pipe in certificates in addition to specified files on the command line:

```console
$ cig pki-chain.pem
$ cat massive-cert-store.pem | cig 
$ app-with-output | cig
```

## Documentation

This documentation can be found at github.com/gesquive/cig

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
