package prerequisites

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Bash represents the modern Bash capability required by SDKMAN.
type Bash struct{}

var errBashOutdated = errors.New("installed bash is older than version 4")

func (Bash) Name() string { return "bash" }
func (Bash) Check() bool {
	state, _, _ := (Bash{}).Readiness(context.Background())
	return state == StateReady
}
func (Bash) Install() error                     { return (Bash{}).RemediateContext(context.Background()) }
func (Bash) PrerequisiteDependencies() []string { return []string{"brew"} }
func (Bash) RemediationDescription() string     { return "install Bash 4 or newer with Homebrew" }
func (Bash) RemediateContext(ctx context.Context) error {
	return runCommand(ctx, "brew", "install", "bash")
}
func (Bash) VerifyContext(ctx context.Context) error {
	_, err := ModernBashPath(ctx)
	return err
}
func (Bash) Readiness(ctx context.Context) (State, string, error) {
	path, err := ModernBashPath(ctx)
	if err == nil {
		return StateReady, path, nil
	}
	if errors.Is(err, exec.ErrNotFound) {
		return StateMissing, "Bash 4 or newer is required", nil
	}
	if errors.Is(err, errBashOutdated) {
		return StateOutdated, "Bash 4 or newer is required", nil
	}
	return StateBroken, "", err
}

// ModernBashPath returns a verified Bash 4+ executable independent of PATH order.
func ModernBashPath(ctx context.Context) (string, error) {
	var candidates []string
	if path, err := exec.LookPath("bash"); err == nil {
		candidates = append(candidates, path)
	}
	candidates = append(candidates, "/opt/homebrew/bin/bash", "/usr/local/bin/bash")
	seen := map[string]bool{}
	found := false
	for _, path := range candidates {
		if seen[path] {
			continue
		}
		seen[path] = true
		if _, err := os.Stat(path); err != nil {
			continue
		}
		found = true
		output, err := exec.CommandContext(ctx, path, "-c", "printf '%s' \"$BASH_VERSINFO\"").Output()
		if err != nil {
			continue
		}
		major, err := strconv.Atoi(strings.TrimSpace(string(output)))
		if err == nil && major >= 4 {
			return path, nil
		}
	}
	if found {
		return "", errBashOutdated
	}
	return "", fmt.Errorf("modern bash: %w", exec.ErrNotFound)
}

func init() { Register(Bash{}) }
