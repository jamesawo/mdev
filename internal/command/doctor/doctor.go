package doctor

import (
	"context"
	"io"
)

// Execute runs the doctor command flow.
func Execute(isFixFlag bool) error {
	return ExecuteContext(context.Background(), isFixFlag)
}

// ExecuteWithStreams runs doctor using command-owned streams for interactive fixes.
func ExecuteWithStreams(ctx context.Context, isFixFlag bool, in io.Reader, out io.Writer) error {
	if isFixFlag {
		return FixContextWithStreams(ctx, in, out)
	}
	return ExecuteContext(ctx, false)
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
