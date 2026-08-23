package shell

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSDKMANCandidateStatusUsesManagedCandidateInsteadOfPath(t *testing.T) {
	home, fakeBin := prepareSDKMANTest(t)
	global := filepath.Join(fakeBin, "gradle")
	if err := os.WriteFile(global, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}

	installed, err := SDKMANCandidateInstallationStatus(context.Background(), "gradle", "gradle", "--version")
	if err != nil {
		t.Fatal(err)
	}
	if installed {
		t.Fatal("unrelated PATH gradle was treated as SDKMAN-managed")
	}

	managed := filepath.Join(home, ".sdkman", "candidates", "gradle", "current", "bin", "gradle")
	if err := os.MkdirAll(filepath.Dir(managed), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(managed, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	installed, err = SDKMANCandidateInstallationStatus(context.Background(), "gradle", "gradle", "--version")
	if err != nil || !installed {
		t.Fatalf("managed status = %v, %v", installed, err)
	}
}

func TestInstallSDKMANCandidateUsesManagedSDKMAN(t *testing.T) {
	home, _ := prepareSDKMANTest(t)
	logPath := filepath.Join(t.TempDir(), "sdk.log")
	t.Setenv("MDEV_TEST_SDK_LOG", logPath)
	initPath := filepath.Join(home, ".sdkman", "bin", "sdkman-init.sh")
	if err := os.WriteFile(initPath, []byte("sdk() { printf '%s\\n' \"$*\" >> \"$MDEV_TEST_SDK_LOG\"; }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := InstallSDKMANCandidateContext(context.Background(), "gradle"); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "install gradle\n" {
		t.Fatalf("sdk call = %q", output)
	}
}

func TestUninstallSDKMANCandidateClearsCurrentBeforeRemoval(t *testing.T) {
	home, _ := prepareSDKMANTest(t)
	logPath := filepath.Join(t.TempDir(), "sdk.log")
	t.Setenv("MDEV_TEST_SDK_LOG", logPath)
	initPath := filepath.Join(home, ".sdkman", "bin", "sdkman-init.sh")
	initScript := "sdk() { printf '%s current=%s\\n' \"$*\" \"$(test -L \"$HOME/.sdkman/candidates/maven/current\" && echo yes || echo no)\" >> \"$MDEV_TEST_SDK_LOG\"; }\n"
	if err := os.WriteFile(initPath, []byte(initScript), 0644); err != nil {
		t.Fatal(err)
	}
	candidateRoot := filepath.Join(home, ".sdkman", "candidates", "maven")
	versionDir := filepath.Join(candidateRoot, "3.9.16", "bin")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("3.9.16", filepath.Join(candidateRoot, "current")); err != nil {
		t.Fatal(err)
	}

	if err := UninstallSDKMANCandidate("maven"); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "uninstall maven 3.9.16 current=no\n" {
		t.Fatalf("sdk call = %q", output)
	}
	if _, err := os.Lstat(filepath.Join(candidateRoot, "current")); !os.IsNotExist(err) {
		t.Fatalf("current selection still exists: %v", err)
	}
}

func TestUninstallSDKMANCandidateRestoresCurrentAfterFailure(t *testing.T) {
	home, _ := prepareSDKMANTest(t)
	initPath := filepath.Join(home, ".sdkman", "bin", "sdkman-init.sh")
	if err := os.WriteFile(initPath, []byte("sdk() { return 1; }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	candidateRoot := filepath.Join(home, ".sdkman", "candidates", "maven")
	if err := os.MkdirAll(filepath.Join(candidateRoot, "3.9.16", "bin"), 0755); err != nil {
		t.Fatal(err)
	}
	currentLink := filepath.Join(candidateRoot, "current")
	if err := os.Symlink("3.9.16", currentLink); err != nil {
		t.Fatal(err)
	}

	if err := UninstallSDKMANCandidate("maven"); err == nil {
		t.Fatal("uninstall unexpectedly succeeded")
	}
	target, err := os.Readlink(currentLink)
	if err != nil {
		t.Fatal(err)
	}
	if target != "3.9.16" {
		t.Fatalf("restored current target = %q", target)
	}
}

func prepareSDKMANTest(t *testing.T) (string, string) {
	t.Helper()
	home := t.TempDir()
	fakeBin := t.TempDir()
	bashPath := filepath.Join(fakeBin, "bash")
	bashScript := `#!/bin/sh
case "$*" in
  *BASH_VERSINFO*) printf 5 ;;
  *) exec /bin/bash "$@" ;;
esac
`
	if err := os.WriteFile(bashPath, []byte(bashScript), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, ".sdkman", "bin"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".sdkman", "bin", "sdkman-init.sh"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("PATH", fakeBin)
	return home, fakeBin
}
