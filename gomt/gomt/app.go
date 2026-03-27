package main

import (
	"embed"
	"flag"
	"fmt"
	"gomt/core"
	"gomt/core/layout"
	"gomt/core/logger/lumberjack"
	"gomt/core/model"
	"gomt/core/proto/priv"
	"gomt/core/utils"
	"gomt/core/watcher"
	"gomt/das/file"
	"gomt/gomt/server"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:embed dist
var embededFiles embed.FS

func main() {
	version := flag.Bool("version", false, "version")
	verbose := flag.Bool("verbose", false, "verbose mode")

	serverSchema := flag.String("schema", "default", "schema")
	serverInterface := flag.String("inf", "", "server interface")
	serverAddr := flag.String("addr", "", "server address")
	serverPort := flag.Uint("port", 0, "server port (default http=80 https=443)")

	privPort := flag.Int("priv-server-port", 9739, "private manamgent server udp port")
	privWorker := flag.Int("priv-server-worker", 10, "private manamgent server worker number")

	deviceTypeName := flag.String("device-type", "", "device type name")
	deviceAddr := flag.String("device-addr", "127.0.0.1", "device address")
	deviceInternalInterface := flag.String("device-inter-inf", "eth0", "device internal interface")
	deviceLocalPort := flag.Uint("device-local-port", 6001, "device local priv udp server port")

	omtLogFile := flag.String("omtlog-file", "/drgfly/webomt/web/www/LogFiles/OmtLog/omt.log", "omt log file")
	omtLogMaxSize := flag.Int("omtlog-maxsize", 200, "omt log file maximum size(kilo bytes)")
	omtLogMaxBackups := flag.Int("omtlog-backups", 5, "omt log file maximum backups")

	metricsDataFile := flag.String("metrics-file", "/drgfly/webomt/web/www/LogFiles/Metrics/metrics.csv", "metrics data file")
	metricsDataMaxSize := flag.Int("metrics-maxsize", 200, "metrics data file maximum size(kilo bytes)")
	metricsDataMaxBackups := flag.Int("metrics-backups", 5, "metrics data file maximum backups")

	auditLogFile := flag.String("auditlog-file", "/drgfly/webomt/web/www/LogFiles/AuditLog/audit.log", "audit log file")
	auditLogMaxSize := flag.Int("auditlog-maxsize", 200, "audit log file maximum size(kilo bytes)")
	auditLogMaxBackups := flag.Int("auditlog-backups", 5, "audit log file maximum backups")

	configPath := flag.String("config-path", "/drgfly/config", "path to the config folder")
	wwwPath := flag.String("www-path", "/drgfly/webomt/web/www", "path to the root folder")
	logPath := flag.String("log-path", "/drgfly/webomt/web/www/LogFiles", "path to the log folder")
	logTypes := flag.String("log-types", "GuardLog,WebLog,DeviceLog,LoginLog", "log types")
	fileTypes := flag.String("file-types",
		"Version|/drgfly/webomt/web/www/Vertion|.txt;"+
			"Config|/drgfly/webomt/web/www/config|.xml,.txt;"+
			"PacketUpdateFile|/drgfly/webomt/web/www/packetUpdate|.txt;"+
			"UpgradeFile|/drgfly/webomt/web/www/UploadFiles|.zip;"+
			"ConfigFile|/drgfly/webomt/web/www/ConfigFiles|.json,.csv",
		"file types: {name}|{dir}|{exts};{name}|{dir}|{exts}...")
	maxUploadSize := flag.Int64("max-upload-size", 83886080, "max file upload size")

	cgiBinPath := flag.String("cgi-bin", "/drgfly/webomt/web/www/cgi-bin", "path to the cgi-bin folder")
	gzipEnable := flag.Bool("gzip-enable", false, "gzip enable")
	gzipLevel := flag.Int("gzip-level", 5, "gzip compress level")
	tlsEnable := flag.Bool("tls-enable", true, "TLS enable")
	certFile := flag.String("cert-file", "", "TLS certificate file")
	keyFile := flag.String("key-file", "", "TLS key file")
	reqLimit := flag.Int("req-limit", 64, "maximum number of requests per second")
	cgiLimit := flag.Int("cgi-limit", 16, "maximum number of cgi requests per second")
	requestTimeout := flag.Int("request-timeout", 60, "timeout duration in seconds for request")
	flag.Parse()

	if version != nil && *version {
		fmt.Printf("Version: %s\n", core.VERSION)
		fmt.Printf("Build: %s\n", core.BUILD)
		return
	}
	{
		//check process pid
		watcher.CheckRunning(os.Args[0])
		watcher.SaveLockFile(os.Args[0], os.Getpid())
	}
	isLocalDevice := *deviceAddr == "127.0.0.1"
	if !isLocalDevice {
		*omtLogFile = "./OmtLog/omt.log"
		*auditLogFile = "./AuditLog/audit.log"
		*metricsDataFile = "./Metrics/metrics.log"
	}

	if verbose != nil && *verbose {
		logrus.SetLevel(logrus.TraceLevel)
	}
	var lumberjackLogger *lumberjack.Logger
	if isLocalDevice {
		lumberjackLogger = &lumberjack.Logger{
			Filename:   *omtLogFile,
			MaxSize:    *omtLogMaxSize, // kilobyte
			MaxBackups: *omtLogMaxBackups,
		}
		logrus.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))
	}

	schema := "default"
	if serverSchema != nil {
		schema = *serverSchema
	}

	addr := ""
	if serverAddr != nil && *serverAddr != "" {
		addr = *serverAddr
	} else if serverInterface != nil {
		if *serverInterface == "" {
			if isLocalDevice {
				*serverInterface = "eth1"
			}
		}
		if *serverInterface != "" {
			if ip, err := utils.GetIpByInterfaceName(*serverInterface); err == nil {
				addr = ip.String()
			}
		}
	}
	if addr == "" {
		addr = "127.0.0.1"
	}

	port := uint16(0)
	if serverPort != nil && *serverPort != 0 {
		port = uint16(*serverPort)
	} else if *tlsEnable {
		port = 443
	}

	timeoutDuration := time.Duration(*requestTimeout) * time.Second

	logrus.Infof("device: %v", *deviceAddr)
	logrus.Infof("server: %v:%v", addr, port)

	var rootFs fs.FS
	if !isLocalDevice {
		if tmp, err := fs.Sub(embededFiles, "dist"); err != nil {
			logrus.Fatal(err)
		} else {
			rootFs = tmp
		}
	} else {
		file.MustDir(*cgiBinPath)
		file.MustDir(*wwwPath)
		file.MustDir(*logPath)
		err := os.Chdir(*cgiBinPath)
		if err != nil {
			logrus.Fatal(err)
		}
		rootFs = os.DirFS(*wwwPath)
	}

	ftypes := []*file.FileType{}
	if args := strings.Split(*logTypes, ","); len(args) > 0 {
		for _, typ := range args {
			if v := strings.TrimSpace(typ); v != "" {
				filename := filepath.Join(*logPath, v)
				filename = filepath.ToSlash(filepath.Clean(filename))
				ftype := &file.FileType{
					Name: v,
					Dir:  filename,
					Exts: ".log,.csv,.txt",
				}
				ftypes = append(ftypes, ftype)
			}
		}
	}
	if args := strings.Split(*fileTypes, ";"); len(args) > 0 {
		for _, arg := range args {
			if args2 := strings.Split(arg, "|"); len(args2) == 3 {
				name := strings.TrimSpace(args2[0])
				dir := strings.TrimSpace(args2[1])
				exts := strings.TrimSpace(args2[2])
				if isLocalDevice {
					if name == "PacketUpdateFile" {
						continue
					}
					file.MustDir(dir)
				}
				if name == "" || exts == "" {
					logrus.Fatalf("invalid file type, %v", arg)
				}
				ftype := &file.FileType{
					Name:          name,
					Dir:           dir,
					Exts:          exts,
					SupportUpload: true,
				}
				if name == "UpgradeFile" {
					ftype.FieldName = "file_upload"
				} else if name == "ConfigFile" {
					ftype.FieldName = "file_config"
				}
				if ftype.Name == "Config" {
					ftype.Path = "/config"
				} else if ftype.Name == "Version" {
					ftype.Path = "/Vertion"
				} else if ftype.Name == "PacketUpdateFile" {
					ftype.Path = "/packetUpdate"
				}
				ftypes = append(ftypes, ftype)
			}
		}
	}

	if err := model.SetupAllParameterDefineMap(rootFs, fmt.Sprintf("models/%v", schema)); err != nil {
		logrus.Fatal(errors.Wrap(err, "setup device model parametrs"))
	}
	allModels := model.GetAllParameterDefinesMap()
	if err := layout.SetupAllLayoutMap(schema, allModels); err != nil {
		logrus.Fatal(errors.Wrap(err, "setup device model layouts"))
	}
	privateOpts := priv.PrivMgmtServerOptions{
		Schema:     schema,
		Addr:       fmt.Sprintf(":%v", *privPort),
		MaxWorkers: *privWorker,
	}
	if err := priv.SetupDefaultPrivMgmtServer(privateOpts); err != nil {
		logrus.Fatal(errors.Wrap(err, "setup private management server"))
	}

	opts := server.ServerOptions{
		Schema:              schema,
		ServerAddress:       addr,
		ServerPort:          port,
		GzipEnable:          *gzipEnable,
		GzipLevel:           *gzipLevel,
		TlsEnable:           *tlsEnable,
		TlsCertFile:         *certFile,
		TlsKeyFile:          *keyFile,
		RequestRate:         *reqLimit,
		RequestExpiresIn:    timeoutDuration,
		CgiRequestRate:      *cgiLimit,
		CgiRequestExpiresIn: timeoutDuration,
		RootFS:              rootFs,
		CgiBinDir:           *cgiBinPath,
		ConfigDir:           *configPath,
		FileTypes:           ftypes,
		MaxFileUploadSize:   *maxUploadSize,

		DeviceTypeName:          *deviceTypeName,
		DeviceAddr:              *deviceAddr,
		DevicePort:              uint16(*deviceLocalPort),
		DeviceInternalInterface: *deviceInternalInterface,

		AuditLogFile:          *auditLogFile,
		AuditLogMaxSize:       *auditLogMaxSize,
		AuditLogMaxBackups:    *auditLogMaxBackups,
		MetricsDataFile:       *metricsDataFile,
		MetricsDataMaxSize:    *metricsDataMaxSize,
		MetricsDataMaxBackups: *metricsDataMaxBackups,
	}
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	defer wg.Wait()
	srv, err := server.NewServer(opts)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	watcher.StartWatcher(cleanFunc, srv)
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Run()
	}()
	<-sigs
	watcher.StopWatcher()
	srv.Stop()

	if lumberjackLogger != nil {
		lumberjackLogger.Close()
	}
	watcher.StopProcess()
}

func cleanFunc(i any) {
	watcher.StopWatcher()
	r, ok := i.(*server.OMTServer)
	if ok {
		r.Stop()
	}
}
