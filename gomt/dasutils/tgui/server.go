package tgui

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"gomt/core/proto/priv"
)

type tguiServer struct {
	appGui      *gui
	privMgmtSrv *priv.PrivMgmtServer
	schema      string

	wg sync.WaitGroup
}

// go build -ldflags="-s -w -H=windowsgui" -o dup_recalibration.exe
// go build -ldflags -H=windowsgui -o dup_recalibration_cn.exe app_gui.go
//
//		-ldflags="-s -w"  -ldflags="-H=windowsgui" -o dup_recalibration.exe
//	 -ldflags="-s -w"  -ldflags="-H=windowsgui" -o dup_recalibration.exe

// https://github.com/golang/go/issues/66998

// GODEBUG=tlsrsakex=1
func NewGuiServer() (*tguiServer, error) {
	// -ldflags -H=windowsgui
	guiSrv := &tguiServer{}
	guiSrv.schema = "corning"
	//guiSrv.schema = "default"
	logDebug := false
	myApp := app.New()
	label := widget.NewMultiLineEntry()
	label.Wrapping = fyne.TextWrapBreak
	input := widget.NewEntry()

	deviceAddrLabel := widget.NewLabel("Device Address:")
	stsLabel := widget.NewLabel("Status:")
	stsEntry := widget.NewEntry()

	pingBtn := widget.NewButton("Ping Check", nil)
	execBtn := widget.NewButton("Start", nil)
	ins := &gui{
		startButton:     execBtn,
		contentEntry:    label,
		deviceAddrEntry: input,
		statusLabel:     stsLabel,
		statusEntry:     stsEntry,
		pingButton:      pingBtn,
		deviceAddrLabel: deviceAddrLabel,
		app:             myApp,
		schema:          guiSrv.schema,
		logDebug:        logDebug,
	}
	guiSrv.appGui = ins

	{
		//privPort := 19739
		privWorker := 10
		privateOpts := priv.PrivMgmtServerOptions{
			Schema: guiSrv.schema,
			Addr:   ":",
			//Addr:       fmt.Sprintf(":%v", privPort),
			MaxWorkers: privWorker,
		}
		if err := priv.SetupDefaultPrivMgmtServer(privateOpts); err != nil {
			fmt.Printf("tgui setup private management server, %v", err)
			return nil, err
		}
		guiSrv.privMgmtSrv = priv.GetDefaultPrivMgmtServer()
	}
	return guiSrv, nil
}

func (s *tguiServer) Run() {
	defer s.wg.Wait()
	defer s.privMgmtSrv.Stop()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.privMgmtSrv.Run()
	}()

	s.appGui.show()
}
