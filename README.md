
![Marker Logo](https://user-images.githubusercontent.com/5354910/196059225-d7bc5236-f247-41da-bee9-3fd317ad207f.png)

# Markers

[![Go Report Card](https://goreportcard.com/badge/github.com/procyon-projects/marker)](https://goreportcard.com/report/github.com/procyon-projects/marker)
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/procyon-projects/markers/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/procyon-projects/markers/tree/main)
[![codecov](https://codecov.io/gh/procyon-projects/markers/branch/main/graph/badge.svg?token=cWW7Ek5ZvD)](https://codecov.io/gh/procyon-projects/markers)

`Marker project aims to make markers available for everyone.`

It is inspired by [Kubernetes Markers](https://book.kubebuilder.io/reference/markers.html), which help avoid boilerplate code and simplify code logic 
while working on kubernetes operators. And also it includes some code snippets from [controller-tools](https://github.com/kubernetes-sigs/controller-tools).

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
