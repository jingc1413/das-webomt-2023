package model

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gomt/core/utils"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var defaultAllParameterDefinesMutex sync.Mutex
var defaultAllParameterDefinesSetup bool
var defaultAllParameterDefinesMap map[string]ParameterDefines
var defaultAllParameterDefinesList []ParameterDefines

func GetAllParameterDefinesMap() map[string]ParameterDefines {
	defaultAllParameterDefinesMutex.Lock()
	defer defaultAllParameterDefinesMutex.Unlock()
	if !defaultAllParameterDefinesSetup {
		logrus.Fatal("all parameter defines is not setup")
	}
	return defaultAllParameterDefinesMap
}

func GetParameterDefinesByDeviceTypeName(schema string, deviceTypeName string, version string) ParameterDefines {
	version = FormatVersion(version)

	productDef := GetProductDefine(schema, deviceTypeName)
	if productDef == nil {
		return nil
	}
	if match := productDef.MatcheVersion(version); !match {
		version = "latest"
	}

	all := GetAllParameterDefinesMap()
	if v, ok := all[deviceTypeName]; ok && v != nil {
		return v.GetParameterDefinesByVersion(schema, version)
	}
	return ParameterDefines{}
}

func SetupAllParameterDefineMap(fsys fs.FS, base string) error {
	defaultAllParameterDefinesMutex.Lock()
	defer defaultAllParameterDefinesMutex.Unlock()

	all, err := loadAllParameterDefineMap(fsys, base)
	if err != nil {
		return err
	}
	list := []ParameterDefines{}
	for _, v := range all {
		tmp := v
		list = append(list, tmp)
	}
	defaultAllParameterDefinesSetup = true
	defaultAllParameterDefinesMap = all
	defaultAllParameterDefinesList = list
	return nil
}

func loadAllParameterDefineMap(fsys fs.FS, base string) (map[string]ParameterDefines, error) {
	out := map[string]ParameterDefines{}

	entries, err := fs.ReadDir(fsys, base)
	if err != nil {
		return out, errors.Wrap(err, "read models dir")
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if deviceTypeName := entry.Name(); deviceTypeName != "" {

				filename := filepath.Join(base, deviceTypeName, "parameters.json")
				filename = filepath.ToSlash(filepath.Clean(filename))
				if utils.ExistsFile(filename) {
					f, err := fsys.Open(filename)
					if err != nil {
						return out, errors.Wrap(err, "open file "+deviceTypeName)
					}
					data, err := io.ReadAll(f)
					if err != nil {
						return out, errors.Wrap(err, "read file "+deviceTypeName)
					}
					defs := ParameterDefines{}
					if err := defs.LoadFromData(data); err != nil {
						return out, errors.Wrap(err, "load "+deviceTypeName)
					}
					switch deviceTypeName {
					case "Primary A2", "Master A2":
						out["Primary A2"] = defs
						out["Master A2"] = defs
					case "Secondary A2", "Slave A2", "SA2":
						out["Secondary A2"] = defs
						out["Slave A2"] = defs
						out["SA2"] = defs
					default:
						out[deviceTypeName] = defs
					}
				} else {
					filename := filepath.Join(base, deviceTypeName, "parameters.json.gz")
					filename = filepath.ToSlash(filepath.Clean(filename))
					f, err := fsys.Open(filename)
					if err != nil {
						return out, errors.Wrap(err, "read file "+deviceTypeName)
					}

					gzipReader, err := gzip.NewReader(f)
					if err != nil {
						return out, errors.Wrap(err, "decompress file "+deviceTypeName)
					}
					defer gzipReader.Close()
					data, err := io.ReadAll(gzipReader)
					if err != nil {
						return out, errors.Wrap(err, "decompress file "+deviceTypeName)
					}
					defs := ParameterDefines{}
					if err := defs.LoadFromData(data); err != nil {
						return out, errors.Wrap(err, "load "+deviceTypeName)
					}
					switch deviceTypeName {
					case "Primary A2", "Master A2":
						out["Primary A2"] = defs
						out["Master A2"] = defs
					case "Secondary A2", "Slave A2", "SA2":
						out["Secondary A2"] = defs
						out["Slave A2"] = defs
						out["SA2"] = defs
					default:
						out[deviceTypeName] = defs
					}
				}
			}
		}
	}
	for k, _ := range out {
		logrus.Tracef("load model parameters: %v", k)
	}
	return out, nil
}

type ParameterDefines []*ParameterDefine

func (m *ParameterDefines) LoadFromData(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return errors.Wrap(err, "unmarshal json data")
	}
	return nil
}

func (m *ParameterDefines) LoadFromFile(filename string) error {
	*m = ParameterDefines{}
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "read file")
	}
	return m.LoadFromData(data)
}

func (m ParameterDefines) GetParameterDefinesByVersion(prefix string, version string) ParameterDefines {
	out := ParameterDefines{}
	for _, def := range m {
		if def.Versions.MatchVersion(prefix, version) {
			out = append(out, def)
		}
	}
	return out
}

