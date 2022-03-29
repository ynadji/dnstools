# `dnstools`

A collection of tools to filter and manipulate domain names on standard input.

## Installation

Make sure `go` is installed through your distro. You can install all the tools
with:

```
$ go install github.com/ynadji/dnstools/...@latest
```

## Usage

```
¡ nld --help
NAME:
   nld - nld < domains.txt > nlds.txt

USAGE:
   nld [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --numLevels value, -n value  Domain level to extract (default: 2)
   --ignoreErrors               Silently ignore errors (default: false)
   --usePrivateSuffixes         Treat known private suffixes as TLDs (default: false)
   --ignoreNonIcann             Ignore domains with TLDs not known by the PSL (default: false)
   --exactOnly                  Only keep domains if exactly matches label level (default: false)
   --help, -h                   show help (default: false)

```

```
¡ dclean --help
NAME:
   dclean - dclean [subcommand] < domains.txt > processed-domains.txt

USAGE:
   dclean [global options] command [command options] [arguments...]

COMMANDS:
   valid, v         Filter to domains that are syntactically valid.
   registerable, r  Filter to domains that could be registered and used. Excludes TLDs.
   ascii, a         Map input domains to ASCII (IDN->punycode).
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)

```
