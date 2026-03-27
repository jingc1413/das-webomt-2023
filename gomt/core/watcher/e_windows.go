package watcher

import (
	"log"
	"os/exec"
	"strconv"
	"syscall"
)

func isRunning(pid int) bool {
	h, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(h)

	var exitCode uint32
	err = syscall.GetExitCodeProcess(h, &exitCode)
	if err != nil {
		return false
	}
	return exitCode == 0x103
}

func killRunning(pid int) {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(pid))
	err := kill.Run()
	if err != nil {
		log.Printf("kill process err, %d, %v", pid, err)
	}
}
