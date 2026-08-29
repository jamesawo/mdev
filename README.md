# mdev

`mdev` is a Go CLI for setting up and managing a macOS development environment. It installs common development tools and
keeps their managed data under a configurable filesystem storage path.

## Managed storage

Run `mdev setup` to choose where mdev stores managed tool data. The storage path
can be a directory on the internal disk, under your home directory, on an
external volume, or at another writable filesystem location.

The interactive command recommends `~/mdev`, lists mounted external volumes,
and supports custom paths. A selected parent directory gets an `mdev`
subdirectory; a path already ending in `mdev` is used as-is. For unattended
setup, pass the location directly:

```sh
mdev setup --storage-path ~/mdev
```

Setup never replaces an existing configuration. Paths are stored as canonical
absolute paths, including resolved symlink destinations. Running through
`sudo` still writes configuration for the invoking user, and mdev never
elevates itself automatically. Running `mdev setup` again reports the existing
storage and ends by recommending `mdev list` as the next command.

Before reporting readiness, interactive setup checks the system prerequisites
needed by mdev and its supported tools. If preparation is required, setup shows
the concrete changes and asks for explicit consent with a default of no. Each
approved change is verified before configuration is committed. Declining,
failure, or cancellation leaves setup incomplete and safe to retry.

Some macOS prerequisites complete through a system installer. In that case,
setup starts the supported installer, explains what remains, and stops before
dependent changes. Complete the macOS installation and rerun `mdev setup`;
mdev rechecks current state and continues the remaining work.

`setup --storage-path` never treats the supplied path as permission to change
system tooling. If prerequisite remediation is required, it exits actionably
and asks the user to rerun interactive setup.

The selection is saved in `~/.mdev/config.yaml`:

```yaml
storage_path: /path/to/mdev
```

Each tool owns a directory directly below that path. There is no intermediate
`data` directory:

```text
/path/to/mdev/
├── gradle/
├── maven/
├── nvm/
└── ...
```

For tools whose normal data lives under the user's home directory, mdev may
relocate that data into the configured storage path and replace the original
location with a symlink.

## Install tools

Install one registered tool, choose tools interactively, or install all
registered tools:

```sh
mdev install gradle
mdev install
mdev install --all
```

mdev resolves dependencies deterministically, then displays only pending
lifecycle work and asks for confirmation. Already-complete dependencies remain
part of validation but are not presented as tools mdev will install. Use the
global `--yes` flag for unattended installation. Progress is written as plain
text while each tool is installed, configured, and verified:

```text
installing podman-compose... ✓
configuring podman-compose... ✓
verifying podman-compose... ✓
podman-compose installed.
```

mdev captures successful output from underlying installers, so Homebrew,
SDKMAN, NVM, and other provider details do not overwhelm normal command output.
If a provider fails, mdev identifies the affected tool and lifecycle phase and
includes a concise useful diagnostic from the captured output.

Uninstall follows the same managed presentation contract. Provider output is
captured while mdev reports truthful tool-removal and managed-storage cleanup
progress, for example:

```text
uninstalling podman... ✓
cleaning podman storage... ✓
podman removed.
```

System requirements such as Bash, Homebrew, Curl, Git, and Xcode Command Line
Tools are protected by their shared prerequisite classification. `mdev
uninstall <system-requirement>` explains that mdev cannot remove it and makes
no changes.

Re-running install is retry-safe: valid installations and completed
dependencies are skipped, while partial work resumes through the owning tool's
lifecycle. Install does not update valid tools or start their normal runtime
state. If configured storage is missing or unavailable, install fails instead
of creating or selecting a replacement location.

Java, Gradle, and Maven are installed through mdev's managed SDKMAN
installation. Existing Homebrew or other unmanaged JVM-tool installations are
left untouched and do not count as mdev-managed completion.

Podman support is split into three explicit tools:

```text
podman
├── podman-desktop
└── podman-compose
```

`mdev install podman` installs the CLI and initializes its macOS machine without
starting it. The machine's large data files remain below the configured mdev
storage root at `<storage_path>/podman`, reached through Podman's conventional
home-directory path. Install `podman-desktop` for the optional GUI or
`podman-compose` for the optional Compose provider; each depends on `podman`.
The add-ons do not own or duplicate the machine data. `mdev install --all`
includes all three tools. mdev does not add Docker aliases or modify shell
configuration.

Install defensively checks only the system prerequisites required by the
resolved tool plan. If the machine has drifted since setup, it fails before tool
mutation and recommends setup or doctor for recovery.

## List supported tools

After setup, run `mdev list` for a read-only overview of every registered tool:

```text
system tools
  ○ brew       missing
  ✓ curl       installed
  ✓ git        installed

tools
  ○ gradle     missing
  ✓ java       installed
  ○ maven      missing
```

System prerequisites and other tools are shown separately and sorted
alphabetically. `✓`, `○`, and `?` mean installed, missing, and unknown. An
unknown status indicates that verification failed unexpectedly; mdev continues
checking the remaining tools, reports the failure details, and exits with a
failure status. Normal terminal output is written progressively in alphabetical
display order, so results remain visible while later tools are still being
checked.

For automation, use `mdev list --json`. It reports the same checks as one
complete, deterministically ordered JSON document:

```json
{
  "system_tools": [
    {"name": "brew", "status": "missing"},
    {"name": "curl", "status": "installed"}
  ],
  "tools": [
    {"name": "gradle", "status": "missing"},
    {"name": "java", "status": "unknown", "error": "verification failed"}
  ]
}
```

