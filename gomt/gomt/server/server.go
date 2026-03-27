package server

import (
	"crypto/tls"
	"fmt"
	"gomt/core/audit"
	"gomt/core/certgen"
	"gomt/core/model"
	"gomt/core/proto/priv"
	"gomt/das"
	"gomt/das/agent"
	"gomt/das/file"
	"gomt/das/system"
	"io/fs"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

type ServerOptions struct {
	Schema        string
	ServerAddress string
	ServerPort    uint16

	GzipEnable bool
	GzipLevel  int

	TlsEnable   bool
	TlsCertFile string
	TlsKeyFile  string

	RequestRate      int
	RequestExpiresIn time.Duration

	CgiRequestRate      int
	CgiRequestExpiresIn time.Duration

	RootFS            fs.FS
	ConfigDir         string
	CgiBinDir         string
	FileTypes         file.FileTypes
	MaxFileUploadSize int64

	DeviceTypeName          string
	DeviceAddr              string
	DeviceInternalInterface string
	DevicePort              uint16

	AuditLogFile       string
	AuditLogMaxSize    int
	AuditLogMaxBackups int

	MetricsDataFile       string
	MetricsDataMaxSize    int
	MetricsDataMaxBackups int
}

type OMTServer struct {
	log           *logrus.Entry
	opts          ServerOptions
	localhostMode bool

	auditLogger            *audit.AuditLogger
	e                      *echo.Echo
	dasDeviceProxyBalancer *deviceProxyBalancer
	dassys                 *system.DasSystem
	srv                    *priv.PrivMgmtServer

	db                     *gorm.DB
	enforcerDriverName     string
	enforcerDataSourceName string
	enforce                *casbin.Enforcer

	certGenOption certgen.CertGenOptions
	certResult    *certgen.CertResult

	wsUpgrader websocket.Upgrader
	wsLock     sync.Mutex
	wsConns    []*websocket.Conn

	queryDevicesLock     sync.RWMutex
	queryDevicesRunning  bool
	queryDevicesTime     time.Time
	queryDevicesProgress QueryDevicesProgress
	queryDevicesInterval time.Duration

	wg       sync.WaitGroup
	quit     chan bool
	shutdown bool
}

func NewServer(options ServerOptions) (*OMTServer, error) {
	s := &OMTServer{
		// deviceTypeName: options.DeviceTypeName,
		opts: options,
		quit: make(chan bool),
	}
	s.localhostMode = options.DeviceAddr == "127.0.0.1"
	s.log = logrus.WithFields(logrus.Fields{"localhost": s.localhostMode})
	s.auditLogger = audit.NewAuditLogger("omt", s.opts.AuditLogFile, s.opts.AuditLogMaxSize, s.opts.MetricsDataMaxBackups)

	s.srv = priv.GetDefaultPrivMgmtServer()
	s.dassys = system.NewDasSystem()

	s.e = echo.New()
	s.e.HideBanner = true
	s.e.Use(s.loggerMiddlewareFunc())
	s.e.Use(middleware.Recover())

	s.e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// s.e.Use(middleware.CORS())

	if s.opts.GzipEnable {
		skipPaths := []string{"/cgi-bin/", "/config/", "/Version/"}
		s.e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				req := c.Request()
				// logrus.Warnf("%v", req.URL.Path)
				for _, v := range skipPaths {
					if strings.HasPrefix(req.URL.Path, v) {
						return true
					}
				}
				return false
			},
			Level: s.opts.GzipLevel,
		}))
	}

	limiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(s.opts.RequestRate),
			Burst:     s.opts.RequestRate * 3,
			ExpiresIn: s.opts.RequestExpiresIn,
		}),
	}
	s.e.Use(middleware.RateLimiterWithConfig(limiterConfig))

	s.dasDeviceProxyBalancer = NewDeviceProxyBalancer(s.dassys)
	deviceProxy := middleware.ProxyWithConfig(
		middleware.ProxyConfig{
			Balancer: s.dasDeviceProxyBalancer,
			RegexRewrite: map[*regexp.Regexp]string{
				regexp.MustCompile(`^\/proxy\/(local|[0-9]+)$`):                                                                "",
				regexp.MustCompile(`^\/proxy\/(local|[0-9]+)\/$`):                                                              "/",
				regexp.MustCompile(`^\/proxy\/(local|[0-9]+)\/(.*)\.(html|ico|properties)$`):                                   "/$2.$3",
				regexp.MustCompile(`^\/proxy\/(local|[0-9]+)\/(v-js|js|css|img|static|cgi-bin|config|models)\/(.*)$`):          "/$2/$3",
				regexp.MustCompile(`^\/proxy\/(local|[0-9]+)\/(ConfigFiles|LogFiles|UploadFiles|Vertion|packetUpdate)\/(.*)$`): "/$2/$3",
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	)

	dasDeviceProxyGroup := s.e.Group("/proxy/:device_sub")
	dasDeviceProxyGroup.Use(
		deviceProxy,
	)

	s.wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  102400,
		WriteBufferSize: 102400,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	//certgen
	s.certGenOption.CommonName = s.opts.ServerAddress
	cert, err := certgen.GetCertResult(&s.certGenOption)
	if err != nil {
		return nil, errors.Wrapf(err, "get cert block")
	}
	s.certResult = cert

	// api
	iamMiddlewares := []echo.MiddlewareFunc{
		s.makeLoginHandler("/api"),
	}

	s.e.GET("/public.crt", s.handleGetCAPublic)

	apiGroup := s.e.Group("/api")
	appGroup := apiGroup.Group("/app")
	appGroup.GET("/info", s.handleGetAppInfo)

	apiDiagGroup := apiGroup.Group("/diag")
	apiDiagGroup.Use(iamMiddlewares...)
	apiPingDiagGroup := apiDiagGroup.Group("/ping")
	apiPingDiagGroup.POST("/jobs", s.handleCreatePingJob)
	//apiPingGroup.GET("/jobs/:token", s.handleGetPingJob)
	apiPingDiagGroup.GET("/jobs/:token/ws", s.handlePingJobWebSocket)
	apiPingDiagGroup.POST("/jobs/:token/run", s.handleRunPingJob)
	apiPingDiagGroup.POST("/jobs/:token/cancel", s.handleCancelPingJob)

	//auth
	authGroup := apiGroup.Group("/auth")
	authGroup.Use(iamMiddlewares...)
	authGroup.POST("/login", s.handleLogin)
	authGroup.GET("/logout", s.handleLogout)

	// users
	currentGroup := apiGroup.Group("/current")
	currentGroup.Use(iamMiddlewares...)
	currentGroup.GET("", s.handleGetCurrentUser)
	currentGroup.POST("/change-password", s.handleChangeCurrentUserPassword)

	// roles
	roleGroup := apiGroup.Group("/iam/roles")
	roleGroup.Use(iamMiddlewares...)
	roleGroup.GET("", s.handleListRoles)
	// users
	userGroup := apiGroup.Group("/iam/users")
	userGroup.Use(iamMiddlewares...)
	userGroup.GET("", s.handleListUsers)
	userGroup.POST("", s.handleCreateUser)
	userGroup.DELETE("/:name", s.handleDeleteUser)
	userGroup.POST("/:name", s.handleUpdateUser)
	userGroup.POST("/:name/set-password", s.handleUpdateUserPassword)

	apiDasGroup := apiGroup.Group("/das")
	apiDasGroup.Use(iamMiddlewares...)
	apiDasGroup.GET("/ws", s.handleDasWebSocket)
	apiDasGroup.GET("/products", s.handleListDasProductModels)
	apiDasGroup.GET("/products/:name", s.handleGetDasProductModel)

	apiDasGroup.GET("/device-types", s.handleGetDasDeviceTypes)
	apiDasGroup.GET("/device-types/:name/:version/model/layout", s.handleGetDasModelLayout)
	apiDasGroup.GET("/device-types/:name/:version/model/parameters", s.handleGetDasModelParamters)

	apiDasGroup.GET("/query-devices/progress", s.handleGetDasQueryDevicesProgress)

	apiDasDevicesGroup := apiDasGroup.Group("/devices")
	apiDasDevicesGroup.GET("", s.handleGetDasDeviceInfos)

	apiDasDeviceGroup := apiDasDevicesGroup.Group("/:device_sub")
	apiDasDeviceGroup.Use(iamMiddlewares...)
	apiDasDeviceGroup.Use(s.agentAvailableMiddleware)
	apiDasDeviceGroup.GET("", s.handleGetDasDeviceInfo)
	apiDasDeviceGroup.GET("/type", s.handleDasDeviceGetTypeInfo)
	apiDasDeviceGroup.GET("/available", s.handleDasDeviceCheckAvailable)

	apiDasDeviceGroup.POST("/parameters/get", s.handleDasDeviceGetParameterValues)
	apiDasDeviceGroup.POST("/parameters/set", s.handleDasDeviceSetParameterValues)
	apiDasDeviceGroup.POST("/register/read", s.handleDasDeviceReadRegister)
	apiDasDeviceGroup.POST("/register/write", s.handleDasDeviceWriteRegister)

	apiDasDeviceGroup.POST("/metrics/data/query", s.handleDasDeviceQueryMetricsData)
	apiDasDeviceGroup.GET("/metrics/current", s.handleeDasDeviceGetMetricsCurrent)

	apiDasDeviceGroup.GET("/files/:ftype", s.handleDasDeviceListFiles)
	apiDasDeviceGroup.POST("/files/:ftype", s.handleDasDeviceUploadFile)
	apiDasDeviceGroup.GET("/files/:ftype/:fname", s.handleDasDeviceGetFile)
	apiDasDeviceGroup.DELETE("/files/:ftype/:fname", s.handleDasDeviceDeleteFile)
	apiDasDeviceGroup.GET("/files/:ftype/:fname/packet-info", s.handleDasDeviceGetUpgradeFilePacketInfo)
	apiDasDeviceGroup.GET("/version/packet-info", s.handleDasDeviceGetVersionPacketInfo)

	apiDasDeviceGroup.GET("/alarm/logs", s.handleDasDeviceListAlarmEventLogs)

	apiDasDeviceGroup.POST("/upgrade/start", s.handleDasDeviceStartUpgrade)
	apiDasDeviceGroup.POST("/upgrade/reboot", s.handleDasDeviceSetUpgradeReboot)
	apiDasDeviceGroup.GET("/firmwares", s.handleDasDeviceListFirmwares)
	apiDasDeviceGroup.DELETE("/firmwares/:name", s.handleDasDeviceDeleteFirmware)

	apiDasDeviceGroup.POST("/delete-key-and-logs", s.handleDasDeviceDeleteKeyAndLogs)

	// local
	if s.localhostMode {
		cgiGroup := s.e.Group("/cgi-bin")
		cgiLimiterConfig := middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(s.opts.CgiRequestRate),
				Burst:     s.opts.CgiRequestRate * 3,
				ExpiresIn: s.opts.CgiRequestExpiresIn,
			}),
		}
		cgiGroup.Use(middleware.RateLimiterWithConfig(cgiLimiterConfig))
		cgiGroup.GET("/*", s.makeLocalCgiHandlerFunc(s.opts.CgiBinDir))
		cgiGroup.POST("/*", s.makeLocalCgiHandlerFunc(s.opts.CgiBinDir))
	}

	// root
	s.e.StaticFS("/*", s.opts.RootFS)
	return s, nil
}

