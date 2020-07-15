#!/bin/bash

set -euo pipefail

# Get the dir the script lives in not just the one it's exec'd from 
SCRIPT_DIR=""
if [ "$(uname)" == 'Darwin' ]; then

    if [ -z "$(which greadlink)" ]; then
        echo "Install brew then run \`brew install coreutils\`"
        exit 1
    fi

    SCRIPT_DIR=$(dirname $(greadlink -f $0))
else
    SCRIPT_DIR=$(dirname $(readlink -f $0))
fi

# needs this privilege to change attribute on a file 
docker run --cap-add LINUX_IMMUTABLE -v $SCRIPT_DIR/../:/go/src/indelible -it $@
