# cig
[![Travis CI](https://img.shields.io/travis/gesquive/cig/master.svg?style=flat-square)](https://travis-ci.org/gesquive/cig)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/cig/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gesquive/cig)

Formats PEM formatted certificate stores in a nice human readable format

## Installing

### Compile
This project requires go 1.10+ to compile. Just run `go get -u github.com/gesquive/cig` and the executable should be built for you automatically in your `$GOPATH`.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/cig/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/cig`.

## Usage

```console
Generate certificate information for multiple PEM formatted certificates

Usage:
  cig [flags] <cert_path> [<cert_path>...]

Flags:
  -h, --help      help for cig
  -v, --verbose   Print logs to stdout instead of file
      --version   Display the version number and exit
```
Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug     Include debug statements in log output
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
