# doctor

## purpose

`mdev doctor` diagnoses the current health of mdev configuration, system
readiness, registered tools, and relevant environment state. it is a recovery
and troubleshooting command, not a required step between setup and install.

doctor reuses the same readiness checks used by setup. it does not own a second
copy of prerequisite rules, and setup does not execute the doctor Cobra command.

## diagnosis

doctor reports real check results progressively and distinguishes ready,
missing, outdated, broken, and unknown/error states where supported. status
errors are not silently converted to missing or healthy results.

doctor should detect incomplete setup and recommend `mdev setup`. after a
machine changes, it should identify readiness drift and provide actionable next
steps.

## remediation

when remediation is explicitly requested, doctor may use the same reusable
readiness remediation operations as setup. system changes require informed
consent, run progressively, preserve errors, and are verified afterward.

doctor performs no broad rollback. cancellation or failure stops future work,
preserves completed changes, and reports unresolved readiness without claiming
success.

## architecture

doctor consumes a shared readiness service below the CLI command layer:

```text
shared readiness checks
      /          \
   setup        doctor
prepare/verify  inspect/report
```

the doctor command remains thin presentation and workflow wiring. prerequisite
inspection, dependency ordering, remediation, verification, and context-aware
execution belong in the reusable readiness/prerequisite boundary.

## testing

cover diagnosis without mutation, consent, successful and failed remediation,
verification, cancellation, deterministic output, and parity with the checks
used by setup. real machine-level journeys run only in the macOS VM.
