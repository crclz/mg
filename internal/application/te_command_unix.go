//go:build !windows

package application

import (
	"fmt"
	"os"
	"syscall"
)

func (p *TeCommand) SysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
}

func (p *TeCommand) TryCascadeKill(process *os.Process) {
	// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	// Solution: In addition to sending a signal to a single PID,
	// kill(2) also supports sending a signal to a Process Group by passing
	// the process group id (PGID) as a negative number.

	var pgid = -1 * process.Pid

	var err = syscall.Kill(pgid, syscall.SIGKILL)
	if err != nil {
		fmt.Printf("Kill pgid %v error: %+v\n", pgid, err)
	} else {
		fmt.Printf("Kill pgid %v success\n", pgid)
	}
}
