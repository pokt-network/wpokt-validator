#!/bin/bash

# This script is used to generate coverage report for the project.
# It will generate a coverage report for the project and open it in the browser.

# Generate coverage report
go test -cover -coverpkg=github.com/dan13ram/wpokt-validator/pokt/util,github.com/dan13ram/wpokt-validator/eth/util -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Open coverage report in browser
open coverage.html
