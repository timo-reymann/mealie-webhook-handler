package validation

import "fmt"

type Check func(map[string]string) error

func RequireKey(key string) Check {
	return func(config map[string]string) error {
		_, exists := config[key]
		if !exists {
			return fmt.Errorf("missing required config key '%s'", key)
		}
		return nil
	}
}

func FailOnFirst(checks ...Check) Check {
	return func(config map[string]string) error {
		for _, check := range checks {
			err := check(config)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
