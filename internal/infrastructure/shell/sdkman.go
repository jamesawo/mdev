package shell

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

var sdkmanIdentifier = regexp.MustCompile(`^[a-z0-9-]+$`)

// RunWithSDKMAN executes a command after loading SDKMAN into the shell.
func RunWithSDKMAN(cmd string) error {
	return RunWithSDKMANContext(context.Background(), cmd)
}

// RunWithSDKMANContext executes an SDKMAN-backed command with cancellation.
func RunWithSDKMANContext(ctx context.Context, cmd string) error {
	bash, err := prerequisites.ModernBashPath(ctx)
	if err != nil {
		return err
	}
	return command.RunContext(
		ctx,
		bash,
		"-c",
		"source $HOME/.sdkman/bin/sdkman-init.sh && "+cmd,
	)
}

// InstallSDKMANCandidateContext installs and selects a candidate through the
// mdev-managed SDKMAN installation.
func InstallSDKMANCandidateContext(ctx context.Context, candidate string) error {
	if err := validateSDKMANIdentifier(candidate); err != nil {
		return err
	}
	return RunWithSDKMANContext(ctx, "sdk install "+candidate)
}

// SDKMANCandidateInstallationStatus verifies the executable owned by a
// candidate's current SDKMAN installation, ignoring unrelated PATH commands.
func SDKMANCandidateInstallationStatus(ctx context.Context, candidate, executable string, args ...string) (bool, error) {
	path, err := sdkmanCandidateExecutable(candidate, executable)
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	cmd, err := sdkmanCandidateCommand(ctx, path, args...)
	if err != nil {
		return false, err
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf(messages.CommonCommandErrorDetail, err, strings.TrimSpace(string(output)))
	}
	return true, nil
}

// RunSDKMANCandidateContext executes the current executable owned by an SDKMAN
// candidate after loading the mdev-managed SDKMAN environment.
func RunSDKMANCandidateContext(ctx context.Context, candidate, executable string, args ...string) error {
	path, err := sdkmanCandidateExecutable(candidate, executable)
	if err != nil {
		return err
	}
	cmdArgs := append([]string{"-c", `source "$HOME/.sdkman/bin/sdkman-init.sh" && exec "$@"`, "mdev-sdkman", path}, args...)
	bash, err := prerequisites.ModernBashPath(ctx)
	if err != nil {
		return err
	}
	return command.RunContext(ctx, bash, cmdArgs...)
}

// UninstallSDKMANCandidate removes only the current candidate version owned by
// the mdev-managed SDKMAN installation.
func UninstallSDKMANCandidate(candidate string) error {
	return UninstallSDKMANCandidateContext(context.Background(), candidate)
}

// UninstallSDKMANCandidateContext removes only the current managed candidate
// version with caller cancellation.
func UninstallSDKMANCandidateContext(ctx context.Context, candidate string) error {
	path, err := sdkmanCandidateExecutable(candidate, candidate)
	if err != nil {
		return err
	}
	currentLink := filepath.Dir(filepath.Dir(path))
	info, err := os.Lstat(currentLink)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf(messages.ToolsSDKMANCurrentNotSymlink, candidate)
	}
	originalTarget, err := os.Readlink(currentLink)
	if err != nil {
		return err
	}
	current, err := filepath.EvalSymlinks(currentLink)
	if err != nil {
		return err
	}
	candidateRoot := filepath.Dir(currentLink)
	resolvedRoot, err := filepath.EvalSymlinks(candidateRoot)
	if err != nil {
		return err
	}
	if current == resolvedRoot || !strings.HasPrefix(current, resolvedRoot+string(os.PathSeparator)) {
		return fmt.Errorf(messages.ToolsSDKMANCurrentOutsideCandidate, candidate, current)
	}
	version := filepath.Base(current)
	if err := os.Remove(currentLink); err != nil {
		return err
	}
	if err := RunWithSDKMANContext(ctx, "sdk uninstall "+candidate+" "+version); err != nil {
		restoreErr := os.Symlink(originalTarget, currentLink)
		if restoreErr != nil {
			return errors.Join(err, fmt.Errorf(messages.ToolsSDKMANRestoreCurrent, candidate, restoreErr))
		}
		return err
	}
	return nil
}

func sdkmanCandidateCommand(ctx context.Context, path string, args ...string) (*exec.Cmd, error) {
	bash, err := prerequisites.ModernBashPath(ctx)
	if err != nil {
		return nil, err
	}
	cmdArgs := append([]string{"-c", `source "$HOME/.sdkman/bin/sdkman-init.sh" && exec "$@"`, "mdev-sdkman", path}, args...)
	return exec.CommandContext(ctx, bash, cmdArgs...), nil
}

func sdkmanCandidateExecutable(candidate, executable string) (string, error) {
	if err := validateSDKMANIdentifier(candidate); err != nil {
		return "", err
	}
	if err := validateSDKMANIdentifier(executable); err != nil {
		return "", err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sdkman", "candidates", candidate, "current", "bin", executable), nil
}

func validateSDKMANIdentifier(value string) error {
	if !sdkmanIdentifier.MatchString(value) {
		return fmt.Errorf(messages.ToolsInvalidSDKMANIdentifier, value)
	}
	return nil
}
