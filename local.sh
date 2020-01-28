#!/bin/bash

set -e -o pipefail
trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT
cd `dirname $0`
CURRENT=`pwd`

function test
{
  go test -v -tags all_test $(go list ./... | grep -v vendor) --count 1 -covermode=atomic -timeout 60s
}

CMD=$1
shift
$CMD $*
