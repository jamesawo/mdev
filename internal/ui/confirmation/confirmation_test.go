package confirmation

import (
	"bytes"
	"testing"
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

func TestConfirmation_WhenInputIsClosed_Declines(t *testing.T) {
	confirmer, _ := givenConfirmer("", false)

	accepted := confirmer.Ask("Continue?")

	assertConfirmationDeclined(t, accepted)
}

func TestConfirmation_DefaultNoRequiresExplicitYes(t *testing.T) {
	for _, answer := range []string{"\n", "no\n"} {
		confirmer, _ := givenConfirmer(answer, false)
		assertConfirmationDeclined(t, confirmer.AskDefaultNo("Apply?"))
	}
	confirmer, _ := givenConfirmer("yes\n", false)
	assertConfirmationAccepted(t, confirmer.AskDefaultNo("Apply?"))
}

func givenConfirmer(answer string, assumeYes bool) (*Confirmer, *bytes.Buffer) {
	output := &bytes.Buffer{}
	confirmer := New(bytes.NewBufferString(answer), output, assumeYes)

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
