#!/bin/bash

go test -v -tags all_test $(go list ./... | grep -v vendor) --count 1 -covermode=atomic -timeout 60s
