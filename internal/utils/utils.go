package utils

import (
	"fmt"
	"os"
)

func EnsureDir(dirName string) error {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking directory: %v", err)
	}

	return nil
}
