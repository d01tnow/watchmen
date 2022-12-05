package main

import (
	"fmt"
	"os"
	"syscall"
)

type Flock struct {
	file *os.File
}

func NewFlock(file string) (*Flock, error) {
	fi, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s. error: %v", file, err)
	}

	var fl Flock
	fl.file = fi
	return &fl, nil
}

func (fl *Flock) TryLock() error {
	// flock 是建议性锁, 不具备强制性. 也就是一个进程将一个文件锁住, 不能阻止另一个进程修改文件数据
	// LOCK_EX: 排他锁
	// LOCK_NB: non-block, 在无法锁定文件时不阻塞, 立刻返回
	err := syscall.Flock(int(fl.file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("cannot lock file: %s. error: %v", fl.file.Name(), err)
	}
	pid := os.Getpid()
	fl.file.Write([]byte(fmt.Sprintf("%d", pid)))
	fmt.Println("lock file path:", fl.file.Name())
	return nil
}

func (f *Flock) Close() error {
	if f.file == nil {
		return nil
	}
	err := syscall.Flock(int(f.file.Fd()), syscall.LOCK_UN)
	if err != nil {
		return fmt.Errorf("failed unlock file: %s. error: %v", f.file.Name(), err)
	}
	fp := f.file.Name()
	err = f.file.Close()
	if err != nil {
		fmt.Printf("failed to close file: %s. error: %v\n", fp, err)
	}
	err = os.Remove(fp)
	if err != nil {
		fmt.Printf("failed to remove file: %s. error: %v\n", fp, err)
	}
	fmt.Printf("file(%s) closed.\n", fp)
	return nil
}
