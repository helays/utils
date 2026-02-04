package fileinclude

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"helay.net/go/utils/v3/close/vclose"
)

// Processor 文件包含处理器
type Processor struct {
	baseDir      string
	prefix       string
	visited      map[string]bool
	maxDepth     int
	currentDepth int
	input        inputSource
}

type inputSource struct {
	file      string    // 用于文件输入
	reader    io.Reader // 用于reader输入
	str       string    // 用于字符串输入
	inputType inputType
}

type inputType int

const (
	inputTypeNone inputType = iota
	inputTypeFile
	inputTypeReader
	inputTypeString
)

// NewProcessor 创建新的处理器
func NewProcessor() *Processor {
	return &Processor{
		visited:  make(map[string]bool),
		maxDepth: 100,
		prefix:   "include",
	}
}

func (p *Processor) SetPrefix(prefix string) *Processor {
	p.prefix = prefix
	return p
}

// SetMaxDepth 设置最大递归深度
func (p *Processor) SetMaxDepth(depth int) *Processor {
	p.maxDepth = depth
	return p
}

// FromFile 从文件路径输入
func (p *Processor) FromFile(filename string) *Processor {
	p.baseDir = filepath.Dir(filename)
	p.input = inputSource{
		file:      filename,
		inputType: inputTypeFile,
	}
	return p
}

// FromReader 从io.Reader输入
func (p *Processor) FromReader(r io.Reader, currentDir string) *Processor {
	p.baseDir = currentDir
	p.input = inputSource{
		reader:    r,
		inputType: inputTypeReader,
	}
	return p
}

// FromString 从字符串输入
func (p *Processor) FromString(content string, currentDir string) *Processor {
	p.baseDir = currentDir
	p.input = inputSource{
		str:       content,
		inputType: inputTypeString,
	}
	return p
}

// ToString 输出为字符串
func (p *Processor) ToString() (string, error) {
	var input io.Reader
	switch p.input.inputType {
	case inputTypeFile:
		file, err := os.Open(p.input.file)
		if err != nil {
			return "", fmt.Errorf("failed to open file: %v", err)
		}
		defer vclose.Close(file)
		input = file
	case inputTypeReader:
		input = p.input.reader
	case inputTypeString:
		input = strings.NewReader(p.input.str)
	default:
		return "", fmt.Errorf("no input source specified")
	}

	return p.processInput(input)
}

// ToReader 输出为io.Reader
func (p *Processor) ToReader() (io.Reader, error) {
	content, err := p.ToString()
	if err != nil {
		return nil, err
	}
	return strings.NewReader(content), nil
}

// ToFile 输出到文件
func (p *Processor) ToFile(outputPath string) error {
	content, err := p.ToString()
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(content), 0644)
}

// 处理输入的核心逻辑
func (p *Processor) processInput(input io.Reader) (string, error) {
	if p.currentDepth >= p.maxDepth {
		return "", fmt.Errorf("max include depth %d exceeded", p.maxDepth)
	}

	var result strings.Builder
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, p.prefix) {
			// 不是包含语句，直接写入结果
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}
		pathPattern := strings.TrimSpace(line[len(p.prefix):])
		matchedFiles, err := filepath.Glob(filepath.Join(p.baseDir, pathPattern))
		if err != nil {
			return "", fmt.Errorf("invalid path pattern %q: %v", pathPattern, err)
		}
		for _, matchedFile := range matchedFiles {
			absPath, _err := filepath.Abs(matchedFile)
			if _err != nil {
				return "", fmt.Errorf("failed to get absolute path for %q: %v", matchedFile, _err)
			}
			if p.visited[absPath] {
				return "", fmt.Errorf("file %q has already been included", matchedFile)
			}

			p.visited[absPath] = true

			p.currentDepth++
			content, __err := p.processFile(matchedFile)
			p.currentDepth--

			if __err != nil {
				return "", fmt.Errorf("error processing included file %q: %v", matchedFile, __err)
			}
			result.WriteString(content)
		}

	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}

	return result.String(), nil
}

func (p *Processor) processFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file stat %s: %v", filename, err)
	}
	if stat.IsDir() {
		return "", nil
	}
	defer vclose.Close(file)

	// 创建子处理器处理包含的文件
	subProcessor := &Processor{
		baseDir:      filepath.Dir(filename),
		visited:      p.visited, // 共享visited map
		maxDepth:     p.maxDepth,
		currentDepth: p.currentDepth,
		prefix:       p.prefix,
	}

	return subProcessor.processInput(file)
}
