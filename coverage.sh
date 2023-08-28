#!/bin/bash

# This script is used to generate coverage report for the project.
# It will generate a coverage report for the project and open it in the browser.

# Generate coverage report
go test -cover -coverprofile=coverage.out ./...

# Remove mock files from coverage report
sed -i '/_mock.go/d' coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open coverage report in browser
open coverage.html
