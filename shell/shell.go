package shell

import (
	"bytes"
	"errors"
	"helay.net/go/utils/v3"
	"helay.net/go/utils/v3/crypto/md5"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/tools"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type CmdResponse struct {
	Stdout, StdErr []byte
}

// CmdRun 执行 shell语句,并返回 response 结构
func CmdRun(name string, args ...string) (CmdResponse, error) {
	cmd := exec.Command(name, args...)
	var stdOut, stdErr = bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	ulogs.Log("[ExecReturnBytes]", "shell 执行命令", cmd.String())
	err := cmd.Run()
	return CmdResponse{
		Stdout: stdOut.Bytes(),
		StdErr: stdErr.Bytes(),
	}, err
}

// ExecShell 执行 shell语句
func ExecShell(name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = bytes.NewBuffer(nil)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	ulogs.Log("shell 执行命令", cmd.String())

	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}

func ExecShellQuit(name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = bytes.NewBuffer(nil)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}

// ExecCtlShell 可手动结束的 shell命令
func ExecCtlShell(stop chan byte, name string, s ...string) (string, error) {
	cmd := exec.Command(name, s...)
	var (
		out = new(bytes.Buffer)
		w   = new(bytes.Buffer)
		end = make(chan byte)
	)
	cmd.Stdout = out
	cmd.Stderr = w
	ulogs.Log("shell 执行命令", cmd.String())
	go func(cmd *exec.Cmd, stop, end chan byte) {
		for {
			select {
			case <-stop:
				if cmd == nil {
					ulogs.Error("分析EXEC 不存在", cmd.String())
					continue
				}
				if cmd.Process == nil {
					ulogs.Error("分析 Process不存在", cmd.String())
					continue
				}
				if err := cmd.Process.Kill(); err != nil {
					ulogs.Error("手动结束进程失败", err)
				} else {
					ulogs.Log("手动结束命令", cmd.String())
				}

			case <-end:
				ulogs.Log("shell 执行完成", cmd.String())
				return
			}
		}
	}(cmd, stop, end)
	defer func() {
		go func() {
			end <- 1
		}()

	}()
	if err := cmd.Run(); err != nil {
		return out.String(), errors.New(err.Error() + "，" + w.String())
	}
	return out.String(), nil
}

// MachineCode 利用硬件信息 生成token
func MachineCode() string {
	_os := strings.ToLower(runtime.GOOS)
	_arch := strings.ToLower(runtime.GOARCH)
	ulogs.Log("系统版本", _os, "平台", _arch, "读取设备信息...")
	var (
		err         error
		info        string
		_macode     string
		machineCode string
	)
	switch _os {
	case "linux":
		info, err = ExecShell("dmidecode")
		if err != nil {
			_macode, err = ExecShell("/bin/bash", "-c", `"/sbin/ip link" | grep link | /usr/bin/sort | /usr/bin/uniq | /usr/bin/sha256sum`)
			if err != nil {
				ulogs.Error("机器码生成失败", err)
				os.Exit(1)
			}
		}
	case "darwin":
		var tmp []byte
		tmp, err = tools.FileGetContents("/Users/helay/go/src/company/vis-device/startUp/vis.agent/run/info")
		if err != nil {
			ulogs.Error("机器码生成失败", err)
			os.Exit(1)
		}
		info = string(tmp)
	default:
		os.Exit(1)
	}
	if _macode != "" {
		machineCode = md5.Md5string(_macode)
	} else {
		cpuPreg := regexp.MustCompile(`Processor[\s\S]+?ID.+?((?:[A-Z0-9]{2} ?){8})`)

		tmp := cpuPreg.FindStringSubmatch(info)
		if len(tmp) != 2 {
			ulogs.Error("系统信息获取失败")
			os.Exit(0)
		}
		cpuid := strings.TrimSpace(tmp[1])

		boardPreg := regexp.MustCompile(`Base Board[\s\S]+?Serial Number.+?([A-Z0-9]+)`)

		tmp = boardPreg.FindStringSubmatch(info)
		if len(tmp) != 2 {
			ulogs.Error("系统信息获取失败")
			os.Exit(0)
		}
		boardid := tmp[1]
		// 生成机器码
		machineCode = md5.Md5string(cpuid + utils.Salt + boardid)
	}
	ulogs.Debug("机器码", machineCode)
	return machineCode
}
