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
elevates itself automatically.

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
        │ VM shared directory
        ▼
macOS VM: macOs-mdev
        │
        │ /usr/local/bin/mdev symlink
        ▼
  interactive SSH session
```

### One-time setup

1. Install [UTM](https://mac.getutm.app/) and create an Apple Silicon macOS VM named `macOs-mdev` with a local user
   named `mdev`. The VM does not need Go, GoLand, or a repository clone.

2. Share the host repository's build directory with the VM. The shared binary must be visible inside the VM; for the
   standard UTM share:

   ```text
   /Volumes/My Shared Files/dist/mdev
   ```

3. Create a stable command inside the VM:

   ```sh
   sudo mkdir -p /usr/local/bin
   sudo ln -sf "/Volumes/My Shared Files/dist/mdev" /usr/local/bin/mdev
   ```

   Rebuilding `dist/mdev` on the host immediately updates the command used by the VM.

4. In the VM, enable **System Settings → General → Sharing → Remote Login** and grant access to the `mdev` user.

5. Configure public-key authentication from the host. Generate or select a local SSH key, append its public key to
   `~/.ssh/authorized_keys` in the VM, and verify that login no longer requires the VM password. Private keys and
   passwords must remain outside this repository.

6. Add the expected alias to the host's `~/.ssh/config`:

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
./scripts/vm/shell.sh                 # build separately, then open a VM shell
./scripts/vm/run.sh mdev --version    # run any command in the VM
./scripts/vm/run.sh mdev doctor
./scripts/vm/test-setup.sh            # run isolated setup journeys in the VM
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
- `run mdev` runs `build mdev`, ensures the VM and SSH are ready, and opens an interactive SSH shell in GoLand's
  Terminal tool window.

The normal GoLand loop is therefore:

1. Edit the project on the host.
2. Run `run mdev`.
3. At the VM prompt, run `mdev --version`, `mdev doctor`, or another command.

### Troubleshooting

- **`utmctl` is unavailable:** install UTM in `/Applications`, or make `utmctl` available on `PATH`.
- **The VM cannot be found:** ensure its UTM name is `macOs-mdev`, or set `MDEV_VM_NAME` locally.
- **SSH times out:** start the VM in UTM, confirm Remote Login is enabled, and verify `ssh mdev-vm` independently.
- **`mdev` is missing in the VM:** confirm `dist/mdev` is visible through the UTM share and recreate the
  `/usr/local/bin/mdev` symlink.