JSON standard output contains JSON only and is written after all checks finish.
Unknown statuses are included in the complete document before the command exits
with a failure status.

List requires valid mdev configuration and available configured storage. It
does not create or repair configuration or storage, install or configure tools,
show dependency relationships, or diagnose installation health. Use
`mdev graph` for dependencies and `mdev doctor` for health diagnosis.

## Version information

Run `mdev version` to identify the binary currently running:

```text
mdev 0.2.0
commit: unknown
built: unknown
```

The version is the product version without a leading `v`. Commit identifies the
source revision when supplied, and built is the build date or another stable
build value when supplied. Local development builds use deterministic fallbacks:
version `0.2.0`, commit `unknown`, and built `unknown`.

For scripts, `mdev version --json` emits the same three values as one JSON-only
document:

```json
{"version":"0.2.0","commit":"unknown","built":"unknown"}
```

The compatibility form `mdev --version` prints only `mdev 0.2.0`. Both forms
are read-only and work without mdev configuration or storage.

## Local development

Run the CLI directly:

```sh
go run . doctor
```

Build a local executable:

```sh
make build
./mdev doctor
```

## macOS VM development

The project uses a disposable macOS VM to test machine-level setup without changing the host Mac. Source code and the Go
toolchain remain on the host; only the compiled binary runs in the VM.

```text
GoLand or host terminal
        │
        │ builds
        ▼
    dist/mdev
        │
        │ SSH deployment
        ▼
macOS VM: macOs-mdev
        │
        │ ~/.local/share/mdev/bin/mdev
        ▼
  interactive SSH session
```

### One-time setup

1. Install [UTM](https://mac.getutm.app/) and create an Apple Silicon macOS VM named `macOs-mdev` with a local user
   named `mdev`. The VM does not need Go, GoLand, or a repository clone.

2. In the VM, create the deployment directory and a stable command symlink:

   ```sh
   mkdir -p "$HOME/.local/share/mdev/bin"
   sudo mkdir -p /usr/local/bin
   sudo ln -sf "$HOME/.local/share/mdev/bin/mdev" /usr/local/bin/mdev
   ```

   The target binary is created automatically the first time `run mdev` deploys
   successfully.

3. In the VM, enable **System Settings → General → Sharing → Remote Login** and grant access to the `mdev` user.

4. Configure public-key authentication from the host. Generate or select a local SSH key, append its public key to
   `~/.ssh/authorized_keys` in the VM, and verify that login no longer requires the VM password. Private keys and
   passwords must remain outside this repository.

5. Add the expected alias to the host's `~/.ssh/config`:

   ```sshconfig
   Host mdev-vm
       HostName <VM_IP>
       User mdev
       IdentityFile ~/.ssh/<VM_KEY>
   ```

   Verify the alias before using the repository automation:

   ```sh
   ssh mdev-vm
   ```

VM addresses can change. If SSH stops connecting, determine the current address inside the VM and update the local
alias.

### Daily workflow

The repository provides these scripts:

```sh
./scripts/vm/start.sh                 # ensure the VM is running
./scripts/vm/wait.sh                  # wait for passwordless SSH
./scripts/vm/deploy-binary.sh         # deploy the active checkout's existing build
./scripts/vm/shell.sh                 # deploy the existing build, then open a VM shell
./scripts/vm/run.sh mdev --version    # run any command in the VM
./scripts/vm/run.sh mdev doctor
./scripts/vm/test-setup.sh            # run isolated setup journeys in the VM
./scripts/vm/test-list.sh             # run the configured list journey in the VM
./scripts/vm/test-install.sh          # run isolated install lifecycle journeys in the VM
./scripts/vm/test-readiness.sh        # run the controlled setup/readiness/install journey
./scripts/vm/test-version.sh          # run isolated version journeys in the VM
```

The scripts use the `mdev-vm` SSH alias and expect a VM named `macOs-mdev`. A developer with a differently named local
VM can override only the VM name:

```sh
MDEV_VM_NAME=my-local-vm ./scripts/vm/shell.sh
```

The default SSH wait timeout is 120 seconds. It can be adjusted for a slow boot:

```sh
MDEV_VM_WAIT_TIMEOUT=240 ./scripts/vm/shell.sh
```

### GoLand configurations

Two shared run configurations are included:

- `build mdev` builds the repository's main package for macOS ARM64 at `dist/mdev`.
- `run mdev` runs `build mdev`, ensures the VM and SSH are ready, deploys the
  active checkout's binary directly over SSH, and opens an interactive SSH
  shell in GoLand's Terminal tool window. This works the same way from the
  primary checkout and linked Codex worktrees.

The normal GoLand loop is therefore:

1. Edit the project on the host.
2. Run `run mdev`.
3. At the VM prompt, run `mdev --version`, `mdev doctor`, or another command.

### Troubleshooting

- **`utmctl` is unavailable:** install UTM in `/Applications`, or make `utmctl` available on `PATH`.
- **The VM cannot be found:** ensure its UTM name is `macOs-mdev`, or set `MDEV_VM_NAME` locally.
- **SSH times out:** start the VM in UTM, confirm Remote Login is enabled, and verify `ssh mdev-vm` independently.
- **`mdev` is missing in the VM:** confirm `~/.local/share/mdev/bin/mdev` exists
  after running `run mdev`, then recreate the `/usr/local/bin/mdev` symlink from
  the one-time setup.
- **The VM reports an older build:** confirm GoLand opened the intended branch
  or worktree, then run `run mdev` again. The deployment message confirms that
  the active checkout's freshly built binary was uploaded.
