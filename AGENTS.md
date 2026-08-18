# AGENTS.md

## Purpose and agent role

`mdev` is a Go CLI for configuring a macOS development environment. It installs
and removes development tools, manages selected tool data below a configured
storage root, and starts local development services.

An agent's job is to help make `mdev` production-grade through investigation,
implementation, scoped refactoring, testing, VM validation, debugging,
documentation, and quality and safety improvements.

The governing principle is: broad implementation autonomy inside a defined
task; strict boundaries around product decisions, architecture, dependencies,
compatibility, destructive behavior, Git publication, and host execution.

Agents do not own product direction, breaking-change decisions, architecture
redesign, dependency adoption, data migration policy, destructive scope, or
release decisions. Those require human authority as specified below.

## Scope control

Implement only the requested task and the smallest supporting changes required
to complete it correctly.

- Read the relevant code, tests, documentation, and current `git status` before
  editing.
- Preserve all pre-existing staged, unstaged, and untracked work.
- Do not add unrelated features, broad cleanup, speculative abstractions,
  package reorganization, or architectural redesign.
- Do not implement TODO comments unless the task requires them.
- Prefer the smallest clean change and existing repository conventions.
- Report useful out-of-scope findings as follow-ups instead of implementing
  them.

## Existing architecture

Protect the pragmatic architecture that exists today:

- `main.go` is the executable entry point.
- `cmd/` defines Cobra commands, flags, arguments, command registration, and
  process-level behavior.
- `internal/command/` coordinates workflows such as doctor, install,
  uninstall, list, graph, and up.
- `internal/tools/` defines the `Tool` lifecycle, registry, and dependency
  resolution. Concrete tools live in subpackages and register through `init`.
- `internal/services/` defines runtime services and their manager.
- `internal/infrastructure/` contains concrete configuration, environment,
  filesystem, storage, package-manager, shell, prerequisite, and subprocess
  interactions.
- `internal/ui/` contains confirmation, interactive selection, user-facing
  messages, and terminal rendering.
- `test/` contains black-box-style tests that import internal packages.

These are practical boundaries. They are not a mandate to invent domain,
application, repository, adapter, or dependency-injection layers.

### Boundary rules

- Keep Cobra-specific parsing and registration in `cmd/`.
- Put reusable workflow and planning logic under the relevant
  `internal/command/<name>` package.
- Keep tool lifecycle behavior behind the existing `tools.Tool` contract.
- Keep service lifecycle behavior behind the existing `services.Service`
  contract.
- Keep OS process, filesystem, package-manager, and persistent-state access in
  the existing infrastructure or concrete tool/service boundaries.
- Use the established `internal/ui` packages for CLI interaction and output;
  keep user-facing text centralized in `internal/ui/messages` where practical.
- Keep interactions with external commands and macOS facilities isolated and
  testable. Do not spread new raw subprocess execution across unrelated
  packages.
- Do not move packages, replace registries, or change core interfaces merely
  to make the architecture look more layered.
- When existing code crosses a boundary, improve it only as far as the current
  task requires. Do not turn a local task into a repository-wide refactor.

## Product specifications

- `product/commands/` contains the committed product source of truth for each
  command. Read the relevant specification before working on a command.
- `product/work/` contains temporary agent progress and checklist state.
- `product/research/` contains unresolved, private product thinking.
- Product specifications define intended behavior and user experience;
  `AGENTS.md` defines engineering and agent rules.
- Do not change product specifications, promote research into them, or change
  product intent without explicit human approval.

### Command work files

`product/work/<command>.md` tracks temporary implementation progress toward
the corresponding `product/commands/<command>.md`. Use this format:

```markdown
# <command> work
## status
in progress
## human decisions / blockers
- [x] <resolved decision> — <decision made>
- [ ] <decision or blocker requiring human input>
## todo
- [ ] <work item>
- [ ] <work item>
```

- Status is `in progress` or `done`.
- Human decisions / blockers records anything requiring human authority or
  preventing progress. Keep resolved decisions and mark them `[x]`; do not
  delete them.
- Codex may add TODOs when implementation reveals additional work required to
  satisfy the product specification.
- Do not delete or rewrite existing TODO lines. Mark them `[x]` when completed.
- Do not add unrelated discoveries or expand the command's scope.
- Keep todo as the final section so required work can be appended naturally.
- Set status to `done` only when the command satisfies its product specification
  and the repository definition of done.
- Work files are temporary and gitignored; command specifications remain the
  durable source of truth.

## Persistent state and backward compatibility

Existing installations depend on:

- configuration created at `~/.mdev/config.yaml` when setup occurs;
- the YAML key `storage_path`;
- managed tool data at `<storage_path>/<tool-specific-path>`, with no mandatory
  intermediate `data` directory;
