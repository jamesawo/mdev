package tools

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

// StatusChecker is implemented by tools whose installation check can report
// an unexpected failure separately from a missing installation.
type StatusChecker interface {
	InstallationStatus(env *environment.Environment) (bool, error)
}

// ManagedSymlinkStatus reports whether source is a symlink resolving to target.
func ManagedSymlinkStatus(source, target string) (bool, error) {
	return storage.ManagedSymlinkStatus(source, target)
}

// InstallationStatus preserves compatibility with existing Tool
// implementations while allowing built-in tools to expose check failures.
func InstallationStatus(tool Tool, env *environment.Environment) (bool, error) {
	checker, ok := tool.(StatusChecker)
	if !ok {
		return tool.IsInstalled(env), nil
	}
	return checker.InstallationStatus(env)
}

// CommandInstallationStatus verifies a command using the same executable path
// as the tool's Verify method. A missing executable is a normal missing state;
// a command that starts but fails is an exceptional unknown state.
func CommandInstallationStatus(name string, args ...string) (bool, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err == nil {
		return true, nil
	}

	var commandErr *exec.Error
	if errors.As(err, &commandErr) && errors.Is(commandErr.Err, exec.ErrNotFound) {
		return false, nil
	}

	detail := strings.TrimSpace(string(output))
	if detail == "" {
		return false, err
	}
	return false, fmt.Errorf(messages.CommonCommandErrorDetail, err, detail)
}
