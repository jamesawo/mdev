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

remote_root="/tmp/mdev-version-e2e.$$"
cleanup() {
    ssh "$SSH_ALIAS" "rm -rf '$remote_root'" >/dev/null 2>&1 || true
}
trap cleanup EXIT HUP INT TERM

ssh "$SSH_ALIAS" "mkdir -m 0700 '$remote_root'"
scp -q "$BINARY" "$SSH_ALIAS:$remote_root/mdev"

ssh "$SSH_ALIAS" "TEST_ROOT='$remote_root' sh -s" <<'REMOTE'
set -eu
PATH=/usr/local/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin
export PATH
chmod 0755 "$TEST_ROOT/mdev"

test_home="$TEST_ROOT/home-without-mdev-configuration"
mkdir -p "$test_home"

env HOME="$test_home" SUDO_USER= "$TEST_ROOT/mdev" version >"$TEST_ROOT/version.out"
printf '%s\n' \
    'mdev 0.4.2' \
    'commit: unknown' \
    'built: unknown' >"$TEST_ROOT/version.expected"
cmp "$TEST_ROOT/version.expected" "$TEST_ROOT/version.out"

env HOME="$test_home" SUDO_USER= "$TEST_ROOT/mdev" version --json >"$TEST_ROOT/version.json"
printf '%s\n' '{"version":"0.4.2","commit":"unknown","built":"unknown"}' >"$TEST_ROOT/version-json.expected"
cmp "$TEST_ROOT/version-json.expected" "$TEST_ROOT/version.json"
/usr/bin/jq empty "$TEST_ROOT/version.json"

env HOME="$test_home" SUDO_USER= "$TEST_ROOT/mdev" --version >"$TEST_ROOT/root-version.out"
printf '%s\n' 'mdev 0.4.2' >"$TEST_ROOT/root-version.expected"
cmp "$TEST_ROOT/root-version.expected" "$TEST_ROOT/root-version.out"

test ! -e "$test_home/.mdev"

echo "version e2e passed"
REMOTE