func (s *OMTServer) Stop() {
	s.log.Trace("stop")
	s.shutdown = true
	s.quit <- true
}

func (s *OMTServer) Run() {
	s.log.Trace("running")
	defer s.log.Trace("stopped")
	defer s.wg.Wait()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.srv.Run()
	}()

	time.Sleep(time.Second * 1)

	if err := s.setupLocalAgent(); err != nil {
		s.log.Fatal(errors.Wrap(err, "setup local device agent"))
	}
	s.log.Tracef("setup local agent")

	s.updateDasDevices(true)
	s.log.Tracef("setup das devices")

	if err := s.setupIAM(); err != nil {
		s.log.Fatal(errors.Wrap(err, "setup iam"))
	}
	s.log.Tracef("setup iam")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		addr := fmt.Sprintf("%v:%v", s.opts.ServerAddress, s.opts.ServerPort)

		if s.opts.TlsEnable {
			if s.opts.TlsCertFile != "" && s.opts.TlsKeyFile != "" {
				if err := s.e.StartTLS(addr, s.opts.TlsCertFile, s.opts.TlsKeyFile); err != nil && err != http.ErrServerClosed {
					s.e.Logger.Fatal(err)
				}
			} else if s.certResult != nil {
				if err := s.e.StartTLS(addr, s.certResult.CertPEM, s.certResult.PrivPEM); err != nil && err != http.ErrServerClosed {
					s.e.Logger.Fatal(err)
				}
			} else {
				if err := s.e.Start(addr); err != nil && err != http.ErrServerClosed {
					s.e.Logger.Fatal(err)
				}
			}
		} else {
			if err := s.e.Start(addr); err != nil && err != http.ErrServerClosed {
				s.e.Logger.Fatal(err)
			}
		}
	}()

	<-s.quit

	s.e.Close()
	s.srv.Stop()
}

