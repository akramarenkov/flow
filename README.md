# Flow

[![Go Reference](https://pkg.go.dev/badge/github.com/akramarenkov/flow.svg)](https://pkg.go.dev/github.com/akramarenkov/flow)
[![Go Report Card](https://goreportcard.com/badge/github.com/akramarenkov/flow)](https://goreportcard.com/report/github.com/akramarenkov/flow)
[![Coverage Status](https://coveralls.io/repos/github/akramarenkov/flow/badge.svg)](https://coveralls.io/github/akramarenkov/flow)

## Purpose

Library that allows you to manage the flow of data between Go channels

## Implemented disciplines

* **join** - accumulates data items from an input channel into a slice and write that slice to an output channel when the maximum slice size or timeout for its accumulation is reached. See [README](join/README.md)

* **limit** - limits the speed of passing data items from the input channel to the output channel. See [README](limit/README.md)

* **priority** - distributes data items between handlers in quantity corresponding to the priority of the data items. See [README](priority/README.md)