- per-tool storage paths directly below the configured storage path;
- symlinks from selected home-directory locations into managed storage;
- registered tool names and dependencies.

Treat the following as public or compatibility-sensitive behavior:

- command and flag names, positional arguments, and exit status;
- meaningful user-facing output and errors;
- confirmation behavior, including `--yes`;
- tool and service identifiers;
- install, configure, verify, and uninstall ordering;
- dependency and reverse-dependency behavior;
- configuration location and schema;
- managed directory and symlink layout;
- defaults such as versions and storage locations;
- both `mdev version` and `mdev --version`.

Preserve backward compatibility unless the task explicitly authorizes a
breaking change or the human approves one. Do not silently migrate, delete,
reinterpret, or overwrite existing configuration or managed data.

Any approved change to public CLI behavior, persisted state, storage layout,
or defaults must include corresponding tests and documentation.

## Go conventions

- Use the Go version declared in `go.mod` and run `gofmt` on changed Go files.
- Prefer clear, focused functions and explicit names over cleverness.
- Follow existing package patterns unless the task requires a local
  improvement.
- Return errors rather than panic for expected operational failures.
- Add useful error context while preserving causes where relevant; do not
  silently swallow errors.
- Prefer deterministic ordering for user-visible output and execution plans.
- Pass `context.Context` through workflows that already use it.
- Do not add interfaces solely for hypothetical future mocking. Add or extend
  one only for a concrete implementation or test seam required by the task.
- Do not introduce generic `utils`, `helpers`, or similar dumping-ground
  packages.
- Avoid premature optimization and broad rewrites of working code.

## Dependencies

The current direct dependencies are Cobra, Survey v2, and YAML v3.

Prefer the standard library and existing dependencies. Agents may identify and
recommend dependencies, but adding, replacing, or materially upgrading a
direct dependency requires prior human approval.

Before requesting approval, explain why the standard library and existing
dependencies are insufficient, plus the maintenance, transitive-dependency,
licensing, platform, and binary impact where relevant.

Routine transitive or `go.sum` effects resulting from an explicitly approved
dependency change are allowed. Do not perform broad dependency upgrades as
incidental cleanup.

## Testing and validation

Tests must be proportional to the behavior changed.

- Add focused unit tests for planning, resolution, validation, and other pure
  logic.
- Add regression tests for bugs.
- Use `t.TempDir` and `t.Setenv` for isolated filesystem and home-directory
  state.
- Do not let host-side tests use the developer's real `HOME`, `~/.mdev`,
  `/Volumes`, Homebrew installation, tool caches, services, or managed data.
- Prefer fakes or injected runners for subprocess behavior.
- Add command-level coverage when flags, arguments, output, confirmation, or
  exit behavior changes.
- Do not weaken or delete tests merely to make a change pass.
- Avoid brittle assertions tied to irrelevant implementation details.

Agents may use host-side development operations that do not execute `mdev` or
exercise real machine-level `mdev` behavior, including reading and editing
files, `git status`, `git diff`, `gofmt`, static analysis, compilation,
building, inspecting artifacts, pure isolated unit tests, and repository-safe
tooling.

Tests that intentionally exercise real machine-level `mdev` behavior must run
in the VM. If an agent cannot establish that a requested host-side operation
stays completely outside real `mdev` runtime behavior, it must not execute it
on the host. Use the VM or stop and ask.

## Hard rule: never execute mdev on the host

Agents MUST NOT execute any form of the `mdev` application directly against
the developer's host machine.

This includes, but is not limited to:

```sh
mdev ...
./mdev ...
./dist/mdev ...
go run . ...
go run ./... ...
```

The prohibition applies even when a command appears harmless, including
`--version`, `version`, `--help`, `list`, `graph`, and `doctor`. Agents must not
decide that a particular `mdev` command is safe enough to run on the host.
Only the human developer may choose to execute `mdev` on the host.

Building `dist/mdev` on the host is allowed. Executing that binary on the host
is not. Permission to compile never implies permission to execute.

Whenever real CLI or runtime validation is required, use the macOS VM boundary
and the repository tooling under `scripts/vm/`, for example:

```sh
./scripts/vm/run.sh mdev --version
```

The VM scripts expect an appropriate already-built binary shared into the VM.
Build the macOS ARM64 binary first when the task requires it.

Before destructive or machine-level validation, state the exact command and
expected effect. If the VM is unavailable, the target or effect is ambiguous,
or the safety boundary cannot be established, stop and ask. Do not weaken this
rule for convenience or speed.

## Filesystem and deletion safety

Before moving, replacing, unlinking, or recursively deleting data:

- resolve and validate the exact target;
- ensure it is within the configured mdev-managed root when that is the
  intended ownership boundary;
- reject empty, root, home, volume-root, or otherwise broad paths;
- preserve unrelated existing contents;
- retain established confirmation behavior;
- add tests for containment and failure cases.

