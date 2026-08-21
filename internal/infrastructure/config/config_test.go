package config

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestUserHomeDirUsesInvokingUserUnderSudo(t *testing.T) {
	home := t.TempDir()
	fakeSudo(t, &user.User{Username: "developer", Uid: "501", Gid: "20", HomeDir: home})
	t.Setenv("HOME", t.TempDir())

	got, err := UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if got != home {
		t.Fatalf("UserHomeDir() = %q, want %q", got, home)
	}
}

func TestUserHomeDirIgnoresSpoofedSudoUserWithoutRootPrivileges(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "developer")
	originalEffectiveUID := effectiveUID
	effectiveUID = func() int { return 501 }
	t.Cleanup(func() { effectiveUID = originalEffectiveUID })

	got, err := UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if got != home {
		t.Fatalf("UserHomeDir() = %q, want %q", got, home)
	}
}

func TestUserHomeDirRejectsUnknownInvokingUser(t *testing.T) {
	t.Setenv("SUDO_USER", "missing")
	originalEffectiveUID := effectiveUID
	originalLookupUser := lookupUser
	effectiveUID = func() int { return 0 }
	lookupUser = func(string) (*user.User, error) { return nil, user.UnknownUserError("missing") }
	t.Cleanup(func() {
		effectiveUID = originalEffectiveUID
		lookupUser = originalLookupUser
	})
	if _, err := UserHomeDir(); err == nil {
		t.Fatal("UserHomeDir() unexpectedly succeeded")
	}
}

func TestSaveOwnsNewConfigArtifactsForInvokingUser(t *testing.T) {
	home := t.TempDir()
	fakeSudo(t, &user.User{Username: "developer", Uid: "501", Gid: "20", HomeDir: home})
	var owned []string
	originalChangeOwner := changeOwner
	changeOwner = func(path string, uid, gid int) error {
		if uid != 501 || gid != 20 {
			t.Fatalf("ownership = %d:%d", uid, gid)
		}
		owned = append(owned, path)
		return nil
	}
	t.Cleanup(func() { changeOwner = originalChangeOwner })

	if err := Save(Config{StoragePath: filepath.Join(home, "mdev")}); err != nil {
		t.Fatal(err)
	}
	if len(owned) != 2 {
		t.Fatalf("owned paths = %#v, want config directory and temporary config file", owned)
	}
	if owned[0] != filepath.Join(home, ".mdev") {
		t.Fatalf("first owned path = %q", owned[0])
	}
	info, err := os.Stat(filepath.Join(home, ".mdev", "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() {
		t.Fatal("config is not a regular file")
	}
}

func TestOwnPathReportsOwnershipFailure(t *testing.T) {
	home := t.TempDir()
	fakeSudo(t, &user.User{Username: "developer", Uid: "501", Gid: "20", HomeDir: home})
	originalChangeOwner := changeOwner
	changeOwner = func(string, int, int) error { return os.ErrPermission }
	t.Cleanup(func() { changeOwner = originalChangeOwner })

	err := OwnPathForInvokingUser(filepath.Join(home, "mdev"))
	if err == nil || !errors.Is(err, os.ErrPermission) {
		t.Fatalf("OwnPathForInvokingUser() error = %v", err)
	}
}

func fakeSudo(t *testing.T, account *user.User) {
	t.Helper()
	t.Setenv("SUDO_USER", account.Username)
	originalEffectiveUID := effectiveUID
	originalLookupUser := lookupUser
	effectiveUID = func() int { return 0 }
	lookupUser = func(name string) (*user.User, error) {
		if name != account.Username {
			return nil, user.UnknownUserError(name)
		}
		return account, nil
	}
	t.Cleanup(func() {
		effectiveUID = originalEffectiveUID
		lookupUser = originalLookupUser
	})
}
