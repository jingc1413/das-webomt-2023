package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PingOptions struct {
	Address  []string
	Timeout  int
	Count    int
	Interval int
}

type TaskJobPing struct {
	sync.Mutex
	Token     string
	Options   PingOptions
	IsRunning bool
	ws        *websocket.Conn `json:"-"`
	log       *logrus.Entry   `json:"-"`
}

func (m *TaskJobPing) WriteMessage(typ string, data any) {
	if m.ws != nil {
		msg := TestMessage{
			Type: typ,
			Data: data,
			Time: time.Now().Unix(),
		}
		content := msg.String()
		logrus.Trace(content)
		m.ws.WriteMessage(websocket.TextMessage, []byte(content))
	}
}

func (m *TaskJobPing) SetWebSocket(ws *websocket.Conn) {
	m.Lock()
	defer m.Unlock()
	m.ws = ws
}

func (m *TaskJobPing) Cancel() error {
	m.Lock()
	defer m.Unlock()
	m.ws = nil
	return nil
}

func (m *TaskJobPing) Run() error {
	m.Lock()
	defer m.Unlock()
	if m.ws == nil {
		return errors.New("ws not connected")
	}
	opts := m.Options

	if m.IsRunning {
		return errors.New("already running")
	}
	m.IsRunning = true
	go func() {
		m.log.Infof("running")
		defer m.log.Infof("stopped")

		defer func() {
			m.IsRunning = false
			if m.ws != nil {
				m.ws.Close()
			}
		}()

		for _, v := range opts.Address {
			_, err := m.doPing(v, opts)
			if err != nil {
				m.log.Errorf("exec ping err, %v", err)
			}
		}
	}()
	return nil
}

func (m *TaskJobPing) doPing(addr string, opts PingOptions) (pingStats, error) {
	stats := pingStats{
		PacketLoss: 100.0,
	}
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return stats, errors.Wrap(err, "create pinger")
	}
	pinger.Count = opts.Count

	if opts.Timeout > 0 {
		pinger.Timeout = time.Duration(opts.Timeout) * time.Second
	}
	if opts.Interval > 0 {
		pinger.Interval = time.Duration(opts.Interval) * time.Second
	}
	pinger.SetPrivileged(true)

	pinger.OnSetup = func() {
		m.WriteMessage("PingOnSetup", nil)
	}
	pinger.OnSend = func(pkt *ping.Packet) {
		m.WriteMessage("PingOnSend", *pkt)
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
		m.WriteMessage("PingOnRecv", *pkt)
	}
	pinger.OnFinish = func(out *ping.Statistics) {
		stats.TTL = pinger.TTL
		stats.Size = pinger.Size
		stats.Interval = int(pinger.Interval.Seconds())
		stats.Addr = pinger.Addr()
		stats.IPAddr = pinger.IPAddr().String()
		if out != nil {
			stats.PacketsSent = out.PacketsSent
			stats.PacketsRecv = out.PacketsRecv
			stats.PacketsRecvDuplicates = out.PacketsRecvDuplicates
			if stats.PacketsSent != 0 {
				stats.PacketLoss = out.PacketLoss
			}
			stats.MinRtt = out.MinRtt
			stats.MaxRtt = out.MaxRtt
			stats.AvgRtt = out.AvgRtt
		}
		m.WriteMessage("PingStats", stats)
	}
	if err := pinger.Run(); err != nil {
		return stats, errors.Wrap(err, "run pinger")
	}
	return stats, nil
}

type pingStats struct {
	Addr                  string
	IPAddr                string
	PacketsSent           int
	PacketsRecv           int
	PacketsRecvDuplicates int
	PacketLoss            float64
	MaxRtt                time.Duration
	MinRtt                time.Duration
	AvgRtt                time.Duration
	Size                  int
	TTL                   int
	Interval              int
}

func isValidIP(str string) bool {
	ip := net.ParseIP(str)
	if ip != nil && ip.To4() != nil {
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
	}
	_, err := resolver.LookupHost(ctx, str)
	return err == nil
}

func CreateTaskJob(opts PingOptions) (*TaskJobPing, error) {
	{
		if opts.Address == nil || len(opts.Address) < 1 {
			return nil, errors.New("Invalid ip address")
		}
		for _, v := range opts.Address {
			if !isValidIP(v) {
				return nil, errors.New("Invalid ip address")
			}
		}
		if opts.Count < 1 {
			opts.Count = 1 // 1-100
		}
		if opts.Count > 100 {
			opts.Count = 100 // 1-100
		}
		if opts.Timeout < 0 {
			opts.Timeout = 0 // 0-60
		}
		if opts.Timeout > 60 {
			opts.Timeout = 60 // 0-60
		}
		if opts.Interval < 1 {
			opts.Interval = 1 // 1-15
		}
		if opts.Interval > 15 {
			opts.Interval = 15 // 1-15
		}
	}
	job := &TaskJobPing{
		Token:   fmt.Sprintf("%d", time.Now().Nanosecond()),
		Options: opts,
	}
	job.log = logrus.WithField("ping-job", job.Token)
	return job, nil
}

type TaskManager struct {
	sync.Mutex
	jobs map[string]*TaskJobPing
}

func (m *TaskManager) SetTaskJob(job *TaskJobPing) {
	m.Lock()
	defer m.Unlock()
	m.jobs[job.Token] = job
}

func (m *TaskManager) GetTaskJob(token string) *TaskJobPing {
	m.Lock()
	defer m.Unlock()
	job, ok := m.jobs[token]
	if !ok {
		return nil
	}
	return job
}

var defaultManager = TaskManager{
	jobs: make(map[string]*TaskJobPing),
}

func SetTaskJob(job *TaskJobPing) {
	defaultManager.SetTaskJob(job)
}

func GetTaskJob(token string) *TaskJobPing {
	return defaultManager.GetTaskJob(token)
}

func RemoveTaskJob(token string) {
	delete(defaultManager.jobs, token)
}

type TestMessage struct {
	Type string
	Time int64
	Data any `json:"Data,omitempty"`
}

func (m TestMessage) String() string {
	v, _ := json.Marshal(m)
	return string(v)
}
