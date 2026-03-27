package model

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:embed asserts/codings.json
var codingsData []byte

type CodingDefineMap struct {
	FamilyMap           map[string]string
	ProductTypeMap      map[string]string
	PlatformMap         map[string]string
	RevisionMap         map[string]string
	DeviceTypeNameMap   map[string]string
	SnmpEnterpriseIdMap map[string]string
	SchemaMap           map[string]string
}

var defaultCodingDefineMap *CodingDefineMap
var defaultCodingDefineMapOnce sync.Once

func GetDefaultCodingDefineMap() *CodingDefineMap {
	defaultCodingDefineMapOnce.Do(func() {
		m, err := loadCodingDefineMap(codingsData)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "load product codings file"))
			return
		}
		defaultCodingDefineMap = m
	})

	return defaultCodingDefineMap
}

func loadCodingDefineMap(data []byte) (*CodingDefineMap, error) {
	m := CodingDefineMap{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errors.Wrap(err, "unmarshal content")
	}
	return &m, nil
}
