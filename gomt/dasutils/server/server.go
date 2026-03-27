package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gomt/core/model"
	"gomt/core/proto/priv"
	"gomt/das"
	"gomt/das/agent"
	"gomt/das/file"
	"gomt/das/system"

	"fyne.io/fyne/v2/widget"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

const (
	SET_SUCCESS = 1
	SET_WARNING = 2
	SET_FAILED  = 3
)

type ServerOptions struct {
	Schema string

	CgiRequestExpiresIn time.Duration

	ConfigDir string
	CgiBinDir string
	FileTypes file.FileTypes

	CsvDir string

	DeviceTypeName          string
	DeviceAddr              string
	DeviceInternalInterface string
	DevicePort              uint16
}

type DasUtilsServer struct {
	opts          ServerOptions
	localhostMode bool

	dassys *system.DasSystem

	csvConfig CsvConfig

	queryDevicesLock    sync.RWMutex
	queryDevicesRunning bool
	queryDevicesTime    time.Time
	//queryDevicesProgress QueryDevicesProgress
	queryDevicesInterval time.Duration

	wg       sync.WaitGroup
	quit     chan bool
	shutdown bool
}

func NewServer(options ServerOptions) (*DasUtilsServer, error) {
	s := &DasUtilsServer{
		opts: options,
		quit: make(chan bool),
	}
	s.localhostMode = options.DeviceAddr == "127.0.0.1"
	s.dassys = system.NewDasSystem()
	return s, nil
}

func (s *DasUtilsServer) Stop() {
	s.shutdown = true
	s.quit <- true
}

