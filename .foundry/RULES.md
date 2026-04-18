# RULES.md

## Purpose

These rules govern all AI-assisted changes in this repository.  
They exist to keep delivery reliable, traceable, and maintainable.

Agents must follow these rules strictly.

---

## 1. Scope Control

Only implement the approved task scope defined in the task spec.

Do not add extra features, refactors, or opinionated changes unless explicitly requested.

If something important is missing, raise it in notes instead of silently expanding scope.

---

## 2. Minimal Change Principle

Prefer the smallest clean change that solves the problem.

Do not rewrite working code without clear benefit.

Avoid broad churn across unrelated files.

---

## 3. Preserve Stability

Do not break existing behavior unless the spec explicitly requires it.

Backward compatibility is the default.

---

## 4. Understand Before Changing

Read relevant files before editing.

Understand existing patterns, structure, and intent.

Do not guess architecture.

---

## 5. Follow Existing Conventions

Match the repository’s current style unless instructed otherwise.

Prefer consistency over personal preference.

---

## 6. Traceability Required

Every change must map to a Foundry task ID.

Use the task ID in:
- branch names
- commits
- pull requests
- notes

Example:

FD-001

---

## 7. Ask When Unclear

If requirements are ambiguous or conflicting, stop and ask for clarification.

Do not invent missing requirements.

---

## 8. Validate Changes

Before submitting work:

- run relevant tests
- build successfully
- check formatting
- review changed files

Do not submit avoidable breakage.

---

## 9. Keep Code Maintainable

Prefer readable solutions over clever ones.

Avoid unnecessary abstractions.

Avoid premature optimization.

---

## 10. Respect Sensitive Areas

Be cautious when modifying:

- build scripts
- CI/CD
- dependency versions
- file structure
- data deletion logic
- authentication/secrets
- public interfaces

If touched, explain why.

---

## 11. Document Meaningful Decisions

If a change introduces new behavior, constraints, or commands, update docs.

Do not leave hidden knowledge in code only.

---

## 12. Leave Clear Handover Notes

When task is complete, provide:

- what changed
- why
- files touched
- risks
- follow-up suggestions

---

## Final Rule

When uncertain, choose safety, clarity, and minimalism.