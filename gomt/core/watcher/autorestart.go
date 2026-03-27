package watcher

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultPeriod      = time.Second * 1
	defaultCleanPeriod = time.Second * 15
	autoRestartFlag    = "auto-restart"
)

var (
	WatchFilename string
	WatchPeriod   = defaultPeriod
	RestartFunc   = restartByExec
	CleanHandler  = func(any) {}
)

var (
	executableArgs []string
	executableEnvs []string
	executablePath string
	ticker         *time.Ticker
	startFileInfo  os.FileInfo

	ctx     context.Context
	ctxFunc context.CancelFunc
)

func init() {
	executableArgs = os.Args
	executableEnvs = os.Environ()
	executablePath, _ = filepath.Abs(os.Args[0])
}

func StartWatcher(f func(any), i any) {
	WatchFilename = executablePath
	ticker = time.NewTicker(WatchPeriod)
	ctx, ctxFunc = context.WithCancel(context.Background())
	go watcher(f, i)
}

func StopWatcher() {
	ctxFunc()
}

func StopProcess() {
	ticker = time.NewTicker(WatchPeriod)
	stopCtx, cancel := context.WithTimeout(context.Background(), defaultCleanPeriod)
	defer cancel()
	stopTimeout := defaultCleanPeriod

	for {
		select {
		case <-stopCtx.Done():
			log.Warn("process timeout interrupt")
			os.Exit(0)
		case <-ticker.C:
			stopTimeout = stopTimeout - time.Second
			log.Debugf("wait all routine exit...%v", stopTimeout)
			if runtime.NumGoroutine() < 3 {
				log.Debug("all routine exit success")
				return
			}
		}
	}
}

func watcher(c func(any), i any) {
	CleanHandler = c
	for {
		select {
		case <-ctx.Done():
			log.Debugf("exit watcher, %s", WatchFilename)
			return
		case <-ticker.C:
			if isChanged() {
				{
					_, err := exec.LookPath(executablePath)
					if err != nil {
						log.Errorf("exec find: %s", err)
						continue
					}
				}
				if CleanHandler != nil {
					CleanHandler(i)
				}
				RestartFunc()
			}
		}
	}
}

func checkIsAutoRestart() bool {
	for _, s := range executableArgs {
		if s == autoRestartFlag {
			return true
		}
	}
	return false
}

func remove() {
	executableArgs = slices.DeleteFunc(executableArgs, func(e string) bool {
		return e == autoRestartFlag
	})
}

func isChanged() bool {
	return isChangedByStat()
}

func isChangedByStat() bool {
	fileinfo, err := os.Stat(WatchFilename)
	if err == nil {
		if startFileInfo == nil {
			startFileInfo = fileinfo
			return false
		}
		if startFileInfo.ModTime() != fileinfo.ModTime() ||
			startFileInfo.Size() != fileinfo.Size() {
			return true
		}
		return false
	}
	log.Debugf("cannot find %s: %s", WatchFilename, err)
	return false
}

func restartByExec() {
	binary, err := exec.LookPath(executablePath)
	if err != nil {
		log.Errorf("lookpath error: %s", err)
		return
	}
	time.Sleep(1 * time.Second)
	remove()
	executableArgs = append(executableArgs, autoRestartFlag)
	//defer remove()
	execErr := syscall.Exec(binary, executableArgs, executableEnvs)
	if execErr != nil {
		log.Errorf("auto restart error: %s %v", binary, execErr)
	}
}
