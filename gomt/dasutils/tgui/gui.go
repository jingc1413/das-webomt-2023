package tgui

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"time"

	"gomt/core/model"
	"gomt/das/file"
	"gomt/dasutils/server"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
)

//go:embed dist
var embededFiles embed.FS

type gui struct {
	startButton     *widget.Button
	deviceAddrLabel *widget.Label
	deviceAddrEntry *widget.Entry
	statusLabel     *widget.Label
	statusEntry     *widget.Entry
	contentEntry    *widget.Entry

	//csvPathLabel *widget.Label
	//csvPathEntry *widget.Entry
	//csvSelectBtn *widget.Button
	pingButton *widget.Button

	//msg            *strings.Builder
	app          fyne.App
	charNums     int
	isExecFinish bool
	schema       string
	logDebug     bool
}

func (s *gui) isLockBtn() {
	go func() {
		isEnable := true
		for {
			if !s.isExecFinish && !isEnable {
				s.startButton.Enable()
				isEnable = true
			}
			if s.isExecFinish && isEnable {
				s.startButton.Disable()
				isEnable = false
			}
		}
	}()
}
func isValidIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil
}

//func (s *gui) setText(contentEntry string) {
//	ipAddr := s.contentEntry.Text
//	if len(ipAddr) < 1 || !isValidIP(ipAddr) {
//		s.contentEntry.SetText("")
//		s.contentEntry.Refresh()
//
//		s.contentEntry.SetText(fmt.Sprintf("Device IP Address error, %v\n", ipAddr))
//		//return
//	}
//	//s.contentEntry.Append(fmt.Sprintf("%v, %v\n", time.Now(), contentEntry))
//}

func Ping(addr string, t time.Duration) bool {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return false
	}
	pinger.Count = 3
	pinger.Timeout = time.Duration(t)

	pinger.SetPrivileged(true)
	pinger.Run()
	stats := pinger.Statistics()
	if stats.PacketsRecv >= 1 {
		return true
	} else {
		return false
	}
}

func (s *gui) initPingButton() {
	s.pingButton.OnTapped = func() {
		s.statusEntry.SetText("")
		s.pingButton.Disable()
		go func() {
			defer s.pingButton.Enable()
			ipAddr := s.deviceAddrEntry.Text
			if len(ipAddr) < 1 {
				s.contentEntry.SetText("")
				s.contentEntry.SetText("Device ip Address empty")
				return
			}
			if !isValidIP(ipAddr) {
				s.contentEntry.SetText("")
				msg := fmt.Sprintf("Device IP Address error, %v\n", ipAddr)
				s.contentEntry.SetText(msg)
				return
			} else {
				s.contentEntry.SetText("")
			}
			s.contentEntry.SetText(fmt.Sprintf("Start ping device %v\n", ipAddr))
			bRet := Ping(ipAddr, time.Second*5)
			if bRet {
				s.contentEntry.Append(fmt.Sprintf("Device IP Address %v is reachable\n", ipAddr))
			} else {
				s.contentEntry.Append(fmt.Sprintf("Device IP Address %v is not reachable\n", ipAddr))
			}
		}()
	}
}

func (s *gui) initStartButton() {
	s.startButton.OnTapped = func() {
		s.statusEntry.SetText("")
		s.contentEntry.SetText("")
		s.startButton.Disable()
		go func() {
			defer s.startButton.Enable()
			ipAddr := s.deviceAddrEntry.Text
			if len(ipAddr) < 1 {
				s.contentEntry.SetText("Device ip Address empty")
				return
			}
			if !isValidIP(ipAddr) {
				msg := fmt.Sprintf("Device IP Address error, %v\n", ipAddr)
				s.contentEntry.SetText(msg)
				return
			}
			bRet := Ping(ipAddr, time.Second*5)
			if !bRet {
				s.contentEntry.SetText("Device not online, please check\n")
				s.statusEntry.SetText("Failed")
				return
			}
			flag := s.appLogic(ipAddr)
			if !flag {
				s.statusEntry.SetText("Failed")
			}
		}()
	}
}

//var srv *priv.PrivMgmtServer

//func initPrivMgmt() bool {
//	schema := "default"
//	privPort := 9739
//	privWorker := 10
//	privateOpts := priv.PrivMgmtServerOptions{
//		Schema:     schema,
//		Addr:       fmt.Sprintf(":%v", privPort),
//		MaxWorkers: privWorker,
//	}
//	if err := priv.SetupDefaultPrivMgmtServer(privateOpts); err != nil {
//		//fmt.Errorf("setup private management server, %v", err)
//		//s.contentEntry.SetText(fmt.Sprintf("Setup private management server, %v", err))
//		return false
//	}
//
//	srv = priv.GetDefaultPrivMgmtServer()
//	go srv.Run()
//
//	return true
//}

