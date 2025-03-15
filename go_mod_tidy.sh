#!/bin/bash

gci write --skip-generated -s default sandwich
gofumpt -d -e -extra -l -w sandwich
go mod tidy
