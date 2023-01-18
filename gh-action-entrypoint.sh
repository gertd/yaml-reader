#!/usr/bin/env bash
set -e

export PRE_RELEASE="$INPUT_PRE_RELEASE"

if [ -n ${GITHUB_WORKSPACE} ]; then
  git config --global --add safe.directory "$GITHUB_WORKSPACE"
fi

values() {
    ./dist/build_darwin_arm64/yaml-reader --filename ${INPUT_FILENAME}
}

while read -r value; do
    printf "%s\n" ${value} >> $GITHUB_OUTPUT
done < <(values)
