package opsgenie

import (
	"fmt"
	"time"
)

func validateDateWithMinutes(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	layoutStr := "2006-01-02T15:04:05Z"
	selectedTime, err := time.Parse(layoutStr, value)
	if err != nil {
		errors = append(errors, fmt.Errorf("wrong date-time format:%s, Its should be like this : %s", value, layoutStr))
		return
	}
	min := selectedTime.Minute()
	if min%30 != 0 {
		errors = append(errors, fmt.Errorf("you can only select 30 or 00 for minutes. you enter : %d", selectedTime.Minute()))
	}

	return
}

func validateDate(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	layoutStr := "2006-01-02T15:04:05Z"
	_, err := time.Parse(layoutStr, value)
	if err != nil {
		errors = append(errors, fmt.Errorf("wrong date-time format:%s, Its should be like this : %s", value, layoutStr))
		return
	}

	return
}
