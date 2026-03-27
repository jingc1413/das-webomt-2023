package priv

import (
	"gomt/core/model"
	"gomt/core/net/udp"

	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var defaultPrivMgmtServer *PrivMgmtServer
var defaultPrivMgmtServerMutex sync.Mutex
var defaultPrivMgmtServerOnce sync.Once

func GetDefaultPrivMgmtServer() *PrivMgmtServer {
	defaultPrivMgmtServerMutex.Lock()
	defer defaultPrivMgmtServerMutex.Unlock()
	return defaultPrivMgmtServer
}

func SetupDefaultPrivMgmtServer(opts PrivMgmtServerOptions) error {
	defaultPrivMgmtServerOnce.Do(func() {
		defaultPrivMgmtServerMutex.Lock()
		defer defaultPrivMgmtServerMutex.Unlock()
		s, err := NewPrivMgmtServer(opts)
		if err != nil {
			logrus.Fatal(err)
		}
		defaultPrivMgmtServer = s
	})
	return nil
}

type PrivMgmtServerOptions struct {
	Schema           string
	Addr             string
	MaxWorkers       int
	DeviceTypePrefix string
	// ApType        ApType
	// VpType        VpType
	// McpType       McpType
}

type PrivMgmtServer struct {
	opts PrivMgmtServerOptions
	log  *logrus.Entry
	srv  *udp.Server

	nextSerialLock     sync.Mutex
	nextSerial         uint16
	waitingResponseMap *sync.Map

	wg       sync.WaitGroup
	quit     chan bool
	shutdown bool
}

func NewPrivMgmtServer(opts PrivMgmtServerOptions) (*PrivMgmtServer, error) {
	s := &PrivMgmtServer{
		opts: opts,
		quit: make(chan bool),
	}
	s.log = logrus.WithFields(logrus.Fields{"server": "oam"})

	srv, err := udp.NewServer("omc", s.opts.Addr, s, 1000, s.opts.MaxWorkers)
	if err != nil {
		return nil, errors.Wrap(err, "create udp server")
	}
	s.srv = srv
	s.waitingResponseMap = &sync.Map{}
	return s, nil
}

func (s *PrivMgmtServer) Run() {
	s.log.Trace("running")
	defer s.log.Trace("stopped")
	defer s.wg.Wait()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.srv.Run()
	}()

	<-s.quit

	s.srv.Stop()
}

func (s *PrivMgmtServer) Stop() {
	s.log.Trace("stop")
	s.shutdown = true
	s.quit <- true
}

func (s *PrivMgmtServer) getSerial() uint16 {
	s.nextSerialLock.Lock()
	defer s.nextSerialLock.Unlock()

	s.nextSerial += 1
	if s.nextSerial > VPACKET_SERIAL_OMC_END {
		s.nextSerial = VPACKET_SERIAL_OMC_BEGIN
	}
	return s.nextSerial
}

func (s *PrivMgmtServer) sendPacket(addr string, pkt ApPacket) error {
	s.log.Tracef("send packet: %v", pkt.Dump())
	data, err := pkt.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "marshal ap packet")
	}

	if err := s.srv.SendMessage(addr, data); err != nil {
		return errors.Wrap(err, "send ap packet")
	}
	return nil
}

func (s *PrivMgmtServer) sendResponse(addr string, req ApPacket, vpFlag VpFlag, cmdStatus CmdStatus, mcpData McpData) error {
	resp := ApPacket{}
	resp.ApType = req.ApType
	resp.VpType = req.VpType
	resp.VpPacket.MainID = req.VpPacket.MainID
	resp.VpPacket.SubID = req.VpPacket.SubID
	resp.VpPacket.Serial = req.VpPacket.Serial
	resp.VpPacket.VpFlag = vpFlag
	resp.VpPacket.CPID = req.VpPacket.CPID
	resp.VpPacket.McpPacket.CmdID = req.VpPacket.McpPacket.CmdID
	resp.VpPacket.McpPacket.CmdStatus = cmdStatus
	resp.VpPacket.McpPacket.Data = mcpData
	if err := s.sendPacket(addr, resp); err != nil {
		return errors.Wrap(err, "send packet")
	}
	return nil
}

func (s *PrivMgmtServer) sendResponseWithoutData(addr string, req ApPacket, vpFlag VpFlag, cmdStatus CmdStatus) error {
	return s.sendResponse(addr, req, vpFlag, cmdStatus, EMPTY_MCP_DATA)
}

