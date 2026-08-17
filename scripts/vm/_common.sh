#!/bin/sh

VM_NAME=${MDEV_VM_NAME:-macOs-mdev}
SSH_ALIAS=${MDEV_VM_SSH_ALIAS:-mdev-vm}

find_utmctl() {
    if command -v utmctl >/dev/null 2>&1; then
        command -v utmctl
        return 0
    fi

    bundled_utmctl=/Applications/UTM.app/Contents/MacOS/utmctl
    if [ -x "$bundled_utmctl" ]; then
        printf '%s\n' "$bundled_utmctl"
        return 0
    fi

    printf '%s\n' "error: utmctl was not found; install UTM or add utmctl to PATH" >&2
    return 1
}
