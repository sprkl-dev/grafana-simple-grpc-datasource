#!/usr/bin/bash

SCRIPTS_DIR=$(realpath "$(dirname "${BASH_SOURCE[0]}")")
ROOT_DIR=$(realpath "$(dirname "${SCRIPTS_DIR}")")

echo "[+] Setup environment"
"${SCRIPTS_DIR}/setup.sh"

echo "[+] Build front"
(cd "${ROOT_DIR}" && yarn install && yarn build)

echo "[+] Build back"
(cd "${ROOT_DIR}" && mage)