func (s *PrivMgmtServer) sendRequset(sess *DeviceSession, cmdID CmdID, mcpData McpData, wait int) (*ApPacket, error) {
	req := ApPacket{}
	req.ApType = sess.ApType
	req.VpType = sess.VpType
	req.VpPacket.MainID = sess.MainID
	req.VpPacket.SubID = sess.SubID
	req.VpPacket.Serial = s.getSerial()
	req.VpPacket.VpFlag = VP_FLAG_REQUSET
	req.VpPacket.CPID = sess.McpType
	req.VpPacket.McpPacket.CmdID = cmdID
	req.VpPacket.McpPacket.CmdStatus = CMD_STATUS_COMMAND
	req.VpPacket.McpPacket.Data = mcpData

	if err := s.sendPacket(sess.Addr, req); err != nil {
		return nil, errors.Wrap(err, "send packet")
	}

	if wait <= 0 {
		return nil, nil
	}

	resp := s.waitingResponse(req.VpPacket.Serial, wait)
	if resp == nil {
		return nil, errors.New("waiting response timeout")
	}
	s.updateDeviceSession(sess, resp)
	return resp, nil
}

func (s *PrivMgmtServer) waitingResponse(serial uint16, wait int) *ApPacket {
	key := fmt.Sprintf("%v", serial)
	_, ok := s.waitingResponseMap.Load(key)
	if ok {
		return nil
	}

	ch := make(chan *ApPacket)
	s.waitingResponseMap.Store(key, ch)
	defer s.waitingResponseMap.Delete(key)

	timeout := make(chan bool, 1)
	go func() {
		for i := 0; i < wait; i++ {
			time.Sleep(time.Second)

			if s.shutdown {
				break
			}
		}
		timeout <- true
	}()

	select {
	case resp := <-ch:
		return resp
	case <-timeout:
		return nil
	}
}

func (s *PrivMgmtServer) updateDeviceSession(sess *DeviceSession, pkt *ApPacket) {
	if sess == nil {
		return
	}
	sess.SetDeviceID(pkt.VpPacket.MainID, pkt.VpPacket.SubID)
	if sess.ApType != pkt.ApType {
		s.log.Error(fmt.Sprintf("incorrect access protocol type, 0x%02x!=0x%02x", sess.ApType, pkt.ApType))
	}
	if sess.VpType != pkt.VpType {
		s.log.Error(fmt.Sprintf("incorrect visit protocol type, 0x%02x!=0x%02x", sess.VpType, pkt.VpType))
	}
	if sess.McpType != pkt.VpPacket.CPID {
		s.log.Error(fmt.Sprintf("incorrect mcp type, 0x%02x!=0x%02x", sess.McpType, pkt.VpPacket.CPID))
	}
	// if sess.Addr != addr {
	// 	s.log.Warnf(fmt.Sprintf("update device address, Addr:%s>>%s", sess.Addr, addr))
	// 	sess.Addr = addr
	// }
}

func (s *PrivMgmtServer) HandleMessage(addr string, payload []byte) {
	pkt := ApPacket{}
	if err := pkt.UnmarshalBinary(payload); err != nil {
		s.log.Error(errors.Wrap(err, "unmarshal ap packet"))
		errString := err.Error()
		if strings.Contains(errString, "invalid checksum") {
			if err := s.sendResponseWithoutData(addr, pkt, VP_FLAG_RESPONSE_OK, CMD_STATUS_INVALID_CRC); err != nil {
				s.log.Error(errors.Wrap(err, "send error response"))
			}
		} else if strings.Contains(errString, "invalid length") || strings.Contains(errString, "invalid packet length") {
			if err := s.sendResponseWithoutData(addr, pkt, VP_FLAG_RESPONSE_OK, CMD_STATUS_INVALID_LENGTH); err != nil {
				s.log.Error(errors.Wrap(err, "send error response"))
			}
		}
		return
	}
	s.log.Tracef("recv packet: %v", pkt.Dump())

	mcppkt := pkt.VpPacket.McpPacket
	switch mcppkt.CmdID {
	case CMD_ID_NOTIFY:
		if err := s.handleNotifyRequest(pkt, addr); err != nil {
			s.log.Error(errors.Wrap(err, "handle notify request"))
		}
	case CMD_ID_QUERY:
		if err := s.handleResponse(&pkt); err != nil {
			s.log.Error(errors.Wrap(err, "handle query response"))
		}
	case CMD_ID_SET:
		if err := s.handleResponse(&pkt); err != nil {
			s.log.Error(errors.Wrap(err, "handle set response"))
		}
	default:
		if err := s.sendResponse(addr, pkt, VP_FLAG_RESPONSE_OK, CMD_STATUS_INVALID_ID, nil); err != nil {
			s.log.Error(errors.Wrap(err, "send response"))
		}
	}
}