func (s *OMTServer) setupLocalAgent() error {
	deviceSub := "local"
	opts := agent.DasDeviceAgentOptions{
		DeviceAddr:        s.opts.DeviceAddr,
		PrivServerPort:    s.opts.DevicePort,
		HttpServerAddr:    das.MakeDeviceUrlBase(s.opts.DeviceAddr, 443, true),
		CGIBinUrlPath:     "/cgi-bin",
		CGIBinFilePath:    s.opts.CgiBinDir,
		FileTypes:         s.opts.FileTypes,
		ConfigPath:        s.opts.ConfigDir,
		MetricsFilename:   s.opts.MetricsDataFile,
		MetricsMaxSize:    s.opts.MetricsDataMaxSize,
		MetricsMaxBackups: s.opts.MetricsDataMaxBackups,
		InternalInterface: s.opts.DeviceInternalInterface,
	}
	a, err := s.dassys.SetupDasDeviceAgent(s.opts.Schema, deviceSub, nil, opts)
	if err != nil {
		return err
	}
	info := a.GetDeviceInfo()
	s.log.Infof("deviceType=%v subID=%v routeAddr=%v ip=%v",
		info.DeviceTypeName, info.SubID, info.RouteAddressString, info.IpAddressString)

	for {
		s.log.Infof("waiting for local service available...")
		if ok := a.IsServiceAvailable(true); ok {
			break
		}
		time.Sleep(time.Second)
	}
	if err := a.UpdateDeviceInfo(); err != nil {
		s.log.Error(errors.Wrap(err, "update local device info"))
	}

	s.log.Infof("local service is available")

	if a.IsLocalDevice() == false {
		a.SetParameterValueOfJumpEnable(true)
	}
	return nil
}

