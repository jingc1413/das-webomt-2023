package agent

import (
	"encoding/hex"
	"fmt"
	"gomt/core/metric"
	"gomt/core/model"
	"gomt/core/proto/priv"
	"gomt/core/utils"
	"gomt/das/arm"
	"gomt/das/cgi"
	"gomt/das/file"
	"gomt/das/um"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type DasDeviceAgentOptions struct {
	LocalAddr string

	DeviceAddr        string
	InternalInterface string
	HttpServerAddr    string
	PrivServerPort    uint16
	CGIBinFilePath    string
	CGIBinUrlPath     string
	ConfigPath        string
	MetricsFilename   string
	MetricsMaxSize    int
	MetricsMaxBackups int

	FileTypes file.FileTypes
}

type DasDeviceAgent struct {
	log *logrus.Entry

	sync.Mutex
	info           *model.DeviceInfo
	schema         string
	deviceTypeName string
	version        string
	deviceSub      string
	isLocalDevice  bool
	defs           model.ParameterDefines
	fileTypes      file.FileTypes

	serviceAvailableLock      sync.RWMutex
	serviceAvailable          bool
	serviceAvailableCheckTime time.Time

	supportPriv         bool
	supportCGI          bool
	supportAPI          bool
	supportArmIni       bool
	deviceAddr          string
	deviceUrlBase       string
	deviceCgiBinUrlBase string
	deviceApiUrlBase    string
	deviceConfigPath    string
	deviceFileTypes     file.FileTypes

	metricsFilename   string
	metricsMaxsize    int
	metricsMaxbackups int

	privSrv    *priv.PrivMgmtServer
	privSess   *priv.DeviceSession
	cgiHandler *cgi.CGIHandler
	fileMgmt   *file.FileMgmt
	userMgmt   *um.UserManager

	metrics *metric.MetricsData

	wg       sync.WaitGroup
	quit     chan bool
	shutdown bool
}

func NewDasDeviceAgent(schema string, deviceSub string, info *model.DeviceInfo, opts DasDeviceAgentOptions) (*DasDeviceAgent, error) {
	s := &DasDeviceAgent{
		deviceSub:         deviceSub,
		schema:            schema,
		info:              info,
		metricsFilename:   opts.MetricsFilename,
		metricsMaxsize:    opts.MetricsMaxSize,
		metricsMaxbackups: opts.MetricsMaxBackups,
		quit:              make(chan bool),
	}
	s.privSrv = priv.GetDefaultPrivMgmtServer()
	s.log = logrus.WithFields(logrus.Fields{"deviceSub": s.deviceSub})

	if err := s.setupDevice(opts); err != nil {
		// s.log.Error(errors.Wrap(err, "setup device agent"))
		return nil, err
	}
	return s, nil
}

func (s *DasDeviceAgent) Close() {
	defer s.log.Tracef("close")
	s.shutdown = true
	s.quit <- true
}

func (s *DasDeviceAgent) Run() {
	s.log.Trace("running")
	defer s.log.Trace("stopped")
	defer s.wg.Wait()

	if s.metricsFilename != "" && s.deviceSub == "local" {
		s.setupMetrics(s.metricsFilename, s.metricsMaxsize, s.metricsMaxbackups)

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for {
				if s.shutdown {
					return
				}
				start := time.Now()
				s.updateMetrics()
				diff := time.Since(start)
				if diff > 0 && diff < time.Second {
					time.Sleep(time.Second - diff)
				}
			}
		}()
	}

	<-s.quit
}

func (s *DasDeviceAgent) GetDeviceSub() string {
	return s.deviceSub
}

func (s *DasDeviceAgent) GetProductTypeName() string {
	if s.info != nil {
		return s.info.ProductType
	}
	return ""
}

func (s *DasDeviceAgent) GetDeviceTypeName() string {
	return s.deviceTypeName
}

func (s *DasDeviceAgent) IsLocalDevice() bool {
	return s.isLocalDevice
}

func (s *DasDeviceAgent) setDeviceTypeName(deviceTypeName string) {
	s.deviceTypeName = deviceTypeName
	s.info.DeviceTypeName = deviceTypeName
	if s.deviceTypeName != "" {
		s.defs = model.GetParameterDefinesByDeviceTypeName(s.schema, s.info.DeviceTypeName, s.info.Version)
	}
	s.info.Setup()
	s.log = logrus.WithFields(logrus.Fields{"deviceType": s.deviceTypeName, "deviceSub": s.deviceSub})
}

