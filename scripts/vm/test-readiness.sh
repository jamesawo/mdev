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

remote_root="/tmp/mdev-readiness-e2e.$$"
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
fake_bin="$TEST_ROOT/bin"
storage_parent="$TEST_ROOT/storage parent"
mkdir -p "$home" "$fake_bin" "$storage_parent"

cat >"$fake_bin/bash" <<'EOF'
#!/bin/sh
case "$*" in
    *BASH_VERSINFO*) printf 5 ;;
    *) exec /bin/bash "$@" ;;
esac
EOF
cat >"$fake_bin/brew" <<'EOF'
#!/bin/sh
exit 0
EOF
cat >"$fake_bin/xcode-select" <<'EOF'
#!/bin/sh
test "$1" = "-p" && exit 0
exit 2
EOF
cat >"$fake_bin/curl" <<'EOF'
#!/bin/sh
cat <<'INSTALLER'
#!/bin/sh
mkdir -p "$HOME/.sdkman/bin"
cat >"$HOME/.sdkman/bin/sdkman-init.sh" <<'INIT'
export JAVA_HOME="$HOME/.sdkman/candidates/java/current"
sdk() { return 0; }
INIT
cat >"$MDEV_FAKE_BIN/java" <<'JAVA'
#!/bin/sh
exit 0
JAVA
chmod 0755 "$MDEV_FAKE_BIN/java"
INSTALLER
EOF
cat >"$fake_bin/gradle" <<'EOF'
#!/bin/sh
exit 0
EOF
chmod 0755 "$fake_bin"/*
PATH="$fake_bin:/usr/bin:/bin:/usr/sbin:/sbin"
export PATH
MDEV_FAKE_BIN="$fake_bin"
export MDEV_FAKE_BIN

HOME="$home" SUDO_USER= "$TEST_ROOT/mdev" setup --storage-path "$storage_parent" >"$TEST_ROOT/setup.out"
grep -q 'checking bash' "$TEST_ROOT/setup.out"
grep -q 'bash ready' "$TEST_ROOT/setup.out"
grep -q 'mdev is ready' "$TEST_ROOT/setup.out"
test -f "$home/.mdev/config.yaml"

HOME="$home" SUDO_USER= "$TEST_ROOT/mdev" setup >"$TEST_ROOT/setup-again.out"
grep -q 'setup is complete' "$TEST_ROOT/setup-again.out"

HOME="$home" SUDO_USER= "$TEST_ROOT/mdev" list >"$TEST_ROOT/list.out"
grep -q '^system tools$' "$TEST_ROOT/list.out"
grep -q '^tools$' "$TEST_ROOT/list.out"

HOME="$home" SUDO_USER= "$TEST_ROOT/mdev" install gradle --yes >"$TEST_ROOT/install.out"
grep -q '^Installing sdkman\.\.\.$' "$TEST_ROOT/install.out"
grep -q '^Installing java\.\.\.$' "$TEST_ROOT/install.out"
grep -q '^Installing gradle\.\.\.$' "$TEST_ROOT/install.out"
grep -q 'gradle installed$' "$TEST_ROOT/install.out"

HOME="$home" SUDO_USER= "$TEST_ROOT/mdev" install gradle --yes >"$TEST_ROOT/retry.out"
grep -q '^gradle is already installed\.$' "$TEST_ROOT/retry.out"
if grep -q '^Installing gradle' "$TEST_ROOT/retry.out"; then
    echo "retry reinstalled completed tool" >&2
    exit 1
fi

echo "readiness fixture e2e passed"
REMOTE
