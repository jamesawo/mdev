# uninstall

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
