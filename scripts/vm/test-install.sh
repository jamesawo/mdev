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

remote_root="/tmp/mdev-install-e2e.$$"
cleanup() {
    ssh "$SSH_ALIAS" "rm -rf '$remote_root'" >/dev/null 2>&1 || true
}
trap cleanup EXIT HUP INT TERM

ssh "$SSH_ALIAS" "mkdir -m 0700 '$remote_root'"
scp -q "$BINARY" "$SSH_ALIAS:$remote_root/mdev"

ssh "$SSH_ALIAS" "TEST_ROOT='$remote_root' sh -s" <<'REMOTE'
set -eu
chmod 0755 "$TEST_ROOT/mdev"
home="$TEST_ROOT/home"
storage="$TEST_ROOT/storage"
fake_bin="$TEST_ROOT/bin"
mkdir -p "$home/.mdev" "$storage" "$fake_bin"
printf 'storage_path: %s\n' "$storage" >"$home/.mdev/config.yaml"

cat >"$fake_bin/brew" <<'EOF'
#!/bin/sh
exit 0
EOF
cat >"$fake_bin/xcode-select" <<'EOF'
#!/bin/sh
test "$1" = "-p"
EOF
cat >"$fake_bin/podman" <<'EOF'
#!/bin/sh
case "$*" in
    "--version") exit 0 ;;
    "machine init --image-path "*) touch "$HOME/.podman-machine-initialized" ;;
    "machine inspect") test -f "$HOME/.podman-machine-initialized" ;;
    *) exit 2 ;;
esac
EOF
chmod 0755 "$fake_bin/brew" "$fake_bin/podman" "$fake_bin/xcode-select"
PATH="$fake_bin:/usr/bin:/bin:/usr/sbin:/sbin"
export PATH

HOME="$home" "$TEST_ROOT/mdev" install podman --yes >"$TEST_ROOT/first.out"
grep -q '^Install plan$' "$TEST_ROOT/first.out"
grep -q '^Installing podman\.\.\.$' "$TEST_ROOT/first.out"
grep -q '^Configuring podman\.\.\.$' "$TEST_ROOT/first.out"
grep -q '^Verifying podman\.\.\.$' "$TEST_ROOT/first.out"
grep -q 'podman installed$' "$TEST_ROOT/first.out"

HOME="$home" "$TEST_ROOT/mdev" install podman --yes >"$TEST_ROOT/retry.out"
grep -q '^podman is already installed\.$' "$TEST_ROOT/retry.out"
grep -q '^Uninstall: mdev uninstall podman$' "$TEST_ROOT/retry.out"
if grep -q '^Installing podman' "$TEST_ROOT/retry.out"; then
    echo "retry reinstalled completed tool" >&2
    exit 1
fi

if HOME="$home" "$TEST_ROOT/mdev" install missing --yes >"$TEST_ROOT/unknown.out" 2>&1; then
    echo "unknown tool unexpectedly succeeded" >&2
    exit 1
fi
grep -q 'Run `mdev list`' "$TEST_ROOT/unknown.out"

if HOME="$home" "$TEST_ROOT/mdev" install podman --all --yes >"$TEST_ROOT/conflict.out" 2>&1; then
    echo "conflicting invocation unexpectedly succeeded" >&2
    exit 1
fi
grep -q 'cannot use --all' "$TEST_ROOT/conflict.out"

echo "install e2e passed"
REMOTE