func (s *DasDeviceAgent) setupDevice(opts DasDeviceAgentOptions) error {
	s.deviceAddr = opts.DeviceAddr
	s.deviceConfigPath = opts.ConfigPath
	s.deviceFileTypes = opts.FileTypes
	s.deviceUrlBase = opts.HttpServerAddr
	s.deviceApiUrlBase = s.deviceUrlBase + "/api/das/devices/local"
	s.deviceCgiBinUrlBase = s.deviceUrlBase + "/cgi-bin"

	s.fileTypes = opts.FileTypes

	s.isLocalDevice = s.deviceAddr == "127.0.0.1"
	if s.isLocalDevice {
		s.supportPriv = true
		s.supportAPI = false
		s.supportCGI = true //das.IsSupportCGI(s.deviceUrlBase)
		s.supportArmIni = arm.IsSupportArmIni(s.deviceConfigPath)
	} else {
		s.supportPriv = false
		s.supportAPI = false //das.IsSupportAPI(s.deviceUrlBase)
		s.supportCGI = true  //das.IsSupportCGI(s.deviceUrlBase)
		s.supportArmIni = false
	}

	if s.info != nil {
		s.setDeviceTypeName(s.info.DeviceTypeName)
	} else {
		s.info = &model.DeviceInfo{
			ConnectState: 1,
		}
		if s.isLocalDevice {
			ip, err := utils.GetIpByInterfaceName(opts.InternalInterface)
			if err != nil {
				return errors.Wrap(err, "get ip address of internal interface")
			}
			s.info.IpAddress = ip.To4()
		}
		s.info.Setup()
	}

	if s.deviceTypeName == "" {
		v, err := s.getDeviceTypeName()
		if err != nil {
			return errors.Wrap(err, "get device type name")
		}
		s.setDeviceTypeName(v)
	}
	if s.deviceTypeName == "" {
		return errors.New("device type name is null")
	}

	if s.IsConnectState() {
		if err := s.setupHanlers(opts); err != nil {
			return errors.Wrap(err, "setup handlers")
		}
		if version, err := s.getSoftwareVersion(); err != nil {
			return errors.Wrap(err, "get software version")
		} else if version != s.version {
			s.version = version
			if err := s.setupHanlers(opts); err != nil {
				return errors.Wrap(err, "setup handlers")
			}
		}
	}

	return nil
}

func (s *DasDeviceAgent) setupHanlers(opts DasDeviceAgentOptions) error {
	if s.isLocalDevice {
		s.cgiHandler = cgi.NewLocalCGIHandler(s.schema, s.deviceTypeName, s.version, s.defs, opts.CGIBinFilePath)
		s.fileMgmt = file.NewLocalFileMgmt(opts.FileTypes)
		s.userMgmt = um.NewUserManager(s.deviceConfigPath)
	} else {
		s.cgiHandler = cgi.NewRemoteCGIHandler(s.schema, s.deviceTypeName, s.version, s.defs, s.deviceCgiBinUrlBase)
		s.fileMgmt = file.NewRemoteFileMgmt(opts.FileTypes, s.cgiHandler, s.deviceUrlBase)
	}
	if s.isLocalDevice && s.supportPriv {
		mainID, err := hex.DecodeString("00000000")
		if err != nil {
			return errors.Wrap(err, "invalid device id")
		}
		subID := cast.ToUint8(s.deviceSub)
		privAddr := fmt.Sprintf("%v:%v", s.deviceAddr, opts.PrivServerPort)
		sess, err := s.privSrv.RegisterDevice(
			priv.AP_C, priv.VP_A1, priv.MCP_A,
			s.deviceTypeName, s.version, s.defs,
			mainID, subID, privAddr)
		if err != nil {
			return errors.Wrap(err, "register device")
		}
		s.privSess = sess
		// } else {
		// privAddr := fmt.Sprintf("%v:%v", s.deviceAddr, opts.PrivServerPort)
		// sess, err := s.privSrv.RegisterDevice(priv.AP_C, priv.VP_A, priv.MCP_A,
		// 	info.DeviceTypeName, info.MainID, info.SubID, privAddr)
		// if err != nil {
		// 	return errors.Wrap(err, "register device")
		// }
		// s.privSess = sess
	}

	s.version = s.info.Version

	s.log.Tracef("type=%v, version=%v, sub=%v, addr=%v, cgi=%v, priv=%v, api=%v",
		s.deviceTypeName, s.version, s.deviceSub,
		s.deviceAddr, s.supportCGI, s.supportPriv, s.supportAPI,
	)
	return nil
}
