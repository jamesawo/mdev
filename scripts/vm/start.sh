#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
. "$SCRIPT_DIR/_common.sh"

UTMCTL=$(find_utmctl)

if ! status=$("$UTMCTL" status "$VM_NAME" 2>&1); then
    printf '%s\n' "error: unable to query UTM VM '$VM_NAME': $status" >&2
    exit 1
fi

if [ "$status" = "started" ]; then
    printf '%s\n' "UTM VM '$VM_NAME' is already running."
    exit 0
fi

printf '%s\n' "Starting UTM VM '$VM_NAME'..."
if ! "$UTMCTL" start --hide "$VM_NAME"; then
    printf '%s\n' "error: unable to start UTM VM '$VM_NAME'" >&2
    exit 1
fi

printf '%s\n' "UTM VM '$VM_NAME' started."