func (s *DasUtilsServer) Run(showContent, stsEntry *widget.Entry) {
	time.Sleep(time.Second * 1)
	if err := s.setupLocalAgent(); err != nil {
		fmt.Printf("setup local device agent, %v\n", err)
		stsEntry.SetText("Failed")
		showContent.Append(fmt.Sprintf("setup local device agent, %v\n", err))
		//showContent.Append("setup local agent return\n")
		return
	}
	a := s.getDasDeviceAgent("local")
	if a != nil {
		devTypeName := a.GetDeviceTypeName()
		if devTypeName != "Primary A3" && devTypeName != "Primary A2" &&
			devTypeName != "Master A2" && devTypeName != "Master A3" {
			showContent.Append(fmt.Sprintf("Current device type name: %v, not support\n", devTypeName))
			fmt.Printf("Current device type name: %v, not support\n", devTypeName)
			stsEntry.SetText("Failed")
			return
		}
		if _, err := a.SetParameterValueOfJumpEnable(true); err != nil {
			fmt.Printf("set parameter of jump enable, %v\n", err)
		}
		defer a.SetParameterValueOfJumpEnable(false)
	} else {
		showContent.Append("get local device agent failed\n")
		stsEntry.SetText("Failed")
		return
	}

	s.updateDasDevices(true)

	stsEntry.SetText("")
	result := s.querySetParameter(showContent)
	if result == SET_SUCCESS {
		stsEntry.SetText("Success")
	} else if result == SET_FAILED {
		stsEntry.SetText("Failed")
	} else if result == SET_WARNING {
		stsEntry.SetText("Warning")
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (s *DasUtilsServer) recordFile(fileName, message string, showContext *widget.Entry) {
	{
		dirName := filepath.Dir(fileName)
		_, err := os.Stat(dirName)
		if os.IsNotExist(err) {
			err := os.Mkdir(dirName, 0755)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
			fmt.Println("Directory created successfully.")
		} else if err != nil {
			fmt.Println("Error checking directory:", err)
			return
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		showContext.Append(fmt.Sprintf("Open file failed, %v, %v", fileName, err))
		return
	}
	defer file.Close()
	currentTime := time.Now()
	timeStr := currentTime.Format("2006-01-02 15:04:05")
	_, err = fmt.Fprintf(file, "[%s] %v", timeStr, message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		showContext.Append(fmt.Sprintf("Writing to file: %v, %v\n", fileName, err))
		return
	}
}

// 1 success, 2 warning, 3 failed
func (s *DasUtilsServer) querySetParameter(showContext *widget.Entry) int {
	cmdp := filepath.Join(s.opts.CsvDir, "combined_data.csv")
	//dpd := filepath.Join(s.opts.CsvDir, "dup_data.csv")
	//dpi := filepath.Join(s.opts.CsvDir, "dup_id.csv")
	//if !fileExists(cmdp) && !fileExists(dpd) && !fileExists(dpi) {
	dupdata := filepath.Join(s.opts.CsvDir, "dup_data.csv")
	dupid := filepath.Join(s.opts.CsvDir, "dup_id.csv")
	if !fileExists(cmdp) && !fileExists(dupdata) && !fileExists(dupid) {
		s.opts.CsvDir = ""
	} else {
		showContext.Append(fmt.Sprintf("Use external csv path: %v\n", s.opts.CsvDir))
	}
	err := s.csvConfig.LoadFromFS(s.opts.CsvDir)
	if err != nil {
		showContext.Append(err.Error())
		return SET_FAILED
	}
	devInfos := s.dassys.GetAllDeviceInfos()
	if devInfos == nil || len(devInfos) < 1 {
		showContext.Append("GET: device infos empty\n")
		return SET_FAILED
	}
	for _, v := range devInfos {
		fmt.Printf("query devInfos: %v, %v\n", v.DeviceTypeName, v.SubID)
	}
	var failCount, succCount int
	for _, v := range s.csvConfig.DumpIDs {
		for _, dev := range devInfos {
			devInfo := dev
			if devInfo.DeviceTypeName != v.DevTypeName {
				continue
			}
			devAgent := s.dassys.GetAgent(fmt.Sprintf("%v", devInfo.SubID))
			if devAgent == nil {
				fmt.Printf("Query device agent empty, %v:%v\n", v.DevTypeName, devInfo.SubID)
				continue
			}
			{
				// check dup crc
				//Combiner1 00e2 894E
				//Combiner2 04e2
				//combiner3 08e2
				//combiner4 0ce2
				//record := "GET: Query Dup crc\n"
				//fmt.Println(record)
				//showContext.Append(record)
				var dupCrcs []priv.Parameter
				for retry := 0; retry < 2; retry++ {
					dupCrcs, err = devAgent.GetParameterValues([]priv.Parameter{
						priv.Parameter{ID: "TB2-COMBINER1.P00E2"},
						priv.Parameter{ID: "TB2-COMBINER2.P04E2"},
						priv.Parameter{ID: "TB2-COMBINER3.P08E2"},
						priv.Parameter{ID: "TB2-COMBINER4.P0CE2"},
					})
					if err != nil {
						record := fmt.Sprintf("GET: Query Dup crc error, %v, begin retry...\n", err)
						fmt.Println(record)
						showContext.Append(record)
						continue
					}
					break
				}

				var cp1Crc, cp2Crc, cp3Crc, cp4Crc string
				for _, v2 := range dupCrcs {
					if v2.ID == "TB2-COMBINER1.P00E2" {
						cp1Crc = cast.ToString(v2.Value)
					} else if v2.ID == "TB2-COMBINER2.P04E2" {
						cp2Crc = cast.ToString(v2.Value)
					} else if v2.ID == "TB2-COMBINER3.P08E2" {
						cp3Crc = cast.ToString(v2.Value)
					} else if v2.ID == "TB2-COMBINER4.P0CE2" {
						cp4Crc = cast.ToString(v2.Value)
					}
				}
				record := fmt.Sprintf("GET: Query Dup crc valus, %v\n", dupCrcs)
				fmt.Println(record)
				showContext.Append(record)
				record = fmt.Sprintf("GET: Query Dup crc valus, cp1: %v, cp2: %v, cp3: %v, cp4: %v\n", cp1Crc, cp2Crc, cp3Crc, cp4Crc)
				fmt.Printf(record)
				showContext.Append(record)
				if (cp1Crc != "0000" && cp1Crc != "894E") || (cp2Crc != "0000" && cp2Crc != "894E") ||
					(cp3Crc != "0000" && cp3Crc != "894E") || (cp4Crc != "0000" && cp4Crc != "894E") {
					record = fmt.Sprintf("dup crc check failed\n")
					fmt.Print(record)
					showContext.Append(record)
					continue
				}
				if cp1Crc == "0000" && cp2Crc == "0000" && cp3Crc == "0000" && cp4Crc == "0000" {
					record = fmt.Sprintf("dup crc check failed, no dup\n")
					fmt.Print(record)
					showContext.Append(record)
					continue
				}
			}
			fmt.Printf("\nGET: DeviceTypeName=%v, SubID=%v, Combiner=%v\n", v.DevTypeName, devInfo.SubID, v.Combiner)
			showContext.Append(fmt.Sprintf("\nGET: DeviceTypeName=%v, SubID=%v, Combiner=%v\n", v.DevTypeName, devInfo.SubID, v.Combiner))

			vals, err := devAgent.GetParameterValues([]priv.Parameter{
				priv.Parameter{ID: v.SID},
				priv.Parameter{ID: v.DownSPer},
				priv.Parameter{ID: v.DownEPer}})
			if err != nil || len(vals) < 3 {
				record := fmt.Sprintf("GET: Query parameter serial number, downsper, downeper error, %v, %v", err, vals)
				fmt.Println(record)
				showContext.Append(record)
				failCount++
				continue
			}
			var ssid, sCode, dfsCode, dfeCode string
			var downlinkFrequencyStart, downlinkFrequencyEnd int64
			for _, v2 := range vals {
				if v2.ID == v.SID {
					if sVal, ok := v2.Value.(string); ok {
						ssid = sVal
						ssid = strings.TrimRightFunc(ssid, func(r rune) bool {
							return r == '*'
						})
					}
					sCode = v2.Code
				} else if v2.ID == v.DownSPer {
					downlinkFrequencyStart = cast.ToInt64(v2.Value)
					dfsCode = v2.Code
				} else if v2.ID == v.DownEPer {
					downlinkFrequencyEnd = cast.ToInt64(v2.Value)
					dfeCode = v2.Code
				}
			}
			flag := false
			for _, v := range vals {
				if v.Code != "00" {
					flag = true
					break
				}
			}
			if flag {
				record := fmt.Sprintf("GET: query sid downsper downeper error code, %v\n", vals)
				fmt.Print(record)
				showContext.Append(record)
				continue
			}
			record := fmt.Sprintf("GET: SerialNumber=%v(%v), DownlinkFrequencyStart=%v(%v), DownlinkFrequencyEnd=%v(%v)\n",
				ssid, sCode, downlinkFrequencyStart, dfsCode, downlinkFrequencyEnd, dfeCode)
			fmt.Print(record)
			showContext.Append(record)

			defs := model.GetAllParameterDefinesMap()
			if defs == nil {
				record = "Module Parameter Defines empty\n"
				fmt.Print(record)
				//s.recordFile(fileName, record, showContext)
				showContext.Append(record)
				continue
			}
			def, ok := defs[devInfo.DeviceTypeName]
			if !ok || def == nil {
				record = fmt.Sprintf("Query Device Define empty, %v\n", devInfo.DeviceTypeName)
				fmt.Print(record)
				//s.recordFile(fileName, record, showContext)
				showContext.Append(record)
				continue
			}
			downSperDefine := def.GetParameterDefine(model.PrivObjectId(v.DownSPer))
			downEperDefine := def.GetParameterDefine(model.PrivObjectId(v.DownEPer))
			if downSperDefine == nil || downEperDefine == nil {
				record = fmt.Sprintf("Query Device Ratio empty, %v %v\n", v.DownSPer, v.DownEPer)
				fmt.Print(record)
				//s.recordFile(fileName, record, showContext)
				showContext.Append(record)
				continue
			}
			ddsper := decimal.NewFromInt(downlinkFrequencyStart)
			if downSperDefine.Ratio != nil && *downSperDefine.Ratio > 0 {
				ddsper = ddsper.Div(decimal.NewFromInt(*downSperDefine.Ratio))
			}
			ddeper := decimal.NewFromInt(downlinkFrequencyEnd)
			if downEperDefine.Ratio != nil && *downEperDefine.Ratio > 0 {
				ddeper = ddeper.Div(decimal.NewFromInt(*downEperDefine.Ratio))
			}
			downlinkFrequencyStart = ddsper.IntPart()
			downlinkFrequencyEnd = ddeper.IntPart()
			record = fmt.Sprintf("Conver end SerialNumber=%v, DownlinkFrequencyStart=%v, DownlinkFrequencyEnd=%v\n",
				ssid, downlinkFrequencyStart, downlinkFrequencyEnd)
			fmt.Print(record)
			showContext.Append(record)

			//moduleTypeId, ok := s.csvConfig.CombinedDatas[ssid]
			moduleTypeIds, ok := s.csvConfig.CombinedDatas[ssid]
			if !ok {
				record = fmt.Sprintf("Query ModuleTypeId from combined data csv file empty, %v=%v(%v)\n",
					v.SID, ssid, sCode)
				fmt.Print(record)
				showContext.Append(record)
				continue
			}
			if moduleTypeIds == nil || len(moduleTypeIds) < 1 {
				continue
			}
			var moduleTypeId string
			for mtId, _ := range moduleTypeIds {
				for _, dupData := range s.csvConfig.DumpDatas {
					if mtId == dupData.ModuleTypeId && int64(dupData.DownSPer) == downlinkFrequencyStart &&
						int64(dupData.DownEPer) == downlinkFrequencyEnd {
						moduleTypeId = mtId
						fmt.Printf("match combined and dupids, %v, [%v - %v - %v]\n", moduleTypeIds, mtId, dupData.DownSPer, dupData.DownEPer)
						break
					}
				}
				if len(moduleTypeId) > 0 {
					break
				}
			}
			if len(moduleTypeId) < 1 {
				record = fmt.Sprintf("find module type id from dupids and combined empty, %v, %v\n", ssid, moduleTypeIds)
				fmt.Print(record)
				showContext.Append(record)
				continue
			}

			serialNumber := devInfo.ElementSerialNumber
			if len(serialNumber) < 1 {
				serialNumber = "null"
			}
			recordContent := fmt.Sprintf("GET: DeviceTypeName=%v, ModuleTypeId=%v\n", v.DevTypeName, moduleTypeId)
			fmt.Println(recordContent)
			fileName := fmt.Sprintf("record/%v_%v_%v_%v.txt", v.DevTypeName, moduleTypeId, serialNumber, time.Now().Format("20060102150405"))
			s.recordFile(fileName, recordContent, showContext)
			showContext.Append(recordContent)

			/*defs := model.GetAllParameterDefinesMap()
			if defs == nil {
				record = "Module Parameter Defines empty\n"
				fmt.Print(record)
				s.recordFile(fileName, record, showContext)
				showContext.Append(record)
				continue
			}
			def, ok := defs[devInfo.DeviceTypeName]
			if !ok || def == nil {
				record = fmt.Sprintf("Query Device Define empty, %v\n", devInfo.DeviceTypeName)
				fmt.Print(record)
				s.recordFile(fileName, record, showContext)
				showContext.Append(record)
				continue
			}
			downSperDefine := def.GetParameterDefine(model.PrivObjectId(v.DownSPer))
			downEperDefine := def.GetParameterDefine(model.PrivObjectId(v.DownEPer))
			if downSperDefine == nil || downEperDefine == nil {
				record = fmt.Sprintf("Query Device Ratio empty, %v %v\n", v.DownSPer, v.DownEPer)
				fmt.Print(record)
				s.recordFile(fileName, record, showContext)
				showContext.Append(record)
				continue
			}
			ddsper := decimal.NewFromInt(downlinkFrequencyStart)
			if downSperDefine.Ratio != nil && *downSperDefine.Ratio > 0 {
				ddsper = ddsper.Div(decimal.NewFromInt(*downSperDefine.Ratio))
			}
			ddeper := decimal.NewFromInt(downlinkFrequencyEnd)
			if downEperDefine.Ratio != nil && *downEperDefine.Ratio > 0 {
				ddeper = ddeper.Div(decimal.NewFromInt(*downEperDefine.Ratio))
			}*/

			setData, ok := s.csvConfig.DumpDataMap[fmt.Sprintf("%v:%v:%v", moduleTypeId, ddsper.String(), ddeper.String())]
			if !ok {
				record = fmt.Sprintf("Query ModuleTypeId DownlinkFrequencyStart DownlinkFrequencyEnd from dup_data.csv empty, %v %v %v %v %v\n", moduleTypeId,
					downlinkFrequencyStart, downlinkFrequencyEnd, ddsper.String(), ddeper.String())
				fmt.Print(record)
				s.recordFile(fileName, record, showContext)
				showContext.Append(record)
				continue
			}
			ratio := decimal.NewFromInt(1000)
			parameterData := &ParameterData{
				DevTypeName:  v.DevTypeName,
				Combiner:     v.Combiner,
				ModuleTypeId: moduleTypeId,
				SID:          vals[0],
				DownSPer:     vals[1],
				DownEPer:     vals[2],
				P1Slope: priv.Parameter{
					ID:    v.P1Slope,
					Value: decimal.NewFromFloat32(setData.P1Slope).Mul(ratio),
				},
				P1Intercept: priv.Parameter{
					ID:    v.P1Intercept,
					Value: decimal.NewFromFloat32(setData.P1Intercept).Mul(ratio),
				},
				P2Slope: priv.Parameter{
					ID:    v.P2Slope,
					Value: decimal.NewFromFloat32(setData.P2Slope).Mul(ratio),
				},
				P2Intercept: priv.Parameter{
					ID:    v.P2Intercept,
					Value: decimal.NewFromFloat32(setData.P2Intercept).Mul(ratio),
				},
				P3Slope: priv.Parameter{
					ID:    v.P3Slope,
					Value: decimal.NewFromFloat32(setData.P3Slope).Mul(ratio),
				},
				P3Intercept: priv.Parameter{
					ID:    v.P3Intercept,
					Value: decimal.NewFromFloat32(setData.P3Intercept).Mul(ratio),
				},
				P4Slope: priv.Parameter{
					ID:    v.P4Slope,
					Value: decimal.NewFromFloat32(setData.P4Slope).Mul(ratio),
				},
				P4Intercept: priv.Parameter{
					ID:    v.P4Intercept,
					Value: decimal.NewFromFloat32(setData.P4Intercept).Mul(ratio),
				},
			}
			record = fmt.Sprintf("SET: DeviceTypeName=%v, SubID=%v, Combiner=%v\n", v.DevTypeName, devInfo.SubID, v.Combiner)
			fmt.Print(record)
			s.recordFile(fileName, record, showContext)
			showContext.Append(record)
			record = fmt.Sprintf("SET: ModuleTypeId=%v, "+
				"P1Slope=%v, P1Intercept=%v,"+
				"P2Slope=%v, P2Intercept=%v,"+
				"P3Slope=%v, P3Intercept=%v,"+
				"P4Slope=%v, P4Intercept=%v\n", moduleTypeId,
				parameterData.P1Slope.Value, parameterData.P1Intercept.Value,
				parameterData.P2Slope.Value, parameterData.P2Intercept.Value,
				parameterData.P3Slope.Value, parameterData.P3Intercept.Value,
				parameterData.P4Slope.Value, parameterData.P4Intercept.Value)
			fmt.Print(record)
			s.recordFile(fileName, record, showContext)
			showContext.Append(record)

			_, err = devAgent.SetParameterValues([]priv.Parameter{
				parameterData.P1Slope, parameterData.P1Intercept, parameterData.P2Slope, parameterData.P2Intercept,
				parameterData.P3Slope, parameterData.P3Intercept, parameterData.P4Slope, parameterData.P4Intercept})
			if err != nil {
				record = fmt.Sprintf("Set paremeter failed, %v\n", err)
				fmt.Print(record)
				s.recordFile(fileName, record, showContext)
				showContext.Append(record)
				failCount++
				continue
			} else {
				time.Sleep(time.Millisecond * 500)
				vals, err = devAgent.GetParameterValues([]priv.Parameter{
					priv.Parameter{ID: parameterData.P1Slope.ID},
					priv.Parameter{ID: parameterData.P1Intercept.ID},
					priv.Parameter{ID: parameterData.P2Slope.ID},
					priv.Parameter{ID: parameterData.P2Intercept.ID},
					priv.Parameter{ID: parameterData.P3Slope.ID},
					priv.Parameter{ID: parameterData.P3Intercept.ID},
					priv.Parameter{ID: parameterData.P4Slope.ID},
					priv.Parameter{ID: parameterData.P4Intercept.ID},
				})
				if err != nil || len(vals) < 8 {
					record = fmt.Sprintf("SET: Query parameter set vals, %v, %v", err, vals)
					fmt.Println(record)
					s.recordFile(fileName, record, showContext)
					showContext.Append(record)
					failCount++
					continue
				}
				record = fmt.Sprintf("SET: query set values, %v\n", vals)
				fmt.Print(record)
				showContext.Append(record)
				var p1slop, p1intercept, p2slope, p2intercept, p3slope, p3intercept, p4slope, p4intercept int64
				for _, v2 := range vals {
					if v2.ID == parameterData.P1Slope.ID {
						p1slop = cast.ToInt64(v2.Value)
					} else if v2.ID == parameterData.P1Intercept.ID {
						p1intercept = cast.ToInt64(v2.Value)
					} else if v2.ID == parameterData.P2Slope.ID {
						p2slope = cast.ToInt64(v2.Value)
					} else if v2.ID == parameterData.P2Intercept.ID {
						p2intercept = cast.ToInt64(v2.Value)
					} else if v2.ID == parameterData.P3Slope.ID {
						p3slope = cast.ToInt64(v2.Value)
					} else if v2.ID == parameterData.P3Intercept.ID {
						p3intercept = cast.ToInt64(v2.Value)
					} else if v2.ID == parameterData.P4Slope.ID {
						p4slope = cast.ToInt64(v2.Value)
					} else if v2.ID == parameterData.P4Intercept.ID {
						p4intercept = cast.ToInt64(v2.Value)
					}
				}
				record = fmt.Sprintf("SET: query parameter finish, "+
					"P1Slope=%v, P1Intercept=%v,"+
					"P2Slope=%v, P2Intercept=%v,"+
					"P3Slope=%v, P3Intercept=%v,"+
					"P4Slope=%v, P4Intercept=%v\n", p1slop, p1intercept,
					p2slope, p2intercept, p3slope,
					p3intercept, p4slope, p4intercept)
				fmt.Print(record)
				s.recordFile(fileName, record, showContext)
				showContext.Append(record)

				oldP1Slope := parameterData.P1Slope.Value.(decimal.Decimal)
				newP1Slope := decimal.NewFromInt(p1slop)
				oldP1Intercept := parameterData.P1Intercept.Value.(decimal.Decimal)
				newP1Intercept := decimal.NewFromInt(p1intercept)

				oldP2slope := parameterData.P2Slope.Value.(decimal.Decimal)
				newP2slope := decimal.NewFromInt(p2slope)
				oldP2Intercept := parameterData.P2Intercept.Value.(decimal.Decimal)
				newP2Intercept := decimal.NewFromInt(p2intercept)

				oldP3slope := parameterData.P3Slope.Value.(decimal.Decimal)
				newP3slope := decimal.NewFromInt(p3slope)
				oldP3Intercept := parameterData.P3Intercept.Value.(decimal.Decimal)
				newP3Intercept := decimal.NewFromInt(p3intercept)
				oldP4slope := parameterData.P4Slope.Value.(decimal.Decimal)
				newP4slope := decimal.NewFromInt(p4slope)
				oldP4Intercept := parameterData.P4Intercept.Value.(decimal.Decimal)
				newP4Intercept := decimal.NewFromInt(p4intercept)

				if oldP1Slope.Equal(newP1Slope) && oldP1Intercept.Equal(newP1Intercept) &&
					oldP2slope.Equal(newP2slope) && oldP2Intercept.Equal(newP2Intercept) &&
					oldP3slope.Equal(newP3slope) && oldP3Intercept.Equal(newP3Intercept) &&
					oldP4slope.Equal(newP4slope) && oldP4Intercept.Equal(newP4Intercept) {
					record = fmt.Sprintf("SET: query/set match success, DeviceTypeName=%v, SubID=%v, Combiner=%v\n", v.DevTypeName, devInfo.SubID, v.Combiner)
					fmt.Println(record)
					s.recordFile(fileName, record, showContext)
					showContext.Append(record)
					succCount++
				} else {
					failCount++
					record = fmt.Sprintf("SET: set parameter failed, DeviceTypeName=%v, SubID=%v, Combiner=%v, set value: %v\nquery value: %v\n",
						v.DevTypeName, devInfo.SubID, v.Combiner, *setData, vals)
					fmt.Println(record)
					s.recordFile(fileName, record, showContext)
					showContext.Append(record)
				}
			}
		}
	}
	if succCount < 1 {
		fmt.Printf("succ, fail: %v %v\n", succCount, failCount)
		return SET_FAILED
	}
	if failCount < 1 && succCount > 0 {
		fmt.Printf("2succ, fail: %v %v\n", succCount, failCount)
		return SET_SUCCESS
	}
	if failCount > 0 && succCount > 0 {
		fmt.Printf("3succ, fail: %v %v\n", succCount, failCount)
		return SET_WARNING
	}
	fmt.Printf("4succ, fail: %v %v\n", succCount, failCount)
	return SET_FAILED
}

func (s *DasUtilsServer) getDeviceInfo(devTypeName string) *model.DeviceInfo {
	for _, v := range s.dassys.GetAllDeviceInfos() {
		if v.DeviceTypeName == devTypeName {
			return v
		}
	}
	return nil
}

func (s *DasUtilsServer) setupLocalAgent() error {
	s.dassys.ClearAllAgent()
	deviceSub := "local"
	opts := agent.DasDeviceAgentOptions{
		DeviceAddr:        s.opts.DeviceAddr,
		PrivServerPort:    s.opts.DevicePort,
		HttpServerAddr:    das.MakeDeviceUrlBase(s.opts.DeviceAddr, 443, true),
		CGIBinUrlPath:     "/cgi-bin",
		CGIBinFilePath:    s.opts.CgiBinDir,
		FileTypes:         s.opts.FileTypes,
		ConfigPath:        s.opts.ConfigDir,
		InternalInterface: s.opts.DeviceInternalInterface,
	}
	a, err := s.dassys.SetupDasDeviceAgent(s.opts.Schema, deviceSub, nil, opts)
	if err != nil {
		return err
	}
	//info := a.GetDeviceInfo()
	//log.Printf("deviceType=%v subID=%v routeAddr=%v ip=%v\n",
	//	info.DeviceTypeName, info.SubID, info.RouteAddressString, info.IpAddressString)

	for {
		//log.Printf("waiting for local service available...\n")
		if ok := a.IsServiceAvailable(true); ok {
			break
		}
		time.Sleep(time.Second)
	}
	if err := a.UpdateDeviceInfo(); err != nil {
		//log.Printf("update local device info, %v\n", err)
	}
	return nil
}
