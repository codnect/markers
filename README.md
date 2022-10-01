
![Marker Logo](https://user-images.githubusercontent.com/5354910/189712369-f647731b-dc34-405a-a7df-ca23f7ea1025.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/procyon-projects/marker)](https://goreportcard.com/report/github.com/procyon-projects/marker)
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/procyon-projects/marker/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/procyon-projects/marker/tree/main)
[![codecov](https://codecov.io/gh/procyon-projects/marker/branch/main/graph/badge.svg?token=OREV0YI8VU)](https://codecov.io/gh/procyon-projects/marker)

Marker project is inspired by [Kubernetes Markers](https://book.kubebuilder.io/reference/markers.html).

`It aims to make markers available for everyone.`

**Markers** are single-comments that start with a plus, followed by a marker name, optionally followed some marker parameters,
**which are used to generate or verify code but add no logic at runtime.**

## Installation
To Install Marker CLI quickly, follow the installation instructions.

1. You first need Go installed (version 1.18+ is required), then you can use the below Go command to install Marker CLI.

    `$ go get -u github.com/procyon-projects/marker/...`
2. Verify that you've installed Marker CLI by typing the following command.

   `$ marker version`
3. Confirm that the command prints the installed version of Marker CLI.

Type the following command to display usage information for the Marker CLI.
`$ marker help`

# License
Marker is released under Apache-2.0 License.