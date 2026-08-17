#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

"$SCRIPT_DIR/start.sh"
"$SCRIPT_DIR/wait.sh"

exec ssh -t "${MDEV_VM_SSH_ALIAS:-mdev-vm}"
