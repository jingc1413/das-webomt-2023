package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func loadLockFile(name string) int {
	name = filepath.Base(name)
	lockFile := fmt.Sprintf(".%s.pid", name)
	pidBytes, err := os.ReadFile(lockFile)
	if err != nil {
		log.Errorf("load lock file, %s, %v", name, err)
		return -1
	}
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return -1
	}
	return pid
}

func SaveLockFile(name string, pid int) {
	name = filepath.Base(name)
	ext := filepath.Ext(name)
	prefix := name[:len(name)-len(ext)]
	lockFile := fmt.Sprintf(".%s.pid", prefix)
	pidString := fmt.Sprintf("%d", pid)
	os.WriteFile(lockFile, []byte(pidString), 0666)
}

func CheckRunning(name string) {
	if !checkIsAutoRestart() {
		name = filepath.Base(name)
		ext := filepath.Ext(name)
		prefix := name[:len(name)-len(ext)]
		pid := loadLockFile(prefix)
		if pid > 0 && isRunning(pid) {
			killRunning(pid)
			log.Printf("kill %s stop, [PID] %d\n", name, pid)
		}
	}
}
