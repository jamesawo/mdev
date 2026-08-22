# version

## purpose

`mdev version` reports the identity of the mdev binary that is currently running.

The command provides a small, stable set of metadata for people diagnosing a local installation and for scripts that
need the same information in machine-readable form. It is read-only and does not require mdev setup, configured storage,
network access, Git, or any installed development tool.

## basic usage

Run:

```bash
mdev version
```

The normal output contains the product name and version, source commit identifier, and build value:

```text
mdev 0.2.0
commit: a81f37c
built: 2026-08-22
```

Keep this order and these labels stable:

1. `mdev <version>`
2. `commit: <commit>`
3. `built: <built>`

Do not include the project author, decorative headings, status symbols, explanatory paragraphs, colors, progress
indicators, or blank lines in the normal output.

## metadata fields

Version metadata contains exactly three public values:

- `version`: the mdev product version, without a leading `v`;
- `commit`: the source commit identifier when supplied, preferably the short form used for human display;
- `built`: the build date or other stable build value when supplied.

The command reports supplied values as-is after applying only the local fallback rules below. It does not query Git,
derive values from repository state, inspect tags, contact a service, or generate metadata at runtime.

## local development defaults

This project is currently local-development only. When real build metadata has not been supplied, use deterministic
defaults:

```text
version: 0.2.0
commit: unknown
built: unknown
```

The corresponding default human output is:

```text
mdev 0.2.0
commit: unknown
built: unknown
```

Do not substitute the current time, current date, working-tree state, Git branch, Git commit, hostname, username, or
filesystem path. Repeated executions of the same binary must report the same metadata.

The implementation should accept metadata as an explicit input at the command boundary so a future build can supply
real values without changing rendering or command behavior. This work does not add linker flags, Git-derived metadata,
tag parsing, release automation, deployment infrastructure, or CI/CD configuration.

## json output

Run:

```bash
mdev version --json
```

JSON output is one complete object containing the same underlying values as normal output:

```json
{
  "version": "0.2.0",
  "commit": "a81f37c",
  "built": "2026-08-22"
}
```

The object contains exactly these required string fields:

- `version`;
- `commit`;
- `built`.

JSON standard output must contain valid JSON only. Do not include human labels, the `mdev` prefix, headings, colors,
progress messages, or explanatory text. Emit the JSON as one complete document followed by a newline.

Human and JSON output must use one metadata source and one command workflow. Do not create a separate metadata-detection
path for JSON.

## root version flag

Preserve support for:

```bash
mdev --version
```

The root flag is the concise compatibility form and prints:

```text
mdev 0.2.0
```

It uses the same `version` value as `mdev version`. Commit and build values remain available through the detailed
subcommand. The root flag does not accept `--json`; machine-readable output belongs to `mdev version --json`.

## cli behavior

`mdev version`:

- accepts no positional arguments;
- supports only its `--json` output flag in addition to Cobra help;
- does not require or use `--yes`;
- does not prompt for input;
- writes through Cobra's configured output writer;
- returns success after producing valid output;
- returns a failure if output cannot be written.

All user-facing copy and formatting templates belong in `internal/ui/messages`.

## read-only behavior

Running either version form must not:

- read or write `~/.mdev/config.yaml`;
- require configured or available mdev storage;
- inspect the Git repository or working tree;
- execute subprocesses;
- use the network;
- create files or directories;
- modify tools, services, shell state, or environment configuration;
- prompt or ask for confirmation.

## ownership and package boundaries

`cmd/version.go` owns only Cobra wiring:

- command name, arguments, help, and `--json` flag registration;
- selection of Cobra's output writer;
- construction of the metadata and options passed to the workflow;
- delegation to `internal/command/version`.

`internal/command/version` owns:

- the metadata value type;
- local-development fallback values;
- human and JSON rendering;
- validation or normalization required by the stated fallback contract;
- propagation of writer and encoding errors.

`internal/ui/messages/version.go` owns human-facing labels, help copy, command/flag names, and formatting templates.

`cmd/root.go` continues to own root Cobra wiring for `mdev --version`, but obtains its version value from the same
metadata source used by the version subcommand.

Do not add a release package, build service, metadata provider interface, repository abstraction, Git adapter, or new
dependency for this feature.

## migration from current behavior

The current implementation must be changed as follows:

- move formatting and execution out of `cmd/version.go` into `internal/command/version`;
- replace direct `fmt.Println` calls with the Cobra-provided output writer;
- remove the hard-coded author field and `Created by ...` output;
- replace the current `0.1.0` value with the approved local default `0.2.0`;
- add deterministic `commit` and `built` fallback values;
- add `--json` without changing the normal three-line output contract;
- preserve both `mdev version` and `mdev --version`;
- keep all user-facing version copy in `internal/ui/messages/version.go`.

## testing

### unit

Tests beside `internal/command/version` cover:

- exact three-line human output and field ordering;
- exact development fallback output;
- supplied version, commit, and built values;
- JSON validity, required fields, values, and absence of human-oriented text;
- human and JSON modes using identical metadata input;
- deterministic repeated output;
- writer failure propagation;
- absence of configuration, filesystem, Git, subprocess, and network dependencies.

### command

Tests beside `cmd/version.go` cover:

- thin delegation with the expected metadata, writer, and JSON option;
- `--json` registration and help copy;
- rejection of positional arguments;
- normal error propagation;
- `mdev --version` remaining available and using the same version value;
- no author output.

### e2e

The version E2E journey builds the macOS ARM64 binary on the host and executes it only inside the macOS VM. It validates:

- exact `mdev version` output with development defaults;
- valid JSON-only output from `mdev version --json` with matching values;
- concise `mdev --version` output;
- successful operation without mdev configuration or storage.

E2E coverage should use an isolated VM location and must not execute mdev on the host.

## documentation

Document the human command, JSON mode, root compatibility flag, metadata meanings, and local-development fallback values
in `README.md`. Do not document future release or build automation as existing functionality.

## done when

Version is done when:

- `mdev version` prints exactly the version, commit, and built lines in the specified order;
- local builds deterministically report `0.2.0`, `unknown`, and `unknown` unless metadata is supplied;
- `mdev version --json` emits one valid JSON object with matching `version`, `commit`, and `built` strings;
- `mdev --version` remains supported and uses the same version value;
- `cmd/version.go` is a thin Cobra layer delegating to `internal/command/version`;
- no author line or direct terminal printing remains in the version command;
- all user-visible version copy is centralized under `internal/ui/messages`;
- the command remains completely read-only and independent of configuration, storage, Git, subprocesses, and network;
- focused unit and command tests pass;
- full host-safe tests, vetting, formatting, compilation, and diff checks pass;
- the human, JSON, and root-flag journeys pass in the macOS VM;
- documentation describes only implemented local-development behavior.
