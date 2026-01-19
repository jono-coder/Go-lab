package validate

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+$`)
)

func NotEmpty(field string, value string) error {
	if value == "" {
		return fmt.Errorf("'%v' may not be empty", field)
	}
	return nil
}

func NotBlank(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("'%v' may not be blank", field)
	}
	return nil
}

func Email(field string, value string) error {
	if !emailRegex.MatchString(value) {
		return fmt.Errorf("'%v' may not be invalid", field)
	}
	return nil
}

func NotNegative(field string, value int) error {
	if value < 0 {
		return fmt.Errorf("'%v' may not be negative", field)
	}
	return nil
}

func NotZero(field string, value int) error {
	if value != 0 {
		return fmt.Errorf("'%v' may only be zero", field)
	}
	return nil
}

func NotPositive(field string, value int) error {
	if value > 0 {
		return fmt.Errorf("'%v' may not be positive", field)
	}
	return nil
}

func NotZeroOrPositive(field string, value int) error {
	if value >= 0 {
		return fmt.Errorf("'%v' may not be zero or positive", field)
	}
	return nil
}
