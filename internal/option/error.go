package option

import "strings"

type optionErrors []error

func (e *optionErrors) Error() string {
	messages := make([]string, 0, len(*e))
	for _, err := range *e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "\n")
}
