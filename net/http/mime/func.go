package mime

import (
	"bytes"
	"fmt"
	"github.com/helays/utils/close/osClose"
	"io"
	"net/http"
	"os"
)

// GetFilePathMimeType 接受文件路径并返回文件的MIME类型。
// 如果发生错误，返回错误信息。
func GetFilePathMimeType(filePath string) (string, error) {
	// 打开指定路径的文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %w", err)
	}
	defer osClose.CloseFile(file)
	return GetFileMimeType(file)
}

// GetFileMimeType 接受一个已打开的文件指针并返回文件的MIME类型。
// 如果发生错误，返回错误信息。在
func GetFileMimeType(file io.Reader) (string, error) {
	// 读取文件的前512个字节
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("读取文件时出错: %w", err)
	}
	// 检测MIME类型
	contentType := http.DetectContentType(buffer[:n])
	return contentType, nil
}

// GetMimeType 接受一个已打开的文件指针并返回文件的MIME类型。
func GetMimeType(file io.Reader) (io.Reader, string, error) {
	// 读取文件的前512个字节
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, "", fmt.Errorf("读取文件时出错: %w", err)
	}
	// 检测MIME类型
	contentType := http.DetectContentType(buffer[:n])

	// 如果需要重新使用这个文件，可以将已读取的数据和原Reader组合起来
	combinedReader := io.MultiReader(bytes.NewReader(buffer[:n]), file)

	// 假设你需要返回这个 combinedReader 以供后续使用
	// 这里仅作为演示，实际使用中可能需要调整函数签名来返回 combinedReader
	// 此处简化处理，直接返回 mimeType 和 nil 错误
	return combinedReader, contentType, nil
}
