# strainscrape

### Table of Contents

1. [Install](#install)
2. [Commands](#commands)

A command line interface to that scrapes cannabis strain information from [CannaConnection](annaconnection.com).

## Install

### Homebrew

Use [Homebrew](https://brew.sh) to install:

```bash
$ brew tap droxey/strainscrape
$ brew install strainscape
```

### Golang

```bash
$ go get github.com/droxey/strainscrape
```

Note: If you use modules (or if `GO111MODULE=on` in your environment), `go get` will not install packages "globally". It will add them to your project's `go.mod` file instead. As of Go 1.11.1, setting `GO111MODULE=off` works to circumvent this behavior:

```bash
$ GO111MODULE=off; go get github.com/droxey/strainscrape
```

## Commands

### List Strains

Returns a list of all strains.

```bash
$ strains list
```

Save the list in a JSON file named `strains.json`.

```bash
$ strains save --filename="strains.json"
```

### Find a Strain

Find a strain by name.

```bash
$ strains find "Strain Name"
```
