# list

## purpose

`mdev list` shows the tools mdev knows about and whether each tool is currently installed.
List should provide a quick, calm overview of what is available on the machine without requiring the user to understand
tool dependencies or mdev internals.
It is read-only. List must never install, uninstall, repair, configure, or otherwise modify tools.

## basic usage

Run:

```bash
mdev list
```

Show all tools registered with mdev, grouped into system tools and other tools.

For example:

```text
system tools
  ✓ curl       installed
  ✓ git        installed
  ✓ homebrew   installed
tools
  ✓ docker     installed
  ✓ go         installed
  ○ node       missing
  ○ postgres   missing
```

A fresh mdev installation should still show all registered tools even when none of the optional tools are installed.

## tool groups

Show tools in two sections.

### system tools

System tools are foundational tools or prerequisites used by mdev or required by other tools.

Examples may include:

- curl
- git
- homebrew

### tools

All other registered tools appear under:

```text
tools
```

Do not call this section development tools. mdev may support tools for broader use cases in the future.

If a section has no registered tools, omit the empty section.

## installation status

List cares only whether a registered tool is installed.

It does not distinguish between:

- tools installed by mdev;
- tools installed manually;
- tools installed by Homebrew or another package manager;
- tools that existed before mdev was configured.

If mdev can verify that the tool is available, show:

```text
✓ <tool>    installed
```

If the tool is not installed, show:

```text
○ <tool>    missing
```

mdev does not claim ownership merely by reporting that a tool is installed.

Install and uninstall commands are responsible for deciding how to act on the detected installation.

## unknown status

If mdev cannot reliably determine whether a particular tool is installed because its verification encounters an
unexpected failure, do not report the tool as missing.

Show:

```text
? <tool>    unknown
```

Continue checking and displaying the remaining tools.

After the list, explain the failed status check with a concise actionable error.

For example:

```text
could not determine postgres status: <reason>
```

If one or more tool statuses are unknown, list should exit as a failure after presenting the available results.

An unexpected verification failure must not prevent statuses already known for other tools from being useful to the
user.

## ordering

Tools are displayed alphabetically within each section.

Do not order tools according to dependency relationships.

Dependency relationships belong to:

```bash
mdev graph
```

The output order must be deterministic.

## output

Use simple status symbols together with text.

Normal states are:

- ✓ installed
- ○ missing
- ? unknown

Do not rely on color alone to communicate status.

Color may be used subtly when supported by the terminal, but the symbols and text must remain understandable without
color.

Do not use decorative emoji.

Keep columns readable and consistently aligned where practical.

Do not show additional information such as versions, installation paths, descriptions, or dependencies in the default
list.

The purpose of the default output is a quick overview, not detailed tool inspection.

For normal human-oriented output, render each section and tool result progressively as verification proceeds so the user
receives immediate and continuing feedback during potentially slow checks. Process entries sequentially in deterministic
alphabetical display order within each section, rather than printing in completion order.

Progressive output must not introduce concurrency, artificial delays, spinners, animation, or unnecessary progress
messages. Unknown verification errors are still summarized after the completed overview.

## json output

`mdev list --json` produces machine-readable JSON while `mdev list` remains the human-oriented default.

JSON output must use the same underlying status checks and results as normal output. Do not implement a separate
status-detection path.

The JSON document preserves the `system_tools` and `tools` groups. Each entry contains at least:

- `name`;
- `status`, with one of `installed`, `missing`, or `unknown`.

An entry with an `unknown` status also contains a concise `error` field explaining its failed status check.

JSON standard output must contain valid JSON only. Do not include headings, status symbols, colors, progress messages,
explanatory text, or any other human-oriented output.

Unlike normal output, JSON is emitted as one complete document after every verification finishes. Ordering remains
deterministic within both groups.

If one or more statuses are unknown, emit the complete JSON document with all results and then return a non-zero exit
status. A configuration or storage failure that prevents listing also returns non-zero and must not write malformed or
partial JSON.

## configuration

`mdev list` requires a valid mdev configuration.

If setup has not been completed, calmly direct the user to:

```bash
mdev setup
```

Do not create configuration automatically.

