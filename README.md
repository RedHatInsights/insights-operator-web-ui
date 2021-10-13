# insights-operator-web-ui
Web UI for insights operator instrumentation service

[![GoDoc](https://godoc.org/github.com/RedHatInsights/insights-operator-web-ui?status.svg)](https://godoc.org/github.com/RedHatInsights/insights-operator-web-ui)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-operator-web-ui)](https://goreportcard.com/report/github.com/RedHatInsights/insights-operator-web-ui)
[![Build Status](https://travis-ci.org/RedHatInsights/insights-operator-web-ui.svg?branch=master)](https://travis-ci.org/RedHatInsights/insights-operator-web-ui)
[![codecov](https://codecov.io/gh/RedHatInsights/insights-operator-web-ui/branch/master/graph/badge.svg)](https://codecov.io/gh/RedHatInsights/insights-operator-web-ui)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/RedHatInsights/insights-operator-web-ui)
[![License](https://img.shields.io/badge/license-Apache-blue)](https://github.com/RedHatInsights/insights-operator-web-ui/blob/master/LICENSE)


<!-- vim-markdown-toc GFM -->

* [Description](#description)
* [How to build it](#how-to-build-it)
* [Start](#start)
* [Configuration](#configuration)
* [CI](#ci)
* [Contribution](#contribution)

<!-- vim-markdown-toc -->

## Description

A simple web-based user interface to the insights operator instrumentation service

## How to build it

Use the standard Go command:

```
go build
```

This command should create an executable file named `insights-operator-web-ui`.

## Start

Just run the executable file created by `go build`:

```
./insights-operator-web-ui
```

## Configuration

Configuration is stored in `config.toml`. ATM two options needs to be specified:

* URL to the insights operator instrumentation service
* port or full address where this tool will be available

## CI

[Travis CI](https://travis-ci.com/) is configured for this repository. Several tests and checks are started for all pull requests:

* Unit tests that use the standard tool `go test`
* `go fmt` tool to check code formatting. That tool is run with `-s` flag to perform [following transformations](https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
* `go vet` to report likely mistakes in source code, for example suspicious constructs, such as Printf calls whose arguments do not align with the format string.
* `golint` as a linter for all Go sources stored in this repository
* `gocyclo` to report all functions and methods with too high cyclomatic complexity. The cyclomatic complexity of a function is calculated according to the following rules: 1 is the base complexity of a function +1 for each 'if', 'for', 'case', '&&' or '||' Go Report Card warns on functions with cyclomatic complexity > 9
* `ineffassign` to detect and print all ineffectual assignments in Go code
* `errcheck` for checking for all unchecked errors in go programs
* `shellcheck` to perform static analysis for all shell scripts used in this repository
* `abcgo` to measure ABC metrics for Go source code and check if the metrics does not exceed specified threshold

Please note that all checks mentioned above have to pass for the change to be merged into master branch.

History of checks done by CI is available at [RedHatInsights / insights-operator-web-ui](https://travis-ci.org/RedHatInsights/insights-operator-web-ui).

## Contribution

Please look into document [CONTRIBUTING.md](CONTRIBUTING.md) that contains all information about how to contribute to this project.

Please look also at [Definition of Done](DoD.md) document with further informations.

## Package manifest

Package manifest is available at [docs/manifest.txt](docs/manifest.txt).
