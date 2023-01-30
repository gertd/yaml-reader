#!/usr/bin/env bash

set -e

# validate if values are set
[ -z "${INPUT_FILE}" ] && echo "INPUT_FILE is not set, exiting." && exit 2

if [ -n ${GITHUB_WORKSPACE} ]; then
  git config --global --add safe.directory "$GITHUB_WORKSPACE"
fi

#
# start execution block
#
e_code=0

# construct commandline arguments 
CMD="/app/yaml-reader --file ${INPUT_FILE}"

if [[ "$INPUT_JSON" == "true" ]]; then
  CMD="${CMD} --json"
fi

# execute command
eval "$CMD" || e_code=1

exit $e_code
