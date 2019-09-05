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
		errors = append(errors, fmt.Errorf("wrong date-time format:%s, Its should be like this : %s", value, layoutStr))
	}

	return
}
