package validation

import "errors"

func RequireFields(fields map[string]string) error {
	for k, v := range fields {
		if v == "" {
			return errors.New("missing required field: " + k)
		}
	}
	return nil
}