func (s *OMTServer) setupDasDeviceAgent(info *model.DeviceInfo) (*agent.DasDeviceAgent, error) {
	if info == nil {
		return nil, errors.New("invalid device info")
	}
	opts := agent.DasDeviceAgentOptions{
		PrivServerPort:    s.opts.DevicePort,
		CGIBinUrlPath:     "/cgi-bin",
		CGIBinFilePath:    s.opts.CgiBinDir,
		FileTypes:         s.opts.FileTypes,
		ConfigPath:        s.opts.ConfigDir,
		MetricsFilename:   s.opts.MetricsDataFile,
		MetricsMaxSize:    s.opts.MetricsDataMaxSize,
		MetricsMaxBackups: s.opts.MetricsDataMaxBackups,
	}
	if s.localhostMode {
		opts.DeviceAddr = info.IpAddressString
		opts.HttpServerAddr = das.MakeDeviceUrlBase(info.IpAddressString, 443, true)
	} else {
		opts.DeviceAddr = s.opts.DeviceAddr
		opts.HttpServerAddr = das.MakeDeviceUrlBase(s.opts.DeviceAddr, info.ForwardingPort, true)
	}
	deviceSub := fmt.Sprintf("%v", info.SubID)
	a, err := s.dassys.SetupDasDeviceAgent(s.opts.Schema, deviceSub, info, opts)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *OMTServer) handleGetCAPublic(c echo.Context) error {
	return c.Blob(http.StatusOK, "application/octet-stream", s.certResult.CaPEM)
}
