package tools

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"helay.net/go/utils/v3/close/osClose"
	"helay.net/go/utils/v3/close/vclose"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/logger/ulogs"
)

type ReadRowCallback func(scanner *bufio.Scanner) error

func ReadRowWithFile(file io.Reader, callback ReadRowCallback) error {
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

func FilePutWithReader(path string, rd io.Reader) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	defer vclose.Close(f)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, rd)
	return err
}

// FilePutContents 快速简易写文件
func FilePutContents(path, content string) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	osClose.CloseFile(file)
	return err
}

func FilePutContentsbytes(path string, content []byte) error {
	_path := filepath.Dir(path)
	if _, err := os.Stat(_path); err != nil {
		if err := Mkdir(_path); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	osClose.CloseFile(file)
	return err
}

// FileAppendContents 快速简易写文件（追加）
func FileAppendContents(path, content string) error {
	_path := filepath.Dir(path)
	if _, err := os.Stat(_path); err != nil {
		if err := Mkdir(_path); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	osClose.CloseFile(file)
	return err
}

// FileGetContents 快速简易读取文件
func FileGetContents(path string) ([]byte, error) {
	file, err := os.Open(path)
	defer vclose.Close(file)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(file)
}

// Mkdir 判断目录是否存在，否则创建目录
func Mkdir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// Fileabs 生成文件的绝对路径
// noinspection SpellCheckingInspection
func Fileabs(cpath string) string {
	if filepath.IsAbs(cpath) {
		return cpath
	}
	return filepath.Join(config.Appath, cpath)
}

// FileAbsWithCurrent 生成文件的绝对路径,根目录手动指定
func FileAbsWithCurrent(current, cpath string) string {
	if filepath.IsAbs(cpath) {
		return cpath
	}
	return filepath.Join(current, cpath)
}

func RemoveAll(path string) {
	ulogs.Checkerr(os.RemoveAll(path), "删除文件失败")
}
