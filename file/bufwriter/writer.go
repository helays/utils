package bufwriter

import (
	"bufio"
	"fmt"
	"os"
)

// 默认间隔的指数
const (
	defaultFlushPower = 14 // 2^14 = 16384
	defaultSyncPower  = 17 // 2^17 = 131072
)

type Writer struct {
	file          *os.File
	buf           *bufio.Writer
	counter       int64
	flushPower    uint // 刷新间隔的指数
	syncPower     uint // 同步间隔的指数
	flushInterval int  // 实际的刷新间隔
	syncInterval  int  // 实际的同步间隔
	flushMask     int  // 刷新掩码
	syncMask      int  // 同步掩码
}

func New(path string) (*Writer, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	w := &Writer{
		file:       file,
		buf:        bufio.NewWriter(file),
		counter:    0,
		flushPower: defaultFlushPower,
		syncPower:  defaultSyncPower,
	}

	// 计算实际的间隔和掩码
	w.flushInterval = 1 << w.flushPower
	w.syncInterval = 1 << w.syncPower
	w.flushMask = w.flushInterval - 1
	w.syncMask = w.syncInterval - 1

	return w, nil
}

func NewWithWriter(writer *os.File) (*Writer, error) {
	w := &Writer{
		file:       writer,
		buf:        bufio.NewWriter(writer),
		counter:    0,
		flushPower: defaultFlushPower,
		syncPower:  defaultSyncPower,
	}

	// 计算实际的间隔和掩码
	w.flushInterval = 1 << w.flushPower
	w.syncInterval = 1 << w.syncPower
	w.flushMask = w.flushInterval - 1
	w.syncMask = w.syncInterval - 1

	return w, nil
}

// SetFlushInterval 设置刷新间隔的指数
func (w *Writer) SetFlushInterval(power uint) {
	w.flushPower = power
	w.flushInterval = 1 << power
	w.flushMask = w.flushInterval - 1
}

// SetSyncInterval 设置同步间隔的指数
func (w *Writer) SetSyncInterval(power uint) {
	w.syncPower = power
	w.syncInterval = 1 << power
	w.syncMask = w.syncInterval - 1
}

// 写入数据
func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.buf.Write(p)
	if err != nil {
		return n, fmt.Errorf("写入数据到bufio失败:%v", err)
	}
	w.counter++

	// 使用位运算替代取余运算
	if w.counter&int64(w.flushMask) == 0 {
		if err = w.buf.Flush(); err != nil {
			return n, fmt.Errorf("刷新bufio失败:%v", err)
		}
	}
	if w.counter&int64(w.syncMask) == 0 {
		if err = w.file.Sync(); err != nil {
			return n, fmt.Errorf("同步文件失败:%v", err)
		}
	}

	return n, nil
}

// Flush 手动刷新缓冲区
func (w *Writer) Flush() error {
	if err := w.buf.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲区失败: %v", err)
	}
	if err := w.file.Sync(); err != nil {
		return fmt.Errorf("同步文件失败: %v", err)
	}
	return nil
}

// Close 关闭写入器
func (w *Writer) Close() error {
	var errs []error

	// 1. 先刷新缓冲区（确保所有数据写入文件）
	if err := w.buf.Flush(); err != nil {
		errs = append(errs, fmt.Errorf("刷新缓冲区失败: %v", err))
		// 继续执行，尝试同步和关闭文件
	}

	// 2. 同步文件到磁盘
	if err := w.file.Sync(); err != nil {
		errs = append(errs, fmt.Errorf("同步文件失败: %v", err))
		// 继续执行，尝试关闭文件
	}

	// 3. 关闭文件
	if err := w.file.Close(); err != nil {
		errs = append(errs, fmt.Errorf("关闭文件失败: %v", err))
	}

	// 如果有多个错误，返回组合错误
	if len(errs) > 0 {
		return fmt.Errorf("关闭写入器时发生多个错误: %v", errs)
	}

	return nil
}

// GetIntervals 获取当前间隔设置（用于调试）
func (w *Writer) GetIntervals() (flushPower, syncPower uint, flushInterval, syncInterval int) {
	return w.flushPower, w.syncPower, w.flushInterval, w.syncInterval
}
