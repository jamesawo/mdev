#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

if [ "$#" -eq 0 ]; then
    printf '%s\n' "usage: $0 <command> [arguments...]" >&2
    exit 2
fi

"$SCRIPT_DIR/start.sh"
"$SCRIPT_DIR/wait.sh"

remote_command='PATH=/usr/local/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin; export PATH;'
for argument do
    escaped=$(printf '%s' "$argument" | sed "s/'/'\\\\''/g")
    remote_command="${remote_command} '${escaped}'"
done

exec ssh "${MDEV_VM_SSH_ALIAS:-mdev-vm}" "$remote_command"