func (m ParameterDefines) GetParameterDefine(privOid PrivObjectId) *ParameterDefine {
	re := regexp.MustCompile(`(.*)\[.*\]`)
	if subs := re.FindStringSubmatch(string(privOid)); len(subs) == 2 {
		privOid = PrivObjectId(subs[1])
	}
	for _, def := range m {
		if def.PrivOid == privOid {
			return def
		}
	}
	return nil
}

func (m ParameterDefines) GetParameterDefineByName(name string) *ParameterDefine {
	for _, def := range m {
		if def.Name == name {
			return def
		}
	}
	return nil
}

func (m ParameterDefines) GetParameterDefinesByRegexpName(re *regexp.Regexp) []*ParameterDefine {
	out := []*ParameterDefine{}
	for _, def := range m {
		if re.MatchString(def.Name) {
			tmp := def
			out = append(out, tmp)
		}
	}
	return out
}

func (m *ParameterDefines) AddParameterDefine(def *ParameterDefine) error {
	re := regexp.MustCompile(`\s+`)

	for _, v := range *m {
		if v.PrivOid == def.PrivOid {
			if v.Name == def.Name && v.Access == def.Access && v.DataType == def.DataType {
				groups := v.Groups
				for _, group := range def.Groups {
					found := false
					for _, group2 := range groups {
						if group == group2 {
							found = true
							break
						}
					}
					if !found {
						groups = append(groups, group)
					}
				}
				v.Groups = []string{}
				for _, group := range groups {
					v.Groups = append(v.Groups, re.ReplaceAllString(group, " "))
				}
				return nil
			}
			return errors.Errorf("parameter already exists, oid=%v", def.PrivOid)
		}
	}
	*m = append(*m, def)
	return nil
}

func (m *ParameterDefines) SetParameterDefine(def ParameterDefine) {
	for i, v := range *m {
		if v.PrivOid == def.PrivOid {
			(*m)[i] = &def
			return
		}
	}
	*m = append(*m, &def)
}

type ParameterDefine struct {
	PrivOid        PrivObjectId           `json:"PrivOid,omitempty"`
	SnmpOid        SnmpObjectId           `json:"SnmpOid,omitempty"`
	Name           string                 `json:"Name"`
	Groups         []string               `json:"Groups"`
	Paths          []string               `json:"Paths"`
	VersionsString string                 `json:"Versions,omitempty"`
	Versions       ParameterVersionRanges `json:"-"`

	Description string `json:"Description,omitempty"`
	UnitName    string `json:"UnitName,omitempty"`
	Tips        string `json:"Tips,omitempty"`

	Access         Access             `json:"Access"`
	DataType       DataType           `json:"DataType"`
	InputType      string             `json:"InputType"`
	ByteSize       int64              `json:"ByteSize"`
	Ratio          *int64             `json:"Ratio,omitempty"`
	Max            *int64             `json:"Max,omitempty"`
	Min            *int64             `json:"Min,omitempty"`
	Step           *int64             `json:"Step,omitempty"`
	Options        Options            `json:"Options,omitempty"`
	MultipleOption bool               `json:"MultipleOption,omitempty"`
	Regexp         string             `json:"Regexp,omitempty"`
	Split          string             `json:"Split,omitempty"`
	Child          []*ParameterDefine `json:"Child,omitempty"`
}

type Access string

const (
	AccessReadWrite     Access = "rw"
	AccessReadOnly      Access = "ro"
	AccessWriteOnly     Access = "wo"
	AccessNotify        Access = "notify"
	AccessReportOnly    Access = "reportonly"
	AccessEventOnly     Access = "eventonly"
	AccessNotAccessible Access = "notaccessible"
)

type DataType string

const (
	DataTypeBinary  DataType = "binary"
	DataTypeBase64  DataType = "base64"
	DataTypeString  DataType = "string"
	DataTypeBool    DataType = "boolean"
	DataTypeInt     DataType = "int"
	DataTypeInt8    DataType = "int8"
	DataTypeInt16   DataType = "int16"
	DataTypeInt32   DataType = "int32"
	DataTypeInt64   DataType = "int64"
	DataTypeUInt    DataType = "uint"
	DataTypeUInt8   DataType = "uint8"
	DataTypeUInt16  DataType = "uint16"
	DataTypeUInt32  DataType = "uint32"
	DataTypeUInt64  DataType = "uint64"
	DataTypeFloat32 DataType = "float32"
	DataTypeFloat64 DataType = "float64"

	DataTypeStringArray DataType = "string:array"
	DataTypeInt32Array  DataType = "int32:array"
	DataTypeUint32Array DataType = "uint32:array"
	DataTypeObject      DataType = "object"

	DataTypeTimestamp DataType = "timestamp"
	DataTypeDateTime  DataType = "datetime"
	DataTypeIPV4      DataType = "ipv4"
	DataTypeIPV6      DataType = "ipv6"

	DataTypeTable DataType = "table"
	DataTypeRow   DataType = "row"
)

