//go:build windows

package application

import (
	"fmt"
	"os"
	"syscall"
)

func (p *TeCommand) SysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func (p *TeCommand) TryCascadeKill(process *os.Process) {
	var err = process.Kill()
	if err != nil {
		fmt.Printf("TryCascadeKill %v error: %v\n", process.Pid, err)
	} else {
		fmt.Printf("TryCascadeKill %v success\n", process.Pid)
	}
}
