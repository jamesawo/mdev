package podman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	brew "github.com/jamesawo/mdev/internal/infrastructure/packagemanager"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type Podman struct{}

// Tool contract methods below describe Podman metadata, authoritative state,
// cancellable provisioning, managed storage, and uninstall behavior.
func (*Podman) Name() string                                   { return "podman" }
func (*Podman) Description() string                            { return messages.ToolsPodmanDescription }
func (*Podman) Dependencies() []string                         { return nil }
func (*Podman) SystemPrerequisites() []string                  { return []string{"brew"} }
func (*Podman) StorageDir(env *environment.Environment) string { return storage.ToolDir(env, "podman") }
func (p *Podman) IsInstalled(env *environment.Environment) bool {
	installed, _ := p.InstallationStatus(env)
	return installed
}
func (p *Podman) InstallationStatus(env *environment.Environment) (bool, error) {
	installed, err := tools.CommandInstallationStatus("podman", "--version")
	if err != nil || !installed {
		return installed, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	managed, err := tools.ManagedSymlinkStatus(podmanMachineDir(home), p.StorageDir(env))
	if err != nil || !managed {
		return managed, err
	}
	state, err := inspectMachineState(context.Background(), podmanMachineDir(home))
	return state.managed && state.complete, err
}
func (p *Podman) Install(env *environment.Environment) error {
	return p.InstallContext(context.Background(), env)
}
func (*Podman) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return brew.InstallContext(ctx, "podman")
}
func (p *Podman) Configure(env *environment.Environment) error {
	return p.ConfigureContext(context.Background(), env)
}
func (p *Podman) ConfigureContext(ctx context.Context, env *environment.Environment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if err := storage.Relocate(podmanMachineDir(home), p.StorageDir(env)); err != nil {
		return err
	}
	state, err := inspectMachineState(ctx, podmanMachineDir(home))
	if err != nil {
		return err
	}
	if !state.registered {
		return command.RunContext(ctx, "podman", "machine", "init")
	}
	if state.complete {
		return nil
	}
	if !state.managed {
		return fmt.Errorf(messages.ToolsPodmanUnmanagedMachine, state.name)
	}
	if err := command.RunContext(ctx, "podman", "machine", "rm", "--force", state.name); err != nil {
		return err
	}
	return command.RunContext(ctx, "podman", "machine", "init")
}

func podmanMachineDir(home string) string {
	return filepath.Join(home, ".local", "share", "containers", "podman", "machine")
}

type podmanMachineState struct {
	registered bool
	managed    bool
	complete   bool
	name       string
}

type podmanMachineInspection struct {
	Name      string `json:"Name"`
	ConfigDir struct {
		Path string `json:"Path"`
	} `json:"ConfigDir"`
}

type podmanMachineConfig struct {
	SSH struct {
		IdentityPath string `json:"IdentityPath"`
	} `json:"SSH"`
	ImagePath struct {
		Path string `json:"Path"`
	} `json:"ImagePath"`
}

// inspectMachineState verifies that Podman's registered machine owns real
// identity and disk artifacts beneath the conventional managed machine path.
func inspectMachineState(ctx context.Context, machineDir string) (podmanMachineState, error) {
	output, err := exec.CommandContext(ctx, "podman", "machine", "inspect").Output()
	if err != nil {
		return podmanMachineState{}, nil
	}
	var inspected []podmanMachineInspection
	if err := json.Unmarshal(output, &inspected); err != nil {
		return podmanMachineState{}, fmt.Errorf(messages.ToolsPodmanInspectOutput, err)
	}
	if len(inspected) == 0 {
		return podmanMachineState{}, nil
	}
	machine := inspected[0]
	state := podmanMachineState{registered: true, name: machine.Name}
	if machine.Name == "" || machine.ConfigDir.Path == "" {
		return state, nil
	}
	content, err := os.ReadFile(filepath.Join(machine.ConfigDir.Path, machine.Name+".json"))
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return state, fmt.Errorf(messages.ToolsPodmanReadConfig, err)
	}
	var config podmanMachineConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return state, fmt.Errorf(messages.ToolsPodmanParseConfig, err)
	}
	artifacts := []string{config.SSH.IdentityPath, config.ImagePath.Path}
	for _, artifact := range artifacts {
		if artifact == "" || !pathWithin(machineDir, artifact) {
			return state, nil
		}
	}
	state.managed = true
	for _, artifact := range artifacts {
		info, err := os.Stat(artifact)
		if errors.Is(err, os.ErrNotExist) {
			return state, nil
		}
		if err != nil {
			return state, err
		}
		if info.IsDir() {
			return state, nil
		}
	}
	state.complete = true
	return state, nil
}

func pathWithin(root, candidate string) bool {
	root, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	candidate, err = filepath.Abs(candidate)
	if err != nil {
		return false
	}
	relative, err := filepath.Rel(root, candidate)
	return err == nil && relative != ".." && !filepath.IsAbs(relative) && relative != "." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}
func (p *Podman) Verify(env *environment.Environment) error {
	return p.VerifyContext(context.Background(), env)
}
func (*Podman) VerifyContext(ctx context.Context, env *environment.Environment) error {
	if err := exec.CommandContext(ctx, "podman", "--version").Run(); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	machineDir := podmanMachineDir(home)
	managed, err := tools.ManagedSymlinkStatus(machineDir, (&Podman{}).StorageDir(env))
	if err != nil {
		return err
	}
	state, err := inspectMachineState(ctx, machineDir)
	if err != nil {
		return err
	}
	if !managed || !state.managed || !state.complete {
		return errors.New(messages.ToolsPodmanManagedStateIncomplete)
	}
	return nil
}

// Uninstall deregisters only a machine whose artifacts are verified beneath
// mdev-managed storage, then removes the Podman CLI without touching add-ons.
func (*Podman) Uninstall(_ *environment.Environment) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	state, err := inspectMachineState(context.Background(), podmanMachineDir(home))
	if err != nil {
		return err
	}
	if state.registered {
		if !state.managed {
			return fmt.Errorf(messages.ToolsPodmanUnmanagedMachine, state.name)
		}
		if err := command.Run("podman", "machine", "rm", "--force", state.name); err != nil {
			return err
		}
	}
	return brew.Uninstall("podman")
}

func init() { tools.Register(&Podman{}) }
