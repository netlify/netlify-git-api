#!/bin/bash

set -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

echo "Building Linux Binary"
$DIR/toolchains/build_linux.sh
echo "Building OSX Binary"
$DIR/toolchains/build_osx.sh
echo "Building Windows Binary"
$DIR/toolchains/build_windows.sh
