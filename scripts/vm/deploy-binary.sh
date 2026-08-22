#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPOSITORY_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/../.." && pwd)
BINARY=${MDEV_VM_BINARY:-$REPOSITORY_ROOT/dist/mdev}
. "$SCRIPT_DIR/_common.sh"

if [ ! -x "$BINARY" ]; then
    printf '%s\n' "error: build the macOS ARM64 binary at $BINARY first" >&2
    exit 1
fi

remote_dir='.local/share/mdev/bin'
remote_binary="$remote_dir/mdev"
remote_temporary="$remote_dir/.mdev.$$"

cleanup() {
    ssh "$SSH_ALIAS" "rm -f '$remote_temporary'" >/dev/null 2>&1 || true
}
trap cleanup EXIT HUP INT TERM

ssh "$SSH_ALIAS" "mkdir -p '$remote_dir'"
scp -q "$BINARY" "$SSH_ALIAS:$remote_temporary"
ssh "$SSH_ALIAS" "chmod 0755 '$remote_temporary' && mv -f '$remote_temporary' '$remote_binary'"
trap - EXIT HUP INT TERM

printf '%s\n' "Deployed active-checkout binary to $SSH_ALIAS:~/$remote_binary"
