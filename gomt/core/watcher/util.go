package watcher

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"sync"
)

type checksumMap struct {
	l sync.Mutex
	m map[string]string
}

func (a *checksumMap) updateFileChecksum(filename, newChecksum string) (ok bool) {
	a.l.Lock()
	defer a.l.Unlock()
	oldChecksum, ok := a.m[filename]
	if !ok || oldChecksum != newChecksum {
		a.m[filename] = newChecksum
		return true
	}
	return false
}

func fileChecksum(filename string) (checksum string, err error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if len(contents) == 0 {
		return "", errors.New("empty file, forcing rebuild without updating checksum")
	}

	h := sha256.New()
	if _, err := h.Write(contents); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
