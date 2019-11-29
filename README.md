# insights-operator-web-ui
Web UI for insights operator instrumentation service

[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-operator-web-ui)](https://goreportcard.com/report/github.com/RedHatInsights/insights-operator-web-ui) [![Build Status](https://travis-ci.org/RedHatInsights/insights-operator-web-ui.svg?branch=master)](https://travis-ci.org/RedHatInsights/insights-web-ui)

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
