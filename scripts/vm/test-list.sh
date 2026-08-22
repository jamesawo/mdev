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

remote_root="/tmp/mdev-list-e2e.$$"
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

test_home="$TEST_ROOT/home"
storage="$TEST_ROOT/storage"
mkdir -p "$test_home/.mdev" "$storage"
cat >"$test_home/.mdev/config.yaml" <<EOF
storage_path: $storage
EOF

env HOME="$test_home" SUDO_USER= "$TEST_ROOT/mdev" list >"$TEST_ROOT/list.out"
env HOME="$test_home" SUDO_USER= "$TEST_ROOT/mdev" list --json >"$TEST_ROOT/list.json"

grep -q '^system tools$' "$TEST_ROOT/list.out"
grep -q '^tools$' "$TEST_ROOT/list.out"
grep -Eq '^  [✓○?] curl + (installed|missing|unknown)$' "$TEST_ROOT/list.out"
grep -Eq '^  [✓○?] git + (installed|missing|unknown)$' "$TEST_ROOT/list.out"
grep -Eq '^  [✓○?] gradle + (installed|missing|unknown)$' "$TEST_ROOT/list.out"
grep -Eq '^  [✓○?] java + (installed|missing|unknown)$' "$TEST_ROOT/list.out"

system_names=$(awk '
    /^system tools$/ { section = 1; next }
    /^tools$/ { section = 0 }
    section && /^  / { print $2 }
' "$TEST_ROOT/list.out")
tool_names=$(awk '
    /^tools$/ { section = 1; next }
    section && /^  / { print $2 }
' "$TEST_ROOT/list.out")
test "$system_names" = "$(printf '%s\n' "$system_names" | LC_ALL=C sort -f)"
test "$tool_names" = "$(printf '%s\n' "$tool_names" | LC_ALL=C sort -f)"

/usr/bin/jq empty "$TEST_ROOT/list.json"
grep -q '^{' "$TEST_ROOT/list.json"
grep -q '"system_tools":\[' "$TEST_ROOT/list.json"
grep -q '"tools":\[' "$TEST_ROOT/list.json"
grep -q '"name":"curl","status":"installed"' "$TEST_ROOT/list.json"
if grep -Eq 'system tools|✓|○|could not determine' "$TEST_ROOT/list.json"; then
    echo "JSON output contains human-oriented content" >&2
    exit 1
fi

env HOME="$test_home" SUDO_USER= "$TEST_ROOT/mdev" list --json >"$TEST_ROOT/list-second.json"
cmp "$TEST_ROOT/list.json" "$TEST_ROOT/list-second.json"

echo "list e2e passed"
REMOTE
