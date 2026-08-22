// Package readiness coordinates reusable system-prerequisite checks and remediation.
package readiness

import (
	"context"
	"fmt"

	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
)

// Item is one checked prerequisite and its current readiness state.
type Item struct {
	Prerequisite prerequisites.Prerequisite
	State        prerequisites.State
	Detail       string
}

// Ready reports whether this prerequisite needs no work.
func (i Item) Ready() bool { return i.State == prerequisites.StateReady }

// Remediable reports whether mdev has an approved remediation capability.
func (i Item) Remediable() bool {
	_, ok := i.Prerequisite.(prerequisites.Remediator)
	return ok
}

// RemediationDescription returns the concrete approved change, when available.
func (i Item) RemediationDescription() string {
	remediator, ok := i.Prerequisite.(prerequisites.Remediator)
	if !ok {
		return "manual remediation required"
	}
	return remediator.RemediationDescription()
}

// Reporter observes real readiness boundaries. Implementations may write plain text.
type Reporter interface {
	Checking(string) error
	Checked(Item) error
	Remediating(string, string) error
	Verified(Item) error
}

type noopReporter struct{}

func (noopReporter) Checking(string) error            { return nil }
func (noopReporter) Checked(Item) error               { return nil }
func (noopReporter) Remediating(string, string) error { return nil }
func (noopReporter) Verified(Item) error              { return nil }

// CheckAll checks the complete registered prerequisite plan.
func CheckAll(ctx context.Context, reporter Reporter) ([]Item, error) {
	return Check(ctx, prerequisites.List(), reporter)
}

// CheckNames checks the dependency closure for canonical prerequisite names.
func CheckNames(ctx context.Context, names []string, reporter Reporter) ([]Item, error) {
	selected := make([]prerequisites.Prerequisite, 0, len(names))
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true
		prerequisite, ok := prerequisites.Get(name)
		if !ok {
			return nil, fmt.Errorf("unknown system prerequisite %s", name)
		}
		selected = append(selected, prerequisite)
	}
	return Check(ctx, selected, reporter)
}

// Check inspects a deterministic dependency-first prerequisite closure.
func Check(ctx context.Context, selected []prerequisites.Prerequisite, reporter Reporter) ([]Item, error) {
	if reporter == nil {
		reporter = noopReporter{}
	}
	plan, err := prerequisites.Resolve(selected)
	if err != nil {
		return nil, err
	}
	items := make([]Item, 0, len(plan))
	for _, prerequisite := range plan {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := reporter.Checking(prerequisite.Name()); err != nil {
			return nil, err
		}
		state, detail, err := prerequisites.ReadinessStatus(ctx, prerequisite)
		if err != nil {
			return nil, fmt.Errorf("check system prerequisite %s: %w", prerequisite.Name(), err)
		}
		item := Item{Prerequisite: prerequisite, State: state, Detail: detail}
		items = append(items, item)
		if err := reporter.Checked(item); err != nil {
			return nil, err
		}
	}
	return items, nil
}

// Unready returns the items which still require remediation or manual recovery.
func Unready(items []Item) []Item {
	var result []Item
	for _, item := range items {
		if !item.Ready() {
			result = append(result, item)
		}
	}
	return result
}

// Remediate performs approved changes dependency-first and verifies each result.
func Remediate(ctx context.Context, items []Item, reporter Reporter) error {
	if reporter == nil {
		reporter = noopReporter{}
	}
	for _, item := range items {
		if item.Ready() {
			continue
		}
		remediator, ok := item.Prerequisite.(prerequisites.Remediator)
		if !ok {
			return fmt.Errorf("system prerequisite %s is %s and requires manual remediation", item.Prerequisite.Name(), item.State)
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := reporter.Remediating(item.Prerequisite.Name(), remediator.RemediationDescription()); err != nil {
			return err
		}
		if err := remediator.RemediateContext(ctx); err != nil {
			return fmt.Errorf("remediate system prerequisite %s: %w", item.Prerequisite.Name(), err)
		}
		if verifier, ok := item.Prerequisite.(prerequisites.Verifier); ok {
			if err := verifier.VerifyContext(ctx); err != nil {
				return fmt.Errorf("verify system prerequisite %s: %w", item.Prerequisite.Name(), err)
			}
		}
		state, detail, err := prerequisites.ReadinessStatus(ctx, item.Prerequisite)
		if err != nil {
			return fmt.Errorf("verify system prerequisite %s: %w", item.Prerequisite.Name(), err)
		}
		verified := Item{Prerequisite: item.Prerequisite, State: state, Detail: detail}
		if !verified.Ready() {
			return fmt.Errorf("verify system prerequisite %s: state is %s", item.Prerequisite.Name(), state)
		}
		if err := reporter.Verified(verified); err != nil {
			return err
		}
	}
	return nil
}
