package opsgenie

import (
	"fmt"
	"time"
)

func validateDate(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	layoutStr := "2006-01-02T15:04:05Z"
	_, err := time.Parse(layoutStr, value)
	if err != nil {
		errors = append(errors, fmt.Errorf("you entererd wrong date-time format. Its should be like this : %s", layoutStr))
	}

	return
}
