package model

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ProductModel struct {
	ID              string
	Family          string //iDAS, PS, DRRU
	Type            string // A, R, E
	Platform        string // AU, A2, A3
	Revision        string // Common, NP, AU-AIR
	SFPType         string
	PowerLeadType   string
	ProductTypeName string // AU, RU, EU, SAU
	DeviceTypeName  string
}

func CreateProductModel(family string, typ string, platform string, revision string, productTypeName string, deviceTypeName string) (*ProductModel, error) {
	model := ProductModel{}
	root := GetDefaultProductMapRoot()

	var familyNode *ProductMapNode = nil
	k := family
	for _, v := range root.Children {
		if v.Code == k {
			familyNode = v
		}
	}
	if familyNode == nil {
		return nil, errors.Errorf("invalid product family '%v'", k)
	}

	k = fmt.Sprintf("%v-%v", family, typ)
	var typeNode *ProductMapNode = nil
	for _, v := range familyNode.Children {
		if v.Code == k {
			typeNode = v
		}
	}
	if typeNode == nil {
		return nil, errors.Errorf("invalid product type '%v'", k)
	}

	k = fmt.Sprintf("%v-%v%v", family, typ, platform)
	var platformNode *ProductMapNode = nil
	for _, v := range typeNode.Children {
		if v.Code == k {
			platformNode = v
		}
	}
	if platformNode == nil {
		return nil, errors.Errorf("invalid product platform '%v'", k)
	}

	k = fmt.Sprintf("%v-%v%v%v", family, typ, platform, revision)
	var revisonNode *ProductMapNode = nil
	for _, v := range platformNode.Children {
		if v.Code == k {
			revisonNode = v
		}
	}
	if revisonNode == nil {
		k = fmt.Sprintf("%v-%v%v%v", family, typ, platform, revision)
		var revisonNode *ProductMapNode = nil
		for _, v := range platformNode.Children {
			if v.Code == k {
				revisonNode = v
			}
		}
		if revisonNode == nil {
			return nil, errors.Errorf("invalid product revision '%v'", k)
		}
	}

	model.Type = typ
	model.ProductTypeName = productTypeName
	model.DeviceTypeName = deviceTypeName

	model.Family = family
	model.Platform = platform
	model.Revision = revision
	model.ID = model.String()

	return &model, nil
}

func (m ProductModel) String() string {
	return fmt.Sprintf("%v-%v%v%v:%v",
		m.Family, m.Type, m.Platform, m.Revision,
		m.ProductTypeName,
	)
}

func (m ProductModel) Equal(m2 ProductModel) bool {
	return m.String() == m2.String()
}

func (m ProductModel) Name() string {
	codings := GetDefaultCodingDefineMap()
	familyName := ""
	typeName := ""
	platformName := ""
	revisionName := ""

	k := m.Family
	familyName = codings.FamilyMap[k]

	k = fmt.Sprintf("%v-%v", m.Family, m.Type)
	typeName = codings.ProductTypeMap[k]

	k = fmt.Sprintf("%v-%v%v", m.Family, m.Type, m.Platform)
	platformName = codings.PlatformMap[k]

	k = fmt.Sprintf("%v-%v%v%v", m.Family, m.Type, m.Platform, m.Revision)
	if v, ok := codings.RevisionMap[k]; ok {
		revisionName = v
	} else {
		k = fmt.Sprintf("%v-%v%v%v", m.Family, m.Type, m.Platform, m.Revision)
		if v, ok := codings.RevisionMap[k]; ok {
			revisionName = v
		}
	}
	return fmt.Sprintf("%v_%v_%v_%v", familyName, typeName, platformName, revisionName)
}

func ParselModelName(name string, deviceTypeName string) ([]*ProductModel, error) {
	out := []*ProductModel{}

	args := strings.Split(name, ":")
	if len(args) != 2 {
		return out, errors.Errorf("invalid product model name '%v'", name)
	}
	name = args[0]
	productTypeName := args[1]

	re := regexp.MustCompile(`^(iDAS|DRRU|PS)-(A|R|E)(\d{1})(\d{2})`)
	if !re.MatchString(name) {
		return out, errors.Errorf("invalid product model name")
	}
	parts := re.FindStringSubmatch(name)
	if model, err := CreateProductModel(parts[1], parts[2], parts[3], parts[4], productTypeName, deviceTypeName); err == nil && model != nil {
		out = append(out, model)
	} else {
		return out, err
	}
	return out, nil
}

func (m ProductModel) MarshalModelName() string {
	codings := GetDefaultCodingDefineMap()
	return fmt.Sprintf("%v-%v%v%v(%v_%v_%v_%v)",
		m.Family, m.Type, m.Platform, m.Revision,
		codings.FamilyMap[m.Family],
		codings.ProductTypeMap[m.Type],
		codings.PlatformMap[m.Platform],
		codings.RevisionMap[m.Revision],
	)
}

type ProductModels []*ProductModel

func (m ProductModels) GetProductModelByModelID(modelID string) *ProductModel {
	modelID2 := strings.ReplaceAll(modelID, "_", "-")
	for _, v := range m {
		if v.ID == modelID || v.ID == modelID2 {
			model := v
			return model
		}
	}
	return nil
}

func (m ProductModels) GetProductModelByDeviceTypeName(deviceTypeName string) *ProductModel {
	for _, v := range m {
		if v.DeviceTypeName == deviceTypeName {
			tmp := *v
			return &tmp
		}
	}
	return nil
}

var defaultProductModels ProductModels
var defaultProductModelsOnce sync.Once

func GetDefaultProductModels() ProductModels {
	defaultProductModelsOnce.Do(func() {
		all := []*ProductModel{}
		keys := []string{}
		codings := GetDefaultCodingDefineMap()
		for k, v := range codings.DeviceTypeNameMap {
			names := strings.Split(v, "/")
			for _, name := range names {
				if objs, err := ParselModelName(k, name); err == nil {
					for _, obj := range objs {
						tmp := *obj
						all = append(all, &tmp)
						keys = append(keys, tmp.ID)
					}
				} else {
					logrus.Error(errors.Wrapf(err, "parse model name '%v'", k))
				}
			}
		}
		list := []*ProductModel{}
		sort.StringSlice(keys).Sort()
		for _, key := range keys {
			for _, obj := range all {
				if obj.ID == key {
					list = append(list, obj)
				}
			}
		}
		defaultProductModels = list
	})

	return defaultProductModels
}

func GetProductModelByDeviceTypeName(deviceTypeName string) *ProductModel {
	models := GetDefaultProductModels()
	return models.GetProductModelByDeviceTypeName(deviceTypeName)
}

func GetProductModelByModelID(modelID string) *ProductModel {
	models := GetDefaultProductModels()
	return models.GetProductModelByModelID(modelID)
}
