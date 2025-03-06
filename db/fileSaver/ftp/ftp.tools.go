package ftp

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"path"
	"strings"
)

// SetPath 设置当前 文件全路径
func SetPath(ftpClient *ftp.ServerConn, sPath string) (string, error) {
	//if path.IsAbs(sPath) {
	//	return sPath, nil
	//}
	// 判断当前目录是否是根目录
	current, err := ftpClient.CurrentDir()
	if err != nil {
		return "", fmt.Errorf("获取当前目录失败：%s", err.Error())
	}
	return path.Join(current, sPath), nil
}

// Exist 判断远程文件是否存在
func Exist(ftpClient *ftp.ServerConn, sPath string) (bool, error) {
	remotePath := path.Dir(sPath)
	lst, err := ftpClient.List(remotePath)
	if err != nil {
		return false, fmt.Errorf("判断文件%s是否存在失败：%s", sPath, err.Error())
	}
	remoteFileName := path.Base(sPath)
	for _, v := range lst {
		if v.Name == remoteFileName {
			return true, nil
		}
	}
	return false, nil
}

func Mkdir(ftpClient *ftp.ServerConn, sPath string) error {
	// 判断文件夹是否存在
	// 如果报错，就直接返回
	// 如果文件夹存在就不创建了

	if ok, err := Exist(ftpClient, sPath); err != nil {
		return err
	} else if ok {
		return nil
	}
	return mkdirAll(ftpClient, sPath)
}

func mkdirAll(ftpClient *ftp.ServerConn, sPath string) error {
	// 这里要处理下目录情况，先获取当前所在目录
	currentDir, err := ftpClient.CurrentDir()
	if err != nil {
		return fmt.Errorf("创建目录%s，获取当前目录失败：%s", sPath, err.Error())
	}
	// 将path中，当前路径部分去掉
	sPath = strings.TrimPrefix(sPath, currentDir)
	// Split the path into its components.
	var currentPath string
	for _, part := range strings.Split(sPath, "/") {
		if part == "" {
			continue // Skip empty parts which can happen with leading/trailing slashes or double slashes.
		}
		currentPath = fmt.Sprintf("%s/%s", currentPath, part)
		// Check if the current directory exists.
		err = ftpClient.ChangeDir(currentPath)
		if err != nil {
			// Directory does not exist, so create it.
			if err = ftpClient.MakeDir(part); err != nil {
				return err
			}
			// Change to the newly created directory.
			if err = ftpClient.ChangeDir(part); err != nil {
				return err
			}
		} else {
			// Directory already exists, continue to next part.
			continue
		}
	}
	return ftpClient.ChangeDir(currentDir)
}
