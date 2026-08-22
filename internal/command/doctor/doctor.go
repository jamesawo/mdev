package doctor

import "context"

// Execute runs the doctor command flow.
func Execute(isFixFlag bool) error {
	return ExecuteContext(context.Background(), isFixFlag)
}

// ExecuteContext runs diagnosis or remediation with cancellation.
func ExecuteContext(ctx context.Context, isFixFlag bool) error {

	if isFixFlag {
		return FixContext(ctx)
	}

	reporter := &progressReporter{}
	report, err := RunContext(ctx, reporter)

	if err != nil {
		return err
	}

	renderSummary(report)
	return nil
}
