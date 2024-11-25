#!/bin/bash

if [ -z "$1" ] || [ -z "$3" ]; then
  echo "This command requires the flags --groupId and --apiKey to be provided, with valid values."
  echo "Usage: $0 --groupId <groupId> --apiKey <apiKey> --filePath <filePath> (or) --groupId <groupId> --filePath <filePath>"
  exit 1
fi

go mod tidy
go run . "$@"