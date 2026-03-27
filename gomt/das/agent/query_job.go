package agent

import (
	"encoding/json"
	"fmt"
	"gomt/core/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type QueryDevicesJobOptions struct {
	Schema      string
	QueryString string
}

type QueryDevicesJob struct {
	sync.Mutex
	Token     string                 `json:"token"`
	Options   QueryDevicesJobOptions `json:"-"`
	IsRunning bool                   `json:"isRunning"`
	ws        *websocket.Conn        `json:"-"`
	log       *logrus.Entry          `json:"-"`
}

func CreateQueryDevicesJob(opts QueryDevicesJobOptions) *QueryDevicesJob {
	job := &QueryDevicesJob{
		Token:   fmt.Sprintf("%d", time.Now().Nanosecond()),
		Options: opts,
	}
	job.log = logrus.WithField("job", job.Token)
	return job
}

type QueryDevicesJobMessage struct {
	Total      int
	Index      int
	DeviceInfo *model.DeviceInfo
}

func (m *QueryDevicesJob) SetWebSocket(ws *websocket.Conn) {
	m.Lock()
	defer m.Unlock()
	m.ws = ws
}

func (m *QueryDevicesJob) HasWebSockect() bool {
	m.Lock()
	defer m.Unlock()
	return m.ws != nil
}

func (m *QueryDevicesJob) WriteMessage(total uint8, index uint8, info *model.DeviceInfo) {
	if m.ws != nil {
		msg := QueryDevicesJobMessage{
			Total: int(total),
			Index: int(index),
		}
		content, _ := json.Marshal(msg)
		m.ws.WriteMessage(websocket.TextMessage, content)
	}
}

func (m *QueryDevicesJob) onQueryDevice(total uint8, index uint8, info *model.DeviceInfo) {
	m.WriteMessage(total, index, info)
}

func (m *QueryDevicesJob) Run(deviceAgent *DasDeviceAgent) error {
	m.Lock()
	defer m.Unlock()

	if m.ws == nil {
		return errors.New("ws not connected")
	}
	if m.IsRunning {
		return errors.New("alreay IsRunning")
	}
	m.IsRunning = true
	go func() {
		m.log.Tracef("running")
		defer m.log.Tracef("stopped")

		defer func() {
			m.IsRunning = false
			if m.ws != nil {
				m.ws.Close()
			}
		}()

		_, err := deviceAgent.QueryDevices(m.Options.Schema, m.Options.QueryString, m.onQueryDevice)
		if err != nil {
			m.log.Error(errors.Wrap(err, "query devices"))
			return
		}
		m.log.Tracef("finish")
	}()
	return nil
}

func (m *QueryDevicesJob) Cancel() error {
	m.Lock()
	defer m.Unlock()
	if m.ws == nil {
		return errors.New("ws not connected")
	}
	if !m.IsRunning {
		return errors.New("not IsRunning")
	}
	m.IsRunning = false

	go m.SetWebSocket(nil)
	return nil
}
