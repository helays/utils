package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"os"
	"path"
)

// SetPath 设置当前 文件全路径
// sPath 如果是绝对路径，那么直接返回sPath
// 如果是相对路径，会跟上当前目录
func SetPath(sftpClient *sftp.Client, sPath string) (string, error) {
	//if path.IsAbs(sPath) {
	//	return sPath, nil
	//}
	// 判断当前目录是否是根目录
	current, err := sftpClient.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前目录失败：%s", err.Error())
	}
	return path.Join(current, sPath), nil
}

// Exist 判断路径是否存在
// 如果存在返回true，
// 当返回false的时候，需要判断error是否未nil，ni来的时候标识文件夹不存在
func Exist(sftpClient *sftp.Client, sPath string) (bool, error) {
	if _, err := sftpClient.Stat(sPath); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}

// Mkdir 创建文件夹
func Mkdir(sftpClient *sftp.Client, sPath string) error {
	if ok, err := Exist(sftpClient, sPath); ok {
		return nil
	} else if err != nil {
		return err
	}
	if err := sftpClient.MkdirAll(sPath); err != nil {
		return fmt.Errorf("创建文件夹%s失败：%s", sPath, err.Error())
	}
	return nil
}
