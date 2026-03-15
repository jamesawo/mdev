package interactive

import "github.com/AlecAivazis/survey/v2"

// RadioSelect presents options as radio buttons and lets the user make a selection
func RadioSelect(title string, options []string) (int, error) {
	var selected string

	prompt := &survey.Select{
		Message: title,
		Options: options,
	}

	err := survey.AskOne(prompt, &selected)
	if err != nil {
		return -1, err
	}

	for i, opt := range options {
		if opt == selected {
			return i, nil
		}
	}

	return -1, nil
}
