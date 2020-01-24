# `nld`

## Installation

Make sure `go` is installed and:

```
$ go get github.com/ynadji/nld
$ go install github.com/ynadji/nld
```

## Usage

```
Usage:
  nld [flags]

Flags:
      --exactOnly            Only keep domains if exactly matches label level
  -h, --help                 help for nld
      --ignoreErrors         Silently ignore errors
      --ignoreNonIcann       Ignore domains with TLDs not known by the PSL
      --n int                Domain level to extract (default 2)
      --usePrivateSuffixes   Treat known private suffixes as TLDs
```
