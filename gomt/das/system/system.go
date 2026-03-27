package system

import (
	"gomt/core/model"
	"gomt/das/agent"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewDasSystem() *DasSystem {
	return &DasSystem{
		agents: sync.Map{},
		log:    logrus.WithFields(logrus.Fields{"manager": "agent"}),
	}
}

type DasSystem struct {
	log    *logrus.Entry
	agents sync.Map
}

func (m *DasSystem) GetAllDeviceInfos() []*model.DeviceInfo {
	infos := []*model.DeviceInfo{}

	m.agents.Range(func(key any, value any) bool {
		if key == "local" {
			return true
		}
		if a := value.(*agent.DasDeviceAgent); a != nil {
			if info := a.GetDeviceInfo(); info.RouteAddressString != "" {
				infos = append(infos, info)
			}
		}
		return true
	})

	return infos
}

func (m *DasSystem) SetupDasDeviceAgent(
	schema string,
	deviceSub string,
	info *model.DeviceInfo,
	opts agent.DasDeviceAgentOptions,
) (*agent.DasDeviceAgent, error) {
	logrus.Tracef("setup device agent %v", deviceSub)
	a, err := agent.NewDasDeviceAgent(schema, deviceSub, info, opts)
	if err != nil {
		return nil, errors.Wrap(err, "create device agent")
	}
	m.agents.Store(deviceSub, a)
	go a.Run()
	return a, nil
}

func (m *DasSystem) DeleteAgent(deviceSub string) {
	logrus.Tracef("delete device agent %v", deviceSub)
	if agt := m.GetAgent(deviceSub); agt != nil {
		agt.Close()
	}
	m.agents.Delete(deviceSub)
}

func (m *DasSystem) ClearAllAgent() {
	m.agents.Range(func(key interface{}, value interface{}) bool {
		if agt, ok := value.(*agent.DasDeviceAgent); ok {
			agt.Close()
		}
		m.agents.Delete(key)
		return true
	})
}

func (m *DasSystem) GetAgent(deviceSub string) *agent.DasDeviceAgent {
	// logrus.Tracef("get device agent %v", deviceSub)
	v, ok := m.agents.Load(deviceSub)
	if ok {
		return v.(*agent.DasDeviceAgent)
	}
	return nil
}
