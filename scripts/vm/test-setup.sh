#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPOSITORY_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/../.." && pwd)
BINARY=${MDEV_E2E_BINARY:-$REPOSITORY_ROOT/dist/mdev}
. "$SCRIPT_DIR/_common.sh"

if [ ! -x "$BINARY" ]; then
    printf '%s\n' "error: build the macOS ARM64 binary at $BINARY first" >&2
    exit 1
fi

"$SCRIPT_DIR/start.sh"
"$SCRIPT_DIR/wait.sh"

remote_root="/tmp/mdev-setup-e2e.$$"
cleanup() {
    ssh "$SSH_ALIAS" "PATH=/usr/bin:/bin; pkill -f '$remote_root/mdev setup' 2>/dev/null || true; rm -rf '$remote_root'" >/dev/null 2>&1 || true
}
trap cleanup EXIT HUP INT TERM

ssh "$SSH_ALIAS" "mkdir -m 0700 '$remote_root'"
scp -q "$BINARY" "$SSH_ALIAS:$remote_root/mdev"

ssh "$SSH_ALIAS" "TEST_ROOT='$remote_root' sh -s" <<'REMOTE'
set -eu
PATH=/usr/local/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin
export PATH
chmod 0755 "$TEST_ROOT/mdev"
canonical_root=$(CDPATH= cd -- "$TEST_ROOT" && pwd -P)

interactive_home="$TEST_ROOT/interactive-home"
mkdir -p "$interactive_home"
INTERACTIVE_HOME=$interactive_home
MDEV_BINARY=$TEST_ROOT/mdev
OUTPUT=$TEST_ROOT/interactive.out
export INTERACTIVE_HOME MDEV_BINARY OUTPUT
expect <<'EXPECT'
set timeout 10
log_file -noappend $env(OUTPUT)
spawn env HOME=$env(INTERACTIVE_HOME) SUDO_USER= $env(MDEV_BINARY) setup
expect "where should mdev store data?"
send "\r"
expect {
    "this location is a symlink." {
        expect "continue"
        send "\r"
        exp_continue
    }
    "mdev is ready." {}
    timeout { exit 124 }
    eof { exit 1 }
}
expect eof
catch wait result
exit [lindex $result 3]
EXPECT
grep -q "mdev is ready" "$TEST_ROOT/interactive.out"
grep -q "storage_path: $canonical_root/interactive-home/mdev" "$interactive_home/.mdev/config.yaml"
test -d "$canonical_root/interactive-home/mdev"

noninteractive_home="$TEST_ROOT/noninteractive-home"
custom_parent="$TEST_ROOT/custom parent"
mkdir -p "$noninteractive_home"
env HOME="$noninteractive_home" SUDO_USER= "$TEST_ROOT/mdev" setup --storage-path "$custom_parent" >"$TEST_ROOT/noninteractive.out"
grep -q "mdev is ready" "$TEST_ROOT/noninteractive.out"
grep -q "storage_path: $canonical_root/custom parent/mdev" "$noninteractive_home/.mdev/config.yaml"
test -d "$canonical_root/custom parent/mdev"

echo "setup e2e passed"
REMOTE