func (s *PrivMgmtServer) handleNotifyRequest(req ApPacket, addr string) error {
	mcppkt := req.VpPacket.McpPacket
	if mcppkt.CmdStatus != CMD_STATUS_COMMAND {
		if err := s.sendResponseWithoutData(addr, req, VP_FLAG_RESPONSE_OK, CMD_STATUS_INVALID_ID); err != nil {
			return errors.Wrap(err, "send error response")
		}
	}
	if err := s.sendResponseWithoutData(addr, req, VP_FLAG_RESPONSE_OK, CMD_STATUS_SUCC); err != nil {
		return errors.Wrap(err, "send notify response")
	}
	return nil
}

func (s *PrivMgmtServer) handleResponse(resp *ApPacket) error {
	key := fmt.Sprintf("%v", resp.VpPacket.Serial)
	if raw, ok := s.waitingResponseMap.Load(key); ok {
		ch, ok := raw.(chan *ApPacket)
		if ok {
			ch <- resp
		}
	}
	return nil
}

func (s *PrivMgmtServer) queryDevices(
	sess *DeviceSession,
	schema string,
	queryString string,
	cb func(total uint8, index uint8, info *model.DeviceInfo),
) ([]*model.DeviceInfo, error) {
	var wg sync.WaitGroup
	defer wg.Wait()

	inputs := []uint8{
		QUERY_DEVICE_INPUT_SUB_ID,
		QUERY_DEVICE_INPUT_ROUTE,
		QUERY_DEVICE_INPUT_IPADDR,
		QUERY_DEVICE_INPUT_CONNECT_STATE,
		QUERY_DEVICE_INPUT_ALARM_STATE,
		QUERY_DEVICE_INPUT_VERSION_STATE,
		QUERY_DEVICE_INPUT_OPTICAL_STATE,
		QUERY_DEVICE_INPUT_DEVICE_TYPE_ID,
		QUERY_DEVICE_INPUT_LOCATION_INFO,
		QUERY_DEVICE_INPUT_DEVICE_TYPE_NAME,
		QUERY_DEVICE_INPUT_VERSION,
		QUERY_DEVICE_INPUT_DEVICE_NAME,
		QUERY_DEVICE_INPUT_ELEMENT_MODEL_NUM,
	}
	infos := []*model.DeviceInfo{}

	if queryString != "" {
		if buf, err := hex.DecodeString(queryString); err == nil {
			_inputs := []uint8{}
			for _, v := range buf {
				_inputs = append(_inputs, uint8(v))
			}
			if len(_inputs) > 0 {
				inputs = _inputs
			}
		}
	}

	data := McpData{}
	data.SetupQueryDevicesRequest(sess.McpType, 1, 1, inputs)

	resp, err := s.sendRequset(sess, CMD_ID_QUERY, data, APC_QUERY_IDS_TIMEOUT)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	total, index, _, outputs, err := resp.VpPacket.McpPacket.Data.ParseQueryDevicesResponse(sess.McpType)
	if err != nil {
		return nil, errors.Wrap(err, "parse response data")
	}
	num := len(outputs)
	for i, output := range outputs {
		if output.DeviceTypeName == "" && output.DeviceTypeID > 0 {
			if _def := model.GetProductDefineByDeviceTypeID(s.opts.DeviceTypePrefix, int(output.DeviceTypeID)); _def != nil {
				output.DeviceTypeName = _def.DeviceTypeName
			}
		}
		info := output.DeviceInfo(schema, sess.MainID)
		infos = append(infos, &info)
		if cb != nil {
			wg.Add(1)
			go func(total uint8, index uint8, info *model.DeviceInfo) {
				defer wg.Done()
				cb(total, index, info)
			}(total, index+uint8(i+1)-uint8(num), &info)
		}
	}

	for {
		// time.Sleep(time.Second)
		index += 1
		if index > total {
			break
		}

		data2 := McpData{}
		data2.SetupQueryDevicesRequest(sess.McpType, total, index, inputs)
		resp, err = s.sendRequset(sess, CMD_ID_QUERY, data2, APC_QUERY_IDS_TIMEOUT)
		if err != nil {
			return nil, errors.Wrap(err, "send request")
		}

		total2, index2, _, outputs, err := resp.VpPacket.McpPacket.Data.ParseQueryDevicesResponse(sess.McpType)
		if err != nil {
			return nil, errors.Wrap(err, "parse response data")
		}
		total = total2
		// if total2 != total {
		// 	return nil, errors.New("total number of query devices is changed")
		// }
		if index2 != index {
			return nil, errors.New("index number of query devices is changed")
		}

		num := len(outputs)
		for i, output := range outputs {
			if output.DeviceTypeName == "" && output.DeviceTypeID > 0 {
				if _def := model.GetProductDefineByDeviceTypeID(s.opts.DeviceTypePrefix, int(output.DeviceTypeID)); _def != nil {
					output.DeviceTypeName = _def.DeviceTypeName
				}
			}
			info := output.DeviceInfo(schema, sess.MainID)
			infos = append(infos, &info)
			if cb != nil {
				wg.Add(1)
				go func(total uint8, index uint8, info *model.DeviceInfo) {
					defer wg.Done()
					cb(total, index, info)
				}(total, index+uint8(i+1)-uint8(num), &info)
			}
		}
	}
	return infos, nil
}

