package timing

import (
	"sync"
	"time"

	"github.com/RussellLuo/timingwheel"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type TimingServer interface {
	Run()
	Stop()

	StartTimer(key string, d time.Duration, f func())
	StopTimer(key string)
	HasTimer(key string) bool
}

type baseTimingServer struct {
	name string
	log  *logrus.Entry

	tw *timingwheel.TimingWheel

	timers     map[string]*timingwheel.Timer
	timersLock sync.RWMutex

	quit     chan bool
	shutdown bool
}

func NewBaseTimingServer(tick time.Duration, size int64) (TimingServer, error) {
	s := &baseTimingServer{
		name: "timing",
		quit: make(chan bool),
	}
	s.log = logrus.WithFields(logrus.Fields{"server": "timing"})
	s.timers = map[string]*timingwheel.Timer{}
	s.tw = timingwheel.NewTimingWheel(tick, size)
	s.log.Info("new timing server")
	return s, nil
}

func (s *baseTimingServer) Run() {
	s.log.Info("start timing server")

	s.tw.Start()

	<-s.quit

	s.tw.Stop()
	s.log.Info("stop timing server")
}

func (s *baseTimingServer) Stop() {
	s.shutdown = true
	s.quit <- true
}

func (s *baseTimingServer) StartTimer(key string, d time.Duration, f func()) {
	s.stopTimer(key)
	if f == nil {
		s.log.Error(errors.New("start timer error, func in null, " + key))
		return
	}

	//s.log.WithFields(log.Fields{"key": key, "expire": d}).Trace("start timer")
	t := s.tw.AfterFunc(d, func() {
		f()
	})
	s.timersLock.Lock()
	s.timers[key] = t
	s.timersLock.Unlock()
}

func (s *baseTimingServer) StopTimer(key string) {
	s.stopTimer(key)
}

func (s *baseTimingServer) HasTimer(key string) bool {
	s.timersLock.RLock()
	_, ok := s.timers[key]
	s.timersLock.RUnlock()
	return ok
}

func (s *baseTimingServer) stopTimer(key string) {
	s.timersLock.Lock()
	t, ok := s.timers[key]
	if ok && t != nil {
		t.Stop()
		delete(s.timers, key)
	}
	s.timersLock.Unlock()
}

var tsrv TimingServer = nil

func Run() {
	if tsrv != nil {
		return
	}
	s, err := NewBaseTimingServer(time.Millisecond, 20)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "init timing server error"))
		return
	}
	tsrv = s
	tsrv.Run()
}

func Stop() {
	if tsrv == nil {
		return
	}
	tsrv.Stop()
	tsrv = nil
}

func StartTimer(key string, d time.Duration, f func()) {
	if tsrv == nil {
		return
	}
	// logrus.Tracef("timing: start timer, %v, %v", key, int64(d/time.Second))
	tsrv.StartTimer(key, d, f)
}

func StopTimer(key string) {
	if tsrv == nil {
		return
	}
	// logrus.Tracef("timing: stop timer, %v", key)
	tsrv.StopTimer(key)
}

func HasTimer(key string) bool {
	if tsrv == nil {
		return false
	}
	return tsrv.HasTimer(key)
}
