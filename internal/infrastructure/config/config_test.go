package config

import (
	"os/user"
	"testing"
)

func TestUserHomeDirUsesInvokingUserUnderSudo(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	if current.Username == "root" {
		t.Skip("test requires a non-root invoking user")
	}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("SUDO_USER", current.Username)
	home, err := UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if home != current.HomeDir {
		t.Fatalf("UserHomeDir() = %q, want %q", home, current.HomeDir)
	}
}

func TestUserHomeDirRejectsUnknownInvokingUser(t *testing.T) {
	t.Setenv("SUDO_USER", "mdev-user-that-does-not-exist")
	if _, err := UserHomeDir(); err == nil {
		t.Fatal("UserHomeDir() unexpectedly succeeded")
	}
}
