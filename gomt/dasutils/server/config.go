package server

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"

	"gomt/core/proto/priv"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

//go:embed csv
var csvFS embed.FS

type CsvConfig struct {
	DumpDatas   []*DumpData
	DumpDataMap map[string]*DumpData
	DumpIDs     []*DumpID
	//CombineDatas []*CombinedData
	//CombinedDatas map[string]string
	CombinedDatas map[string]map[string]int

	//ParameterList map[string]*ParameterData
}

type DumpData struct {
	ModuleTypeId string  `csv:"Module Type ID"`
	DownSPer     int     `csv:"Downlink Frequency Start"`
	DownEPer     int     `csv:"Downlink Frequency End"`
	P1Slope      float32 `csv:"Port 1 Input Power Slope"`
	P1Intercept  float32 `csv:"Port 1 Input Power Nodal Increment"`
	P2Slope      float32 `csv:"Port 2 Input Power Slope"`
	P2Intercept  float32 `csv:"Port 2 Input Power Nodal Increment"`
	P3Slope      float32 `csv:"Port 3 Input Power Slope"`
	P3Intercept  float32 `csv:"Port 3 Input Power Nodal Increment"`
	P4Slope      float32 `csv:"Port 4 Input Power Slope"`
	P4Intercept  float32 `csv:"Port 4 Input Power Nodal Increment"`
}

type DumpID struct {
	DevTypeName string `csv:"Device Type Name"`
	Combiner    int    `csv:"Combiner"`
	SID         string `csv:"Serial Number"`
	DownSPer    string `csv:"Downlink Frequency Start"`
	DownEPer    string `csv:"Downlink Frequency End"`
	P1Slope     string `csv:"Port 1 Input Power Slope"`
	P1Intercept string `csv:"Port 1 Input Power Nodal Increment"`
	P2Slope     string `csv:"Port 2 Input Power Slope"`
	P2Intercept string `csv:"Port 2 Input Power Nodal Increment"`
	P3Slope     string `csv:"Port 3 Input Power Slope"`
	P3Intercept string `csv:"Port 3 Input Power Nodal Increment"`
	P4Slope     string `csv:"Port 4 Input Power Slope"`
	P4Intercept string `csv:"Port 4 Input Power Nodal Increment"`
}

type CombinedData struct {
	ModuleTypeId string `csv:"Module Type ID"`
	SerialNumber string `csv:"Serial Number"`
}

type ParameterData struct {
	DevTypeName  string
	Combiner     int
	ModuleTypeId string
	SID          priv.Parameter
	DownSPer     priv.Parameter
	DownEPer     priv.Parameter
	P1Slope      priv.Parameter
	P1Intercept  priv.Parameter
	P2Slope      priv.Parameter
	P2Intercept  priv.Parameter
	P3Slope      priv.Parameter
	P3Intercept  priv.Parameter
	P4Slope      priv.Parameter
	P4Intercept  priv.Parameter
}

func (s *CsvConfig) LoadFromFS(csvPath string) error {
	var rootFs fs.FS
	if len(csvPath) < 1 {
		if tmp, err := fs.Sub(csvFS, "csv"); err != nil {
			return errors.Errorf("fs sub csv error, %v", err)
		} else {
			rootFs = tmp
		}
	} else {
		rootFs = os.DirFS(csvPath)
	}

	fDup, err := rootFs.Open("dup_data.csv")
	if err != nil {
		return errors.Errorf("open dup_data.csv error, %v", err)
	}
	defer fDup.Close()

	d, err := io.ReadAll(fDup)
	if err != nil {
		return errors.Errorf("load dup_data.csv error, %v", err)
	}
	//fmt.Println(string(d))
	s.DumpDatas = []*DumpData{}
	s.DumpDataMap = map[string]*DumpData{}
	err = gocsv.UnmarshalBytes(d, &s.DumpDatas)
	if err != nil {
		return errors.Errorf("unmarshal dup_data.csv error, %v", err)
	}
	for _, v := range s.DumpDatas {
		s.DumpDataMap[fmt.Sprintf("%v:%v:%v", v.ModuleTypeId, v.DownSPer, v.DownEPer)] = v
	}

	fDupId, err := rootFs.Open("dup_id.csv")
	if err != nil {
		return errors.Errorf("open dup_id.csv error, %v", err)
	}
	defer fDupId.Close()
	d, err = io.ReadAll(fDupId)
	if err != nil {
		return errors.Errorf("load dup_id.csv error, %v", err)
	}
	s.DumpIDs = []*DumpID{}
	err = gocsv.UnmarshalBytes(d, &s.DumpIDs)
	if err != nil {
		return errors.Errorf("unmarshal dup_id.csv error, %v", err)
	}

	fComData, err := rootFs.Open("combined_data.csv")
	if err != nil {
		return errors.Errorf("open combined_data.csv error, %v", err)
	}
	defer fComData.Close()
	d, err = io.ReadAll(fComData)
	if err != nil {
		return errors.Errorf("load combined_data.csv error, %v", err)
	}
	//s.CombinedDatas = map[string]string{}
	s.CombinedDatas = map[string]map[string]int{}
	cdata := []*CombinedData{}
	err = gocsv.UnmarshalBytes(d, &cdata)
	if err != nil {
		return errors.Errorf("unmarshal combined_data.csv error, %v", err)
	}
	for _, v := range cdata {
		if len(v.SerialNumber) < 1 {
			continue
		}
		//if vv, ok := s.CombinedDatas[v.SerialNumber]; ok {
		//fmt.Printf("exist, %v-%v, %v\n", vv, v.SerialNumber, v.ModuleTypeId)
		//}

		//s.CombinedDatas[v.SerialNumber] = v.ModuleTypeId
		if _, ok := s.CombinedDatas[v.SerialNumber]; !ok {
			s.CombinedDatas[v.SerialNumber] = map[string]int{v.ModuleTypeId: 1}
		} else {
			if _, ok := s.CombinedDatas[v.SerialNumber][v.ModuleTypeId]; !ok {
				s.CombinedDatas[v.SerialNumber][v.ModuleTypeId] = 1
			} else {
				fmt.Printf("repeat combined data, %v - %v\n", v.SerialNumber, v.ModuleTypeId)
			}
		}
	}
	//fmt.Printf("load finish, %d, %d\n", len(s.CombinedDatas), i)
	//s.ParameterList = map[string]*ParameterData{}
	return nil
}