func (s *gui) appLogic(deviceAddr string) bool {
	if s.logDebug {
		logrus.SetLevel(logrus.TraceLevel)
	}

	requestTimeout := 60
	timeoutDuration := time.Duration(requestTimeout) * time.Second
	var rootFs fs.FS
	isLocalDevice := deviceAddr == "127.0.0.1"
	if !isLocalDevice {
		if tmp, err := fs.Sub(embededFiles, "dist"); err != nil {
			fmt.Printf("fs sub error, %v", err)
			s.contentEntry.SetText(fmt.Sprintf("fs sub error, %v\n", err))
			return false
		} else {
			rootFs = tmp
		}
	}
	if err := model.SetupAllParameterDefineMap(rootFs, fmt.Sprintf("models/%v", s.schema)); err != nil {
		fmt.Printf("setup device model parametrs, %v\n", err)
		s.contentEntry.SetText(fmt.Sprintf("Setup device model parametrs, %v", err))
		return false
	}
	//allModels := model.GetAllParameterDefinesMap()
	//if err := selflayout.SetupAllLayoutMap(schema, allModels); err != nil {
	//	fmt.Printf("setup device model layouts, %v\n", err)
	//	s.contentEntry.SetText(fmt.Sprintf("Setup device model layouts, %v", err))
	//	return false
	//}

	csvPath, _ := os.Executable()
	csvPath = filepath.Dir(csvPath)
	cgiBinPath := "/drgfly/webomt/web/www/cgi-bin"
	configPath := "/drgfly/config"
	deviceTypeName := ""
	deviceLocalPort := 6001
	opts := server.ServerOptions{
		Schema:              s.schema,
		CgiRequestExpiresIn: timeoutDuration,
		CgiBinDir:           cgiBinPath,
		ConfigDir:           configPath,
		FileTypes:           []*file.FileType{},
		DeviceTypeName:      deviceTypeName,
		DeviceAddr:          deviceAddr,
		DevicePort:          uint16(deviceLocalPort),
		CsvDir:              csvPath,
	}
	srv, err := server.NewServer(opts)
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}
	srv.Run(s.contentEntry, s.statusEntry)
	//srv.Stop()
	return true
}

func (s *gui) show() {
	appWindow := s.app.NewWindow("DUP recalibration")
	appWindow.CenterOnScreen()
	appWindow.SetIcon(theme.SettingsIcon())

	s.initStartButton()
	s.initPingButton()
	deviceBorders := container.NewBorder(nil, nil, s.deviceAddrLabel, s.pingButton, s.deviceAddrEntry)

	statusForm := container.New(layout.NewFormLayout())
	statusForm.Add(s.statusLabel)
	statusForm.Add(s.statusEntry)

	vcontent := container.NewVBox(deviceBorders, s.startButton)
	border := container.NewBorder(vcontent, statusForm, nil, nil, s.contentEntry)

	appWindow.SetContent(border)
	appWindow.Resize(fyne.NewSize(600, 700))
	appWindow.SetFixedSize(true)
	appWindow.Show()
	s.app.Lifecycle().SetOnStopped(func() {
		//srv.Stop()
	})
	s.app.Run()
}

func InitGui() {
	//	initPrivMgmt()
	myApp := app.New()
	label := widget.NewMultiLineEntry()
	label.Wrapping = fyne.TextWrapBreak
	input := widget.NewEntry()

	deviceAddrLabel := widget.NewLabel("Device Address:")
	stsLabel := widget.NewLabel("Status:")
	stsEntry := widget.NewEntry()
	//stsEntry.Disable()

	//csvPathLabel := widget.NewLabel("Csv Path:")
	//csvEntry := widget.NewEntry()
	pingBtn := widget.NewButton("Ping Check", nil)
	//deviceAddrEntry := widget.NewMultiLineEntry()
	//deviceAddrEntry.SetPlaceHolder("Input Device Address")
	//deviceAddrEntry.Wrapping = fyne.TextWrapBreak //文字自动换行
	execBtn := widget.NewButton("Start", nil)
	ins := &gui{
		startButton:     execBtn,
		contentEntry:    label,
		deviceAddrEntry: input,
		statusLabel:     stsLabel,
		statusEntry:     stsEntry,
		//csvPathLabel:    csvPathLabel,
		//csvPathEntry:    csvEntry,
		pingButton:      pingBtn,
		deviceAddrLabel: deviceAddrLabel,
		//msg:     &strings.Builder{},
		app: myApp,
	}
	ins.isLockBtn()
	ins.show()
}