If `~/.mdev/config.yaml` is malformed or unreadable, leave it untouched and report the problem as a failure.

## unavailable storage

If the configured mdev storage is unavailable, do not interpret that condition as tools being missing.

Report that the configured storage is unavailable and show the expected path.

Follow the established unavailable-storage behavior rather than silently falling back to another location, creating
replacement storage, or modifying configuration.

List must remain read-only.

## verification behavior

Use each registered tool's established verification behavior to determine whether it is installed.

List should not introduce separate installation-detection rules that can disagree with the rest of mdev.

Verification must not modify the machine.

List must not attempt automatic repair when verification fails.

## dependencies

List does not display dependency relationships or dependency ordering.

A tool may be shown as installed even if one of its dependencies is currently missing.

Diagnosing inconsistent or unhealthy installations belongs to commands such as:

```bash
mdev doctor
```

Dependency visualization belongs to:

```bash
mdev graph
```

## read-only behavior

Running list must not:

- install or uninstall tools;
- create tool storage;
- repair tools;
- modify configuration;
- modify shell configuration;
- start or stop services;
- change dependency state;
- prompt for destructive actions.

The command observes and reports current state only.

## non-interactive behavior

`mdev list` does not require interactive input during normal operation.

It should be suitable for running in a normal terminal or with its output redirected.

When terminal capabilities such as color are unavailable, preserve the same information using plain text and status
symbols.

List does not require `--yes`.

## cli experience

All user-facing copy should follow mdev's simple, calm, lowercase style.

The normal successful command should remain compact.

Do not add progress indicators or spinners for ordinary fast verification.

If checking tool state takes noticeable time because of an individual tool's verification mechanism, the command should
still avoid unnecessary visual noise.

Cobra's normal help system should provide concise and useful help.

Errors should explain what failed without hiding successfully determined tool statuses when those results remain useful.

## relationship with other commands

Setup introduces list as the next step after first-time configuration:

```text
see what's available:
mdev list
```

List provides discovery of all tools registered with mdev.

Install acts on tools shown by list.

Uninstall may act on an installed tool regardless of how that tool was originally installed.

Graph explains dependency relationships.

Doctor diagnoses unhealthy, inconsistent, or otherwise problematic installations.

List should not absorb the responsibilities of those commands.

## testing

### unit

Cover smaller behavior such as:

- grouping system tools and other tools;
- alphabetical deterministic ordering;
- installed status;
- missing status;
- unknown status;
- continuation after an individual verification failure;
- failure exit status when verification is unknown;
- empty-section omission;
- missing configuration;
- malformed configuration;
- unavailable configured storage;
- read-only behavior where relevant;
- output formatting and status representation.
- progressive human output before later verification completes;
- JSON grouping, fields, statuses, unknown error details, and deterministic ordering;
- complete valid JSON before an unknown-status failure;
- no partial JSON for configuration or storage failures;
- the normal and JSON modes using the same status-detection path.

Tests should use the registered tool abstractions rather than depending on the developer's actual installed tools.

### e2e

The list E2E test exercises the real user journey using the compiled mdev binary inside the macOS VM.

The happy path should confirm that a configured user can run:

```bash
mdev list
```

and receive the expected grouped tool overview and installation statuses.

E2E testing should validate the real CLI experience rather than duplicating every verification and formatting edge case
covered by unit/component tests.

## done when

List is done when:

- all registered tools are shown;
- system tools and other tools are clearly separated;
- tools are alphabetically ordered within each section;
- installed tools are reported as installed regardless of how they were installed;
- missing tools are reported as missing;
- verification failures are reported as unknown rather than incorrectly as missing;
- one failed verification does not hide useful results for other tools;
- normal output is rendered progressively in deterministic display order;
- `--json` emits one complete valid, deterministically ordered document from the same underlying tool states;
- unknown JSON entries include errors and still produce a complete document before the command fails;
- the command remains completely read-only;
- configuration and unavailable-storage behavior are consistent with setup;
- dependencies are left to graph and health diagnosis is left to doctor;
- user-facing copy has been reviewed;
- relevant unit/component tests pass;
- the happy-path E2E journey passes in the macOS VM.
