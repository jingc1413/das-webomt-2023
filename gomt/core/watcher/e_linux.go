package watcher

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"
)

//func isRunning(pid int) bool {
//	if killErr := syscall.Kill(pid, syscall.Signal(0)); killErr == nil {
//		return true
//	}
//	return false
//}

func isRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}
	if err.Error() == "os: process already finished" {
		return false
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	switch errno {
	case syscall.ESRCH:
		return false
	case syscall.EPERM:
		return true
	}
	return false
}

func killRunning(pid int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	p, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("find process err, %d, %v\n", pid, err)
		return
	}
	p.Signal(syscall.SIGINT)
	done := make(chan error)
	go func() {
		_, e := p.Wait()
		done <- e
	}()

	select {
	case <-ctx.Done():
		log.Printf("process is still running...%d", pid)
		p.Signal(syscall.SIGKILL)
	case <-done:
		return
	}
	////////
	/*timeout := time.After(3 * time.Second)
		ticker := time.NewTicker(time.Second)

		p, e := os.FindProcess(pid)
		if e != nil {
			log.Printf("find process err, %d, %v\n", pid, e)
			return
		}
		p.Signal(syscall.SIGINT)
		p.Wait()
	loop:
		for {
			select {
			case <-timeout:
				break loop
			case <-ticker.C:
				if !isRunning(pid) {
					log.Printf("process is not running...%d", pid)
					return
				}
				log.Printf("process is staill running...%d", pid)
			}
		}
		log.Printf("process is staill running, sending SIGKILL...%d\n", pid)
		p.Signal(syscall.SIGKILL)
		return*/
}
