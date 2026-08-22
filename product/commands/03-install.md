# install

## purpose

`mdev install <tool>` makes a registered development tool available to the
user. A tool owns how it is installed, configured, verified, stored, and
uninstalled. The command owns only generic resolution and lifecycle
orchestration.

Install never contains package-manager, download, version-manager, path,
symlink, or other tool-specific strategy rules. It does not install or manage
services and does not start a tool's normal runtime lifecycle.

## usage

Supported invocation modes are:

```sh
mdev install <tool>
mdev install
mdev install --all
```

The no-argument form presents a deterministic interactive multi-select. `--all`
installs every registered tool. `--all` cannot be combined with a positional
tool. The global `--yes` flag accepts install confirmation without prompting.

Tool names use each registered Tool's canonical lowercase `Name()`. Resolution
is exact; install does not add fuzzy matching or aliases. Unknown tools fail
with an actionable recommendation to run `mdev list`.

Install has no JSON mode and does not expose update behavior. Re-running install
must never implicitly upgrade a validly installed tool.

## preflight and planning

Install requires readable mdev configuration and an existing, writable storage
directory. It does not create replacement configuration/storage or fall back to
another path when configured storage is unavailable.

The requested tools and their transitive `Dependencies()` form a deterministic,
deduplicated, dependency-first plan. A missing dependency or cycle fails before
mutation. Invalid unrelated registry entries must not prevent a single-tool
plan.

The plan is shown before mutation and requires confirmation unless `--yes` is
active. Declining confirmation is a successful no-op.

## lifecycle

For each planned tool, install uses the Tool abstraction and the error-aware
installation-status path:

1. determine authoritative completed-tool state;
2. if complete, skip all mutation;
3. call `Install`;
4. call `Configure`;
5. call `Verify`;
6. report success only after verification succeeds.

Installed state means the complete tool lifecycle is valid for mdev, not merely
that a command or partially created directory exists. Concrete tools own that
determination. Generic install must not call `StorageDir`, inspect tool paths, or
infer filesystem state.

Already-installed output remains version-free until reliable version metadata
exists:

```text
<tool> is already installed.
Uninstall: mdev uninstall <tool>
```

Install does not fabricate versions, perform updates, or start normal runtime
state. Podman provisioning may initialize required machine configuration, but
install does not start the machine.

## progress and output

Meaningful stable text is written progressively when real work begins or
becomes known. Output is not buffered until completion. Install, configure,
verify, skip, completion, cancellation, and failure boundaries are observable
through an injected reporter so presentation can evolve independently from
orchestration.

This iteration uses stable plain text in terminals and redirected/non-interactive
output. It does not reuse the terminal-only doctor spinner, add a spinner
dependency, or introduce a broad event/streaming framework.

All mdev-owned output uses Cobra's configured streams. Tool-native subprocess
output may continue to stream through the tool's infrastructure.

## cancellation and failure recovery

Install accepts context and supports Ctrl+C. Built-in subprocess-backed tools
use an optional context-aware lifecycle capability while the existing `Tool`
interface remains the compatibility contract.

Cancellation stops the active operation where supported and starts no later
stage or tool. Successfully completed dependencies and operations remain. There
is no broad rollback. Generic install cleans no guessed tool state; a concrete
tool may clean only temporary artifacts it clearly owns.

Retry favors reconciliation. Completed dependencies are skipped, while a
partially completed requested tool is not falsely reported installed. Ambiguous
filesystem state is left untouched with an actionable error.

Errors identify the tool and lifecycle phase and preserve the underlying cause.
Normal cancellation is concise (`Installation cancelled.`), not a stack trace.
Operational failures return through Cobra with non-zero exit status.

## storage and symlinks

Tool-managed state remains directly under `<storage_path>/<tool>` where
appropriate. Relocation and conventional-path symlinks remain tool-specific.

When relocating state:

- move the source when the destination does not exist;
- an empty destination created for the current operation may be removed only
  after validation before moving the source;
- fail without mutation when both source and destination contain data;
- never merge, overwrite, or delete existing user data;
- accept an existing symlink only when it resolves to the expected managed
  target.

## architecture

The dependency direction is:

```text
cmd/install.go
  -> internal/command/install
  -> tools.Tool and tools lifecycle/dependency abstractions
  -> concrete tool implementation
  -> infrastructure
```

`cmd/install.go` is thin Cobra wiring. `internal/command/install` owns preflight,
selection, planning, confirmation, orchestration, cancellation, progressive
reporting, and contextual errors. Concrete tools own all installation strategy
and tool-specific state.

Do not add PostgreSQL, `mdev update`, JSON output, service lifecycle behavior,
tool-name switches, a new spinner dependency, or speculative architectural
layers as part of this work.