func (m DataType) IsNumeric() bool {
	if m == DataTypeInt ||
		m == DataTypeInt8 ||
		m == DataTypeInt16 ||
		m == DataTypeInt32 ||
		m == DataTypeInt64 ||
		m == DataTypeUInt ||
		m == DataTypeUInt8 ||
		m == DataTypeUInt16 ||
		m == DataTypeUInt32 ||
		m == DataTypeUInt64 ||
		m == DataTypeTimestamp {
		return true
	}
	return false
}

func (m ParameterDefine) Equal(v *ParameterDefine) bool {
	if v.PrivOid != m.PrivOid {
		return false
	}
	if v.Access != m.Access {
		return false
	}
	if v.DataType != m.DataType {
		return false
	}
	if v.ByteSize != m.ByteSize {
		return false
	}
	if v.Name != m.Name {
		return false
	}
	if v.UnitName != m.UnitName {
		return false
	}
	if (v.Ratio == nil && m.Ratio != nil) || (v.Ratio != nil && m.Ratio == nil) {
		return false
	}
	if v.Ratio != nil && m.Ratio != nil && *(v.Ratio) != *(m.Ratio) {
		return false
	}
	if (v.Step == nil && m.Step != nil) || (v.Step != nil && m.Step == nil) {
		return false
	}
	if v.Step != nil && m.Step != nil && *(v.Step) != *(m.Step) {
		return false
	}
	if (v.Max == nil && m.Max != nil) || (v.Max != nil && m.Max == nil) {
		return false
	}
	if v.Max != nil && m.Max != nil && *(v.Max) != *(m.Max) {
		return false
	}
	if (v.Min == nil && m.Min != nil) || (v.Min != nil && m.Min == nil) {
		return false
	}
	if v.Min != nil && m.Min != nil && *(v.Min) != *(m.Min) {
		return false
	}
	if (v.Options == nil && m.Options != nil) || (v.Options != nil && m.Options == nil) {
		return false
	}
	if v.Options != nil && m.Options != nil && v.Options.Equal(m.Options) == false {
		return false
	}
	if v.MultipleOption != m.MultipleOption {
		return false
	}
	return true
}

func (m *ParameterDefine) AddVersion(prefix string, version string) {
	if m.Versions == nil {
		m.Versions = ParameterVersionRanges{}
	}
	if err := m.Versions.AddVersion(prefix, version); err != nil {
		logrus.Fatalf(errors.Wrapf(err, "invalid version %v", version).Error())
	}
	m.VersionsString = m.Versions.String()
}

func (m *ParameterDefine) SetPath(path []string) {
	if path == nil || len(path) == 0 {
		return
	}
	if m.Paths == nil {
		m.Paths = []string{}
	}
	value := strings.Join(path, ".")
	value = strings.ReplaceAll(value, " ", "")
	for _, v := range m.Paths {
		if v == value {
			return
		}
	}
	m.Paths = append(m.Paths, value)
}

func (m *ParameterDefine) SetName(name string) {
	re := regexp.MustCompile(`\s+`)
	m.Name = strings.TrimSpace(re.ReplaceAllString(name, " "))
}

func (m *ParameterDefine) FixName(appendPrefix string, removePrefixList []string) {

	appendPrefix = strings.TrimSpace(appendPrefix)
	lastName := m.Name
	name := m.Name
	for _, v := range removePrefixList {
		if strings.HasPrefix(name, v) {
			name = strings.Replace(name, v, "", 1)
		}
	}
	if appendPrefix != "" && (!strings.HasPrefix(name, appendPrefix)) {
		name = strings.TrimSpace(name)
		name = appendPrefix + " " + name
	}
	if strings.Compare(name, lastName) != 0 {
		m.SetName(name)
	}
}

func (m *ParameterDefine) SetupArrayWithNumber(num int, split string) {
	m.DataType = DataType(fmt.Sprintf("%v:%v", DataTypeStringArray, num))
	m.Split = split
	m.Child = []*ParameterDefine{}
	for i := 0; i < num; i++ {
		child := *m
		child.PrivOid = PrivObjectId(fmt.Sprintf("%v[%v]", child.PrivOid, i))
		child.Child = nil
		child.Split = ""
		child.DataType = DataTypeString
		m.Child = append(m.Child, &child)
	}
	if num == 2 {
		rangeRe := regexp.MustCompile(`^(.*)\s+[\(]{0,1}([a-zA-Z0-9\s]+)\s*[-~/]{1}\s*([a-zA-Z0-9\s]+)[\)]{0,1}(.*)\s*$`)
		if subs := rangeRe.FindStringSubmatch(m.Name); len(subs) == 5 {
			if len(subs[4]) > 0 {
				m.Child[0].SetName(strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[2]) + " " + strings.TrimSpace(subs[4]))
				m.Child[1].SetName(strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[3]) + " " + strings.TrimSpace(subs[4]))
			} else {
				m.Child[0].SetName(strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[2]))
				m.Child[1].SetName(strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[3]))
			}
		}
	}
}
