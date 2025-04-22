package tools

import (
	"bufio"
	"os"
)

type ReadRowCallback func(scanner *bufio.Scanner) error

func ReadRowWithFile(file *os.File, callback ReadRowCallback) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if err := callback(scanner); err != nil {
			return err
		}
	}
	return nil
}
