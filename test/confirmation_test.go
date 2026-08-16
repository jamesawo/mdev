package test

import (
	"bytes"
	"testing"

	"github.com/jamesawo/mdev/internal/ui/confirmation"
)

func TestConfirmation_WhenAssumeYesIsEnabled_AcceptsWithoutPrompting(t *testing.T) {
	confirmer, output := givenConfirmer("no\n", true)

	accepted := confirmer.Ask("Continue?")

	assertConfirmationAccepted(t, accepted)
	assertNoPromptWritten(t, output)
}

func TestConfirmation_WhenUserAnswersNo_Declines(t *testing.T) {
	confirmer, _ := givenConfirmer("no\n", false)

	accepted := confirmer.Ask("Continue?")

	assertConfirmationDeclined(t, accepted)
}

func givenConfirmer(answer string, assumeYes bool) (*confirmation.Confirmer, *bytes.Buffer) {
	output := &bytes.Buffer{}
	confirmer := confirmation.New(bytes.NewBufferString(answer), output, assumeYes)

	return confirmer, output
}

func assertConfirmationAccepted(t *testing.T, accepted bool) {
	t.Helper()

	if !accepted {
		t.Fatal("expected confirmation to be accepted")
	}
}

func assertConfirmationDeclined(t *testing.T, accepted bool) {
	t.Helper()

	if accepted {
		t.Fatal("expected confirmation to be declined")
	}
}

func assertNoPromptWritten(t *testing.T, output *bytes.Buffer) {
	t.Helper()

	if output.Len() != 0 {
		t.Fatalf("expected no confirmation prompt, got %q", output.String())
	}
}