Never introduce destructive behavior based solely on an unchecked environment
variable, glob, unresolved symlink, or empty configuration value.

## Autonomous implementation authority

Inside an approved task, agents may autonomously:

- investigate the relevant code and behavior;
- fix compile errors, formatting, imports, and narrowly related static-analysis
  issues caused by their changes;
- add focused tests and test seams required to verify the requested behavior;
- fix a directly encountered defect when the intended behavior is
  unambiguous, backward-compatible, local to the task, and necessary for
  completion;
- improve local error context and naming without changing public behavior;
- make small refactors within touched packages when required for a clean,
  testable solution;
- update documentation directly affected by the implementation.

Keep autonomous fixes minimal and report them in the handoff.

Prior human approval is required for:

- new commands, flags, tools, services, configuration keys, or product
  defaults;
- changes to command semantics, output contracts, confirmation, or exit
  status;
- changes to the `Tool` or `Service` interfaces or registry model;
- package moves, new architectural layers, or broad cross-package refactors;
- configuration migration or storage/symlink layout changes;
- changes to dependency direction or install/uninstall ordering;
- adding, replacing, or materially upgrading a direct dependency;
- destructive behavior beyond the task's explicitly approved scope;
- automatic repair, migration, telemetry, network services, or background
  work;
- CI, release, signing, installation, distribution, or publication changes;
- broad cleanup unrelated to the task.

## Stop and ask

Stop and request a human decision when:

- requirements conflict or have more than one materially different product
  interpretation;
- completion requires choosing a new command, flag, default, config schema,
  storage layout, tool version, or migration policy;
- preserving compatibility conflicts with the requested result;
- a fix would change behavior for existing installations or delete or move
  user data;
- the correct package boundary is unclear and the alternatives require broad
  changes;
- a new direct dependency or core interface change appears necessary;
- real runtime validation is required but the VM is unavailable;
- the exact destructive target cannot be validated;
- pre-existing user changes overlap required edits and cannot safely be
  preserved;
- repository documentation and code disagree in a way that materially affects
  the implementation;
- a security, data-loss, or irreversible side effect is discovered outside the
  approved task;
- the decision would establish product direction, architecture, dependency
  adoption, data migration policy, destructive scope, or release policy.

When asking, provide evidence, viable options, compatibility and safety impact,
and a recommendation. Do not choose a product or architecture direction
silently.

## Feature branch and PR workflow

- One human-approved `product/work/<command>.md` corresponds to one feature
  branch and one pull request.
- Name the branch `feature/<command>`.
- The human must review and approve the work plan before implementation starts.
- Agents may break TODO items into smaller implementation steps.
- Prefer multiple meaningful commits where useful; do not force one TODO to
  equal one commit.
- Agents may push the feature branch regularly to preserve work remotely.
- When the work-file status reaches `done`, open the pull request.
- Address pull-request feedback on the same branch and pull request.
- Agents must never merge their own pull requests. The human owns approval and
  merging to `main`.
- Never work directly on or push directly to `main`.

## Git behavior

- Inspect `git status` before editing and before handoff.
- Preserve all pre-existing staged, unstaged, and untracked work; do not revert
  or overwrite it.
- Within an approved feature workflow, follow the branch and pull-request rules
  above. Otherwise, do not stage, commit, amend, rebase, merge, push, publish,
  or create a pull request unless explicitly requested.
- Do not use destructive Git commands.
- Keep changes limited to files required by the task.
- If asked to commit, include only task-owned changes and follow the required
  repository commit-message format.
- Never commit secrets, credentials, VM keys, machine-specific paths, or local
  SSH configuration.

## Documentation

Update `README.md` or other relevant documentation when implementation changes
commands, flags, examples, output expectations, configuration, storage,
supported tools or services, prerequisites, or build/test/VM workflows.

Documentation must describe implemented behavior. Do not document speculative
architecture or unfinished TODOs as current capability. Record meaningful
behavioral constraints and decisions rather than leaving them only in code.

## Definition of done

A task is complete only when:

- the requested behavior is implemented without unrelated scope expansion;
- relevant success, failure, regression, and safety-sensitive tests exist;
- changed Go files are formatted;
- permitted host-side checks pass;
- required real CLI or OS-facing behavior is validated through the VM;
- compatibility is preserved, or approved changes are tested and documented;
- relevant documentation is updated;
- no new panic, ignored error, nondeterministic public output, or unsafe path
  operation is introduced;
- `git diff` contains only intended task changes and pre-existing work remains
  intact;
- the handoff states what changed, why, files touched, validation performed and
  not performed, risks, and intentionally deferred follow-ups.

If validation cannot be completed, report the exact limitation. Do not claim
the task is fully validated.
