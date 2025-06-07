#!/bin/bash

delete_lines_matching() {
  local pattern="$1"
  local file="$2"

  if [[ -z "$pattern" || -z "$file" ]]; then
    echo "Usage: delete_lines_matching <pattern> <file>"
    return 1
  fi

  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "/$pattern/d" "$file"
  else
    sed -i "/$pattern/d" "$file"
  fi
}

# This script is used to generate coverage report for the project.
# It will generate a coverage report for the project and open it in the browser.

# Generate coverage report
go test -cover -coverprofile=coverage.out ./...

# Remove mock files from coverage report
delete_lines_matching 'mock' coverage.out

# Remove script files from coverage report
delete_lines_matching 'scripts' coverage.out

# Remove autogen files from coverage report
delete_lines_matching 'autogen' coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open coverage report in browser
open coverage.html
