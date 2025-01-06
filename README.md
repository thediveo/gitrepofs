# gitrepofs

[![PkgGoDev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/gitrepofs)
[![GitHub](https://img.shields.io/github/license/thediveo/gitrepofs)](https://img.shields.io/github/license/thediveo/gitrepofs)
![build and test](https://github.com/thediveo/gitrepofs/actions/workflows/buildandtest.yaml/badge.svg?branch=master)
![Coverage](https://img.shields.io/badge/Coverage-93.9%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/gitrepofs)](https://goreportcard.com/report/github.com/thediveo/gitrepofs)

A Go [fs.FS](https://pkg.go.dev/io/fs#FS) _git repository file system_ to easily
access a repository at a specific tag (or other git reference).

For devcontainer instructions, please see the [section "DevContainer"
below](#devcontainer).

## Usage

The main envisioned use case is for `go generate`-based updates from upstream
repositories to fetch the latest C definitions without the need to integrate
upstream C libraries.

```golang
remoteURL := "https://gohub.org/froozle/baduzle"
latest, latestref, err := version.LatestReleaseTag(context.Background(),
    remoteURL, version.SemverMatcher)
gfs, err := NewForRevision(context.Background(), remoteURL, latestref)
contents, err := fs.ReadFile(gfs, "some/useful/file.h")
```

## Other Implementations

- [@hairyhenderson/go-fsimpl](https://github.com/hairyhenderson/go-fsimpl)
  includes a git file system implementation (beside others). It uses the neat
  trick of cloning a repository into a local work tree in memory (instead of on
  "disk") and then serves from this memory-based file system. This design
  requires files to be present twice: once in the in-memory cloned repository
  and another time in the memory-based file system. Actively maintained at the
  time of this writing, as well as equiped with lots of unit tests.

- [@ear7h/go-git-fs](https://github.com/ear7h/go-git-fs) serves directly from an
  in-memory git repository clone without a work tree. It doesn't come with any
  unit tests and hasn't been maintained since May 2021. The code does the
  [fs.ValidPath](https://pkg.go.dev/io/fs#ValidPath) checks but then adds an
  unnecessary [path.Clean](https://pkg.go.dev/path#Clean) because `fs.ValidPath`
  blocks all the things that `path.Clean` is supposed to sanitize.

- out of competition: [@posener/gitfs](https://github.com/posener/gitfs) tackles
  [http.FileSystem](https://pkg.go.dev/net/http#FileSystem) instead.

## DevContainer

> [!CAUTION]
>
> Do **not** use VSCode's "~~Dev Containers: Clone Repository in Container
> Volume~~" command, as it is utterly broken by design, ignoring
> `.devcontainer/devcontainer.json`.

1. `git clone https://github.com/thediveo/enumflag`
2. in VSCode: Ctrl+Shift+P, "Dev Containers: Open Workspace in Container..."
3. select `enumflag.code-workspace` and off you go...

## Supported Go Versions

`gitrepofs` supports versions of Go that are noted by the Go release policy,
that is, major versions _N_ and _N_-1 (where _N_ is the current major version).

## Copyright and License

`gitrepofs` is Copyright 2023, 2025 Harald Albrecht, and licensed under the
Apache License, Version 2.0.
