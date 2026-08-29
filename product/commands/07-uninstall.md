# uninstall

## purpose and presentation

`mdev uninstall <tool>` removes a registered tool and any approved installed
dependants in safe reverse dependency order. Tool implementations own removal
mechanics. The uninstall workflow owns planning, consent, managed-storage
cleanup, and the user experience.

Normal output presents mdev lifecycle work rather than provider mechanics:

```text
uninstalling podman-compose... ✓
uninstalling podman... ✓
cleaning podman storage... ✓
podman-compose removed.
podman removed.
```

Successful Homebrew, SDKMAN, NVM, and other subprocess stdout or stderr is
captured and does not stream into normal output. A failure closes the active
phase without a success mark and returns a concise error identifying the tool
and operation with bounded useful provider diagnostics. Successful provider
chatter before a failure is not dumped merely because the final command fails.

The same context-scoped execution boundary used by install applies to
uninstall. It preserves cancellation and affects only workflows that opt into
managed presentation; setup and doctor retain their approved behavior. The
core `tools.Tool` contract remains compatible, with optional context-aware
uninstall capability used where a provider subprocess must receive caller
cancellation.

Dependency warnings, plans, directory previews, prompts, progress, and final
results use lowercase mdev product language and deterministic ordering.
Provider-internal dependencies are never presented as mdev tools. A tool is
reported removed only after its uninstall operation and any owned managed
storage cleanup succeed.

## podman ownership

The registered Podman tools have independent uninstall boundaries:

- `podman` owns the Homebrew Podman CLI, the registered managed machine, and
  machine data below `<storage_path>/podman`;
- `podman-desktop` owns only a Podman Desktop cask installed through mdev;
- `podman-compose` owns only the Homebrew Compose provider installed through
  mdev.

Because `podman-desktop` and `podman-compose` depend on `podman`, uninstalling
`podman` must warn about installed dependants and remove approved dependants
before the base tool. Removing either add-on independently must not remove the
Podman CLI, machine, or the other add-on.

mdev must not remove an application or package it does not own. In particular,
discovering an existing Podman Desktop application must not make a later mdev
uninstall delete a Desktop installation that Homebrew does not identify as its
cask. Existing Podman machines and managed data retain the established
confirmation and safe-storage behavior.
