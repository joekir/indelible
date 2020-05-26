#!/bin/bash

set -euo pipefail

# needs this privilege to change attribute on a file 
docker run --cap-add LINUX_IMMUTABLE -v $PWD:/go/src/indelible -it $@