func (s *PrivMgmtServer) queryAllObjectIds(sess *DeviceSession) error {
	data := McpData{}
	data.SetupQueryObjectIdsRequest(sess.McpType, 1, 1)
	resp, err := s.sendRequset(sess, CMD_ID_QUERY, data, APC_QUERY_IDS_TIMEOUT)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	total, index, ids, err := resp.VpPacket.McpPacket.Data.ParseQueryObjectIdsResponse(sess.McpType)
	if err != nil {
		return errors.Wrap(err, "parse response data")
	}

	for {
		time.Sleep(time.Millisecond * 100)
		index += 1
		if index >= total {
			break
		}

		data2 := McpData{}
		data2.SetupQueryObjectIdsRequest(sess.McpType, total, index)
		resp, err = s.sendRequset(sess, CMD_ID_QUERY, data2, APC_QUERY_IDS_TIMEOUT)
		if err != nil {
			return errors.Wrap(err, "send request")
		}
		total2, index2, ids2, err := resp.VpPacket.McpPacket.Data.ParseQueryObjectIdsResponse(sess.McpType)
		if err != nil {
			return errors.Wrap(err, "parse response data")
		}
		if total2 != total {
			return errors.New("total number of query object ids is changed")
		}
		if index2 != index {
			return errors.New("index number of query devices is changed")
		}
		ids = append(ids, ids2...)
	}

	sort.Sort(ObjectIDSort(ids))
	// s.log.Warnf("%v", ids)
	return nil
}

func (s *PrivMgmtServer) queryObjects(sess *DeviceSession, cmdId string, objects []Object) ([]Object, error) {
	_cmdId := GetQueryCommandId(cmdId)

	data := McpData{}
	if err := data.SetupObjects(sess.McpType, _cmdId, objects, sess.Params); err != nil {
		return nil, errors.Wrap(err, "setup objects")
	}

	resp, err := s.sendRequset(sess, _cmdId, data, APC_QUERY_ONE_TIMEOUT)

	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	objects2, err := resp.VpPacket.McpPacket.Data.ParseObjects(sess.McpType, _cmdId, sess.Params)
	if err != nil {
		return nil, errors.Wrap(err, "parse objects")
	}
	return objects2, nil
}

func (s *PrivMgmtServer) setObjects(sess *DeviceSession, cmdId string, objects []Object) ([]Object, error) {
	_cmdId := GetSetCommandId(cmdId)

	data := McpData{}
	if err := data.SetupObjects(sess.McpType, _cmdId, objects, sess.Params); err != nil {
		return nil, errors.Wrap(err, "setup objects")
	}

	resp, err := s.sendRequset(sess, _cmdId, data, APC_QUERY_ONE_TIMEOUT)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	objects2, err := resp.VpPacket.McpPacket.Data.ParseObjects(sess.McpType, _cmdId, sess.Params)
	if err != nil {
		return nil, errors.Wrap(err, "parse objects")
	}
	return objects2, nil
}
