package xml

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type MachineRoot struct {
	XMLName      xml.Name `xml:"matchineroot"`
	NodeType     string   `xml:"nodetype,attr"`
	DeviceId     string   `xml:"deviceid,attr"`
	StationId    string   `xml:"stationid,attr"`
	HaveTopo     int64    `xml:"havetopo,attr"`
	NeedLogin    int64    `xml:"needlogin,attr"`
	IsMaster     int64    `xml:"ismaster,attr"`
	HaveAuto     int64    `xml:"haveauto,attr"`
	HaveSms      int64    `xml:"havesms,attr"`
	RepeaterName string   `xml:"repeatername,attr"`
	RepeaterType string   `xml:"repeatertype,attr"`
	Telephone    string   `xml:"telephone,attr"`
	Protocol     string   `xml:"protocol,attr"`
	AutoNum      string   `xml:"AutoNum,attr"`
	ParaLength   int64    `xml:"paralength,attr"`
	SupportEn    int64    `xml:"support_en,attr"`
	Modules      []Module `xml:"module"`
}

type Module struct {
	ID       string `xml:"id,attr"`
	NodeType string `xml:"nodetype,attr"`
	NameEn   string `xml:"name_en,attr"`
	Name     string `xml:"name,attr"`
	Address  string `xml:"address,attr"`
	CfgLevel string `xml:"cfglevel,attr"`
	Pages    []Page `xml:"page"`
}

func (m Module) GetPage(name string) *Page {
	for _, page := range m.Pages {
		if page.NameEn == name {
			tmp := page
			return &tmp
		}
	}
	return nil
}

type Page struct {
	ID        string  `xml:"id,attr"`
	NodeType  string  `xml:"nodetype,attr"`
	NameEn    string  `xml:"name_en,attr"`
	Name      string  `xml:"name,attr"`
	ChmCmd    string  `xml:"ChmCmd,attr"`
	Pachannel string  `xml:"Pachannel,attr"`
	UnCmd     string  `xml:"UnCmd,attr"`
	SwCmd     string  `xml:"SwCmd,attr"`
	SaveTime  string  `xml:"savetime,attr"`
	Groups    []Group `xml:"group"`
	Pages     []Group `xml:"page"`
}

func (m Page) CID() string {
	cid := m.ChmCmd
	if objectName := m.objectName(); objectName != "" {
		cid = cid + "-" + objectName
	}
	return strings.ToUpper(cid)
}

func (m Page) objectName() string {
	re := regexp.MustCompile(`(PA|POI|Combiner)([0-9]*)(\s*Info|)`)
	objectName := ""
	if subs := re.FindSubmatch([]byte(m.NameEn)); subs != nil {
		if m.Pachannel != "" && string(subs[2]) == "" {
			objectName = fmt.Sprintf("%v%v", string(subs[1]), m.Pachannel)
		} else {
			objectName = fmt.Sprintf("%v%v", string(subs[1]), string(subs[2]))
		}
	} else if m.NameEn == "FreqGainComp Test" {
		objectName = "FreqGainComp"
	} else if m.NameEn == "Debug Parameters" {
		objectName = "Debug"
	} else if m.NameEn == "MCU Parameter" {
		objectName = "MCU"
	}
	return objectName
}

type Group struct {
	ID       string  `xml:"id,attr"`
	NodeType string  `xml:"nodetype,attr"`
	NameEn   string  `xml:"name_en,attr"`
	Name     string  `xml:"name,attr"`
	GroupId  string  `xml:"groupid,attr"`
	Params   []Param `xml:"param"`
}

type Param struct {
	ID          string  `xml:"id,attr"`
	NodeType    string  `xml:"nodetype,attr"`
	DataType    string  `xml:"datatype,attr"`
	SectionEn   string  `xml:"section_en,attr"`
	Section     string  `xml:"section,attr"`
	EditType    string  `xml:"edittype,attr"`
	Rate        int64   `xml:"rate,attr"`
	Checked     int64   `xml:"checked,attr"`
	Len         int64   `xml:"len,attr"`
	Limit       string  `xml:"limit,attr"`
	CaptionEn   string  `xml:"caption_en,attr"`
	Caption     string  `xml:"caption,attr"`
	UnitEn      string  `xml:"unit_en,attr"`
	Unit        string  `xml:"unit,attr"`
	Range       string  `xml:"range,attr"`
	FaceValueEn string  `xml:"facevalue_en,attr"`
	FaceValue   string  `xml:"facevalue,attr"`
	Params      []Param `xml:"param"`
}
