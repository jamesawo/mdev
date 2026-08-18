# setup

## purpose

`mdev setup` prepares mdev for first use by configuring where mdev stores its managed data.

setup should feel easy, calm, and deliberate. a new user should be able to complete it without understanding how mdev works internally.

## first run

run:

```bash
mdev setup
```

show a simple welcome and ask where mdev should store data.

```text
welcome to mdev.

where should mdev store data?

> ~/mdev (recommended)
  sandisk                  /Volumes/SanDisk
  choose another location
```

- `~/mdev` is the default.
- pressing enter accepts the default.
- detected external drives are offered as convenient choices.
- show the drive name with its path as subdued secondary text.
- show multiple external drives when available.
- read-only drives may be shown but cannot be selected.
- if no external drive is available, show only the default and custom-location options.

external storage is optional, not required.

## storage path

the persisted configuration is:

```yaml
storage_path: /absolute/path/to/mdev
```

`storage_path` always points to the final mdev-owned directory.

examples:

```text
~/mdev
/Volumes/SanDisk/mdev
/Users/james/Documents/mdev
```

managed tool data lives directly below it:

```text
<storage_path>/<tool>
```

there is no intermediate `data` directory.

setup creates only `storage_path`. tool-specific directories are created when the corresponding tools need them.

## choosing a location

when a parent location is selected, mdev creates an `mdev` directory inside it.

for example:

```text
/Volumes/SanDisk
```

becomes:

```text
/Volumes/SanDisk/mdev
```

if the supplied path already ends in `mdev`, do not append another `mdev`.

the user may choose any suitable writable filesystem location, including:

- the internal disk;
- an external drive;
- a directory under their home;
- a network-mounted filesystem;
- another writable macOS filesystem location.

mdev should not assume storage must live under `/Volumes`.

## custom locations

`choose another location` allows the user to enter a filesystem path.

support:

- `~`;
- relative paths;
- paths containing spaces.

normalize the final path before storing it:

- expand `~`;
- resolve relative paths;
- clean the path;
- resolve symlinks as described below;
- persist a canonical absolute path.

do not expand arbitrary environment variables such as `$HOME`.

if an entered location does not exist:

```text
that location doesn't exist.

<path>

> create it
  try again
  cancel
```

the preferred/default action comes first.

if the location is not writable, explain the problem and offer another location or cancellation.

if the location requires administrator privileges, explain that the user may rerun setup with `sudo`. mdev must never elevate itself automatically.

when running under `sudo`, the user's home and mdev configuration still belong to the actual invoking user, not root.

## existing mdev directory

if the selected location already contains an `mdev` directory, do not overwrite, clean, or delete it.

tell the user:

```text
an mdev folder already exists here.

> use existing folder
  choose another location
  cancel
```

using the existing folder is the default.

## symlinked locations

symlinked storage locations are supported.

before accepting one, show that it resolves somewhere else:

```text
this location is a symlink.

~/storage  → /Volumes/SanDisk/dev

mdev will use:  /Volumes/SanDisk/dev/mdev

> continue
  choose another location
  cancel
```

if accepted, persist the resolved physical path.

## already configured

if mdev is already configured, show that setup is complete and display the current storage path.

setup must not change storage, reset configuration, migrate data, or move or copy managed data.

## non-interactive setup

advanced users can configure setup directly:

```bash
mdev setup --storage-path ~/mdev
```

this performs the same validation and safety checks without interactive prompts when no decision is required.

if the supplied location does not exist, create it automatically.

if mdev is already configured, `--storage-path` must refuse to replace the existing configuration and show the current storage path.

`setup` does not support `--yes`.

## successful setup

only create `~/.mdev/` and persist configuration once setup can complete successfully.

a successful setup creates:

```text
~/.mdev/config.yaml
<storage_path>/
```

it does not pre-create tool directories, modify shell configuration, or install tools.

after completion:

```text
mdev is ready.

storage: ~/mdev

see what's available:

mdev list
```

## cancellation and interruption

users should be able to change their mind safely.

when `cancel` is offered, cancelling must not leave partial configuration, directories, or unnecessary data behind.

Ctrl+C before successful completion follows the same principle.

a canceled setup does not remember previous selections.

## unavailable storage

if configured storage later becomes unavailable, never silently fall back to another location or create another storage directory.

acknowledge the missing storage and offer a way forward:

```text
storage is not available.

expected: /Volumes/SanDisk/mdev

> connect drive and try again
  cancel
```

do not fall back to another location, create new storage, or change configuration.

## malformed or unreadable configuration

if `~/.mdev/config.yaml` is malformed or cannot be read, leave it untouched.

calmly tell the user to fix or remove the configuration manually and try again. exit as a failure.

## filesystem behavior

setup cares about the filesystem capabilities mdev requires rather than the filesystem's brand or format.

a location must be suitable for mdev required filesystem behavior.

if the selected path exists as a file rather than a directory, explain the conflict and allow the user to choose another location.

setup does not perform disk-capacity checks. commands such as `install` should evaluate capacity when they know what storage an operation requires.

setup creates no marker or metadata files inside `storage_path`.

## relationship with other commands

setup is the first step before commands that depend on mdev storage.

commands requiring configuration should direct an unconfigured user to:

```bash
mdev setup
```

`doctor` should also detect when setup has not been completed and recommend it.

setup does not install development tools.

setup does not own reset/reinstallation behavior.

## cli experience

all mdev user-facing CLI copy should use simple, calm language and lowercase styling.

output should feel deliberate rather than dumping many messages onto the terminal at once. preserve and improve the existing progressive/streamed presentation where appropriate.

use spinners for operations that genuinely take noticeable time. normal setup operations should be fast and do not need artificial progress indicators.

do not clear or redraw the terminal between setup steps. preserve normal terminal history.

Cobra's normal help system should provide setup help with concise, useful command and flag descriptions.

errors should be graceful and help the user understand what they can do next.

## testing

### unit

cover smaller behavior such as:

- path normalization;
- path validation;
- `mdev` directory resolution;
- configuration persistence;
- external-drive discovery;
- symlink handling;
- writable/unwritable locations;
- cancellation behavior;
- relevant failure and edge cases.

### e2e

the setup E2E test exercises the real happy-path user journey using the compiled `mdev` binary inside the macOS VM.

E2E testing should validate the experience a real user goes through rather than duplicating every small edge case covered by unit tests.

## done when

setup is done when:

- a new user can complete the happy path easily;
- internal and external storage work;
- custom paths work;
- `--storage-path` works;
- configuration contains the canonical `storage_path`;
- cancellation leaves no partial setup behind;
- setup never unexpectedly overwrites user data;
- user-facing copy has been reviewed;
- relevant unit/component tests pass;
- the happy-path E2E journey passes in the macOS VM.
