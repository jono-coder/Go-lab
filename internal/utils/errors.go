package utils

import "fmt"

var ErrNotFound = fmt.Errorf("not found")

func NPE(entity string) error {
	return fmt.Errorf("'%v' may not be nil", entity)
}
