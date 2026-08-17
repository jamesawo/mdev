#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
. "$SCRIPT_DIR/_common.sh"

timeout=${MDEV_VM_WAIT_TIMEOUT:-120}
interval=${MDEV_VM_WAIT_INTERVAL:-2}
elapsed=0

printf '%s\n' "Waiting for SSH alias '$SSH_ALIAS'..."

while [ "$elapsed" -lt "$timeout" ]; do
    if ssh -o BatchMode=yes -o ConnectTimeout=5 "$SSH_ALIAS" true >/dev/null 2>&1; then
        printf '%s\n' "SSH alias '$SSH_ALIAS' is ready."
        exit 0
    fi

    sleep "$interval"
    elapsed=$((elapsed + interval))
    printf '%s\n' "Still waiting for SSH ($elapsed/$timeout seconds)..."
done

printf '%s\n' "error: SSH alias '$SSH_ALIAS' was not reachable after $timeout seconds" >&2
exit 1
