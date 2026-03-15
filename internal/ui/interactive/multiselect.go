package interactive

import (
	"github.com/AlecAivazis/survey/v2"
)

// todo: fix or replace this component, it does not tell user how to quit

// MultiSelect presents a checkbox selection list
// and returns the selected option values.
//
// Example usage:
// options := []string{"sdkman","java","maven"}
// selected, _ := MultiSelect("Select tools", options)
func MultiSelect(message string, options []string) ([]string, error) {

	var selected []string

	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}

	err := survey.AskOne(prompt, &selected)

	return selected, err
}
