# list

## purpose

`mdev list` shows the tools mdev knows about and whether each tool is currently installed.
list should provide a quick, calm overview of what is available on the machine without requiring the user to understand
tool dependencies or mdev internals.
it is read-only. list must never install, uninstall, repair, configure, or otherwise modify tools.

## basic usage

run:

```bash
mdev list

show all tools registered with mdev, grouped into system tools and other tools.

for example:

system tools
  ✓ curl       installed
  ✓ git        installed
  ✓ homebrew   installed
tools
  ✓ docker     installed
  ✓ go         installed
  ○ node       missing
  ○ postgres   missing

a fresh mdev installation should still show all registered tools even when none of the optional tools are installed.

tool groups

show tools in two sections.

system tools

system tools are foundational tools or prerequisites used by mdev or required by other tools.

examples may include:

curl
git
homebrew

tools

all other registered tools appear under:

tools

do not call this section development tools. mdev may support tools for broader use cases in the future.

if a section has no registered tools, omit the empty section.

installation status

list cares only whether a registered tool is installed.

it does not distinguish between:

* tools installed by mdev;
* tools installed manually;
* tools installed by Homebrew or another package manager;
* tools that existed before mdev was configured.

if mdev can verify that the tool is available, show:

✓ <tool>    installed

if the tool is not installed, show:

○ <tool>    missing

mdev does not claim ownership merely by reporting that a tool is installed.

install and uninstall commands are responsible for deciding how to act on the detected installation.

unknown status

if mdev cannot reliably determine whether a particular tool is installed because its verification encounters an unexpected failure, do not report the tool as missing.

show:

? <tool>    unknown

continue checking and displaying the remaining tools.

after the list, explain the failed status check with a concise actionable error.

for example:

could not determine postgres status: <reason>

if one or more tool statuses are unknown, list should exit as a failure after presenting the available results.

an unexpected verification failure must not prevent statuses already known for other tools from being useful to the user.

ordering

tools are displayed alphabetically within each section.

do not order tools according to dependency relationships.

dependency relationships belong to:

mdev graph

the output order must be deterministic.

output

use simple status symbols together with text.

normal states are:

✓ installed
○ missing
? unknown

do not rely on color alone to communicate status.

color may be used subtly when supported by the terminal, but the symbols and text must remain understandable without color.

do not use decorative emoji.

keep columns readable and consistently aligned where practical.

do not show additional information such as versions, installation paths, descriptions, or dependencies in the default list.

the purpose of the default output is a quick overview, not detailed tool inspection.

configuration

mdev list requires a valid mdev configuration.

if setup has not been completed, calmly direct the user to:

mdev setup

do not create configuration automatically.

if ~/.mdev/config.yaml is malformed or unreadable, leave it untouched and report the problem as a failure.

unavailable storage

if the configured mdev storage is unavailable, do not interpret that condition as tools being missing.

report that the configured storage is unavailable and show the expected path.

follow the established unavailable-storage behavior rather than silently falling back to another location, creating replacement storage, or modifying configuration.

list must remain read-only.

verification behavior

use each registered tool’s established verification behavior to determine whether it is installed.

list should not introduce separate installation-detection rules that can disagree with the rest of mdev.

verification must not modify the machine.

list must not attempt automatic repair when verification fails.

dependencies

list does not display dependency relationships or dependency ordering.

a tool may be shown as installed even if one of its dependencies is currently missing.

diagnosing inconsistent or unhealthy installations belongs to commands such as:

mdev doctor

dependency visualization belongs to:

mdev graph

read-only behavior

running list must not:

* install or uninstall tools;
* create tool storage;
* repair tools;
* modify configuration;
* modify shell configuration;
* start or stop services;
* change dependency state;
* prompt for destructive actions.

the command observes and reports current state only.

non-interactive behavior

mdev list does not require interactive input during normal operation.

it should be suitable for running in a normal terminal or with its output redirected.

when terminal capabilities such as color are unavailable, preserve the same information using plain text and status symbols.

list does not require --yes.

cli experience

all user-facing copy should follow mdev’s simple, calm, lowercase style.

the normal successful command should remain compact.

do not add progress indicators or spinners for ordinary fast verification.

if checking tool state takes noticeable time because of an individual tool’s verification mechanism, the command should still avoid unnecessary visual noise.

Cobra’s normal help system should provide concise and useful help.

errors should explain what failed without hiding successfully determined tool statuses when those results remain useful.

relationship with other commands

setup introduces list as the next step after first-time configuration:

see what's available:
mdev list

list provides discovery of all tools registered with mdev.

install acts on tools shown by list.

uninstall may act on an installed tool regardless of how that tool was originally installed.

graph explains dependency relationships.

doctor diagnoses unhealthy, inconsistent, or otherwise problematic installations.

list should not absorb the responsibilities of those commands.

testing

unit

cover smaller behavior such as:

* grouping system tools and other tools;
* alphabetical deterministic ordering;
* installed status;
* missing status;
* unknown status;
* continuation after an individual verification failure;
* failure exit status when verification is unknown;
* empty-section omission;
* missing configuration;
* malformed configuration;
* unavailable configured storage;
* read-only behavior where relevant;
* output formatting and status representation.

tests should use the registered tool abstractions rather than depending on the developer’s actual installed tools.

e2e

the list E2E test exercises the real user journey using the compiled mdev binary inside the macOS VM.

the happy path should confirm that a configured user can run:

mdev list

and receive the expected grouped tool overview and installation statuses.

E2E testing should validate the real CLI experience rather than duplicating every verification and formatting edge case covered by unit/component tests.

done when

list is done when:

* all registered tools are shown;
* system tools and other tools are clearly separated;
* tools are alphabetically ordered within each section;
* installed tools are reported as installed regardless of how they were installed;
* missing tools are reported as missing;
* verification failures are reported as unknown rather than incorrectly as missing;
* one failed verification does not hide useful results for other tools;
* the command remains completely read-only;
* configuration and unavailable-storage behavior are consistent with setup;
* dependencies are left to graph and health diagnosis is left to doctor;
* user-facing copy has been reviewed;
* relevant unit/component tests pass;
* the happy-path E2E journey passes in the macOS VM.

One thing I deliberately preserved from `setup`: the spec describes **product behavior**, not how Codex should implement it. That should make the later work-plan comparison much cleaner.