# CONVENTIONS
## Purpose
These conventions define how code should be written, structured, and evolved in this repository.
 - Rules protect quality.  
 - Conventions create consistency.

All contributors and agents must follow these standards unless a task explicitly overrides them.

# 1. Core Engineering Standard
Build software that is:
- simple to understand
- easy to change
- safe to extend
- pleasant to use
- clear to maintain
  
Prefer mature engineering judgment over cleverness.
  If two solutions work, choose the simpler one.

# 2. Architecture Direction
This project uses layered boundaries.
```text
CLI Layer
→ Application Layer
→ Domain / Use Case Layer
→ Infrastructure Layer
```

## Layer Responsibilities

### CLI Layer

Responsible only for:

* parsing commands
* reading flags / args
* invoking application handlers
* rendering output

CLI commands must not contain business logic.

Example:
```bash
mdev doctor
```
Cobra command receives input and delegates immediately.


### Application Layer

Coordinates use cases.

Examples:
```bash
RunDoctor()
InstallTool()
SetupDisk()
```

This layer orchestrates workflow and dependencies.


### Domain / Use Case Layer

Contains core logic and decisions.

Examples:

- validate system state
- determine missing tools
- choose install steps

Must be testable and independent from Cobra or shell details.


### Infrastructure Layer

Handles external systems.

Examples:

* Homebrew
* shell execution
* filesystem
* OS detection
* network calls

External tools must be wrapped behind interfaces or adapters.


# 3. Command Design

Commands must be:

* clear
* memorable
* minimal
* predictable

Examples:
```bash
mdev doctor
mdev install go
mdev disk setup
```

Avoid command sprawl.

Prefer evolving existing commands over adding noisy ones.

New commands should be added with minimal impact to existing code.


# File and Naming Conventions

Use explicit names.

Good: ✅

* doctor_command.go
* doctor_service.go
* brew_client.go
* output_stream.go

Bad: ❌

* utils.go
* helper.go
* misc.go
* manager.go

Names must reveal responsibility.


# Function Design

Functions should:

* do one thing
* have clear inputs
* return useful outputs/errors
* remain small where practical
* read top-to-bottom logically

Prefer readable code over compressed code.

Code should read like a story.


# Output / User Experience Standard

CLI UX is a product feature.

Messages must be:

* clear
* concise
* calm
* useful
* human-readable

Avoid noisy logs and cryptic errors.

Good: ✅

Checking Homebrew...
Homebrew found.
Checking Go...
Go not installed.

Bad: ❌

exec failed code=127 subprocess error

Translate technical failures into useful user guidance.


# Output Text Centralization

User-facing text should be defined in a structured output/messages package.

Avoid scattering CLI strings across files.

Example:

internal/output/messages/

Benefits:

* consistency
* easier edits
* localization later
* testing easier


# Streaming Output Model

Long-running operations should show progress progressively.

Users should see movement before completion.

### Preferred flow:

Starting doctor checks...
Checking Homebrew...
Checking Go...
Checking Xcode tools...
Doctor complete.

Implementation must support flexible output granularity:

* line-based today
* token/chunk streaming later if needed

Output transport must be decoupled from business logic.

# Error Handling

Return contextual errors.

Good: ✅

Homebrew check failed: brew command not found

Bad: ❌

failed

Never swallow errors silently.


# Testing Standard

Tests are required for meaningful logic.

Use pragmatic layers:

- Unit Tests 
  - For pure logic and use cases.

- Integration Tests 
  - For shell, brew, filesystem adapters.

- Command Tests 
  - For CLI behavior and output. 
  - Prefer fast tests first. 
  - Avoid brittle tests tied to irrelevant implementation detail.

# External Tool Boundaries

Interactions with external tools must be isolated.

Examples:

* brew
* shell commands
* macOS APIs

Never spread raw shell execution across the codebase.

Use dedicated adapters.

Example:

internal/platform/brew/
internal/platform/shell/
internal/platform/os/

# Multi-Platform Future Readiness

Current priority is macOS.

But architecture should allow future support for:

* Linux
* Windows

Use OS-specific adapters behind stable interfaces.

Do not hardcode macOS assumptions into core logic.


# Dependency Discipline

Add dependencies only when justified.

Prefer standard library when reasonable.

Avoid heavy frameworks for small problems.


# Onboarding Standard

A new engineer should understand the project quickly.

Codebase must feel navigable.

Prefer obvious structure over personal genius.

If only the author understands it, it is bad design.


# Refactoring Standard

Improve code in scoped increments.

Do not rewrite large areas casually.

Use Foundry tasks to modernize command-by-command over time.


# Definition of Good Code

Good code is:

* readable
* intentional
* testable
* maintainable
* boring in the best way


