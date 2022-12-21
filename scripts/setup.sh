#!/usr/bin/bash

set -e

verify_and_install() {
    if command -v "${1}" > /dev/null; then
        return
    fi

    if [ "$(uname)" != "Darwin" ]; then
        echo "Please install ${1} and re-run the script"
        exit 1
    fi

    echo "${1} is not installed, do you wish to install ${1} [y/n]"
    read -r OK <&1
    [ "$OK" != "y" ] && [ "$OK" != "Y" ] && exit 1
    brew install "${1}"
}

verify_and_install mage
verify_and_install yarn
