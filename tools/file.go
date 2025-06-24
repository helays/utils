package tools

import (
	"bufio"
	"os"
	"strings"
)

type ReadRowCallback func(scanner *bufio.Scanner) error

func ReadRowWithFile(file *os.File, callback ReadRowCallback) error {
	scanner := bufio.NewScanner(file)
	// 初始 4KB，最大 10MB（可根据需求调整）
	scanner.Buffer(make([]byte, 4096), 10*1024*1024)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if err := callback(scanner); err != nil {
			return err
		}
	}
	return nil
}

func ContainsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }
