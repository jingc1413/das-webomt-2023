package model

type ProductDefine struct {
	Schema          string   `json:"Schema"`
	ProductTypeName string   `json:"ProductTypeName"`
	DeviceTypeName  string   `json:"DeviceTypeName"`
	DeviceTypeID    int      `json:"DeviceTypeID"`
	MibModule       string   `json:"MibModule"`
	Versions        []string `json:"Versions"`
	// Charset         string `json:"Charset"`
	SupportLayout bool
}

func (m ProductDefine) GetMatcheVersion(ver string) string {
	for _, v := range m.Versions {
		if v == ver {
			return v
		}
	}
	return ""
}
func (m ProductDefine) MatcheVersion(ver string) bool {
	v := m.GetMatcheVersion(ver)
	return v == ver
}

var ProductDefines []*ProductDefine = []*ProductDefine{
	{Schema: "kddi", ProductTypeName: "AU", DeviceTypeName: "MU", DeviceTypeID: 200, MibModule: ""},
	{Schema: "kddi", ProductTypeName: "SAU", DeviceTypeName: "Slave MU", DeviceTypeID: 201, MibModule: ""},
	{Schema: "kddi", ProductTypeName: "EU", DeviceTypeName: "HU", DeviceTypeID: 202, MibModule: ""},
	{Schema: "kddi", ProductTypeName: "RU", DeviceTypeName: "RU", DeviceTypeID: 203, MibModule: ""},
	{Schema: "kddi", ProductTypeName: "AU", DeviceTypeName: "MUII", DeviceTypeID: 190, MibModule: ""},
	{Schema: "kddi", ProductTypeName: "EU", DeviceTypeName: "HUII", DeviceTypeID: 192, MibModule: ""},
	{Schema: "kddi", ProductTypeName: "RU", DeviceTypeName: "RUII", DeviceTypeID: 193, MibModule: ""},

	{Schema: "default", ProductTypeName: "AU", DeviceTypeName: "Primary AU", DeviceTypeID: 225, MibModule: "SUNWAVE-DAS-MAU-MIB"},
	{Schema: "default", ProductTypeName: "SAU", DeviceTypeName: "Secondary AU", DeviceTypeID: 224, MibModule: "SUNWAVE-DAS-SAU-MIB"},
	{Schema: "default", ProductTypeName: "EU", DeviceTypeName: "EU", DeviceTypeID: 242, MibModule: "SUNWAVE-DAS-EU-MIB"},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "LP-RU", DeviceTypeID: 228, MibModule: "SUNWAVE-DAS-LP-RU-MIB"},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "HPRU-2G", DeviceTypeID: 246, MibModule: "SUNWAVE-DAS-HP-2G-RU-MIB"},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "HP-RU", DeviceTypeID: 245, MibModule: "SUNWAVE-DAS-HP-RU-MIB"},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "HP-F-RU", DeviceTypeID: 232, MibModule: "SUNWAVE-DAS-HP-F-HU-MIB"},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "XP-RU", DeviceTypeID: 251, MibModule: "SUNWAVE-DAS-XP-RU-MIB"},

	{Schema: "default", ProductTypeName: "AU", DeviceTypeName: "Primary A2", DeviceTypeID: 230, MibModule: "SUNWAVE-DAS-MA2-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "AU", DeviceTypeName: "Master A2", DeviceTypeID: 230, MibModule: "SUNWAVE-DAS-MA2-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "SAU", DeviceTypeName: "Secondary A2", DeviceTypeID: 231, MibModule: "SUNWAVE-DAS-SA2-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "SAU", DeviceTypeName: "Slave A2", DeviceTypeID: 231, MibModule: "SUNWAVE-DAS-SA2-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "SAU", DeviceTypeName: "SA2", DeviceTypeID: 229, MibModule: "SUNWAVE-DAS-SA2-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "EU", DeviceTypeName: "EU-O", DeviceTypeID: 227, MibModule: "SUNWAVE-DAS-EU212-MIB", Versions: []string{"2.1"}},
	{Schema: "default", ProductTypeName: "EU", DeviceTypeName: "E2-O", DeviceTypeID: 210, MibModule: "SUNWAVE-DAS-e2O-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "N2-RU", DeviceTypeID: 235, MibModule: "SUNWAVE-DAS-N2RU-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "M2-RU", DeviceTypeID: 219, MibModule: "SUNWAVE-DAS-M2RU-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "H2-RU", DeviceTypeID: 211, MibModule: "SUNWAVE-DAS-H2RU-MIB", Versions: []string{"2.1.2"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "X2-RU", DeviceTypeID: 212, MibModule: "SUNWAVE-DAS-X2RU-MIB", Versions: []string{"2.1.3"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "HP-RU", DeviceTypeID: 245, MibModule: "SUNWAVE-DAS-HP-RU-MIB", Versions: []string{"2.2.1"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "HP-F-RU", DeviceTypeID: 232, MibModule: "SUNWAVE-DAS-HP-F-RU-MIB", Versions: []string{"2.2.1"}},

	{Schema: "corning", ProductTypeName: "AU", DeviceTypeName: "Primary A2", DeviceTypeID: 230, MibModule: "SUNWAVE-DAS-MA2-MIB", Versions: []string{"2.2"}},
	{Schema: "corning", ProductTypeName: "AU", DeviceTypeName: "Master A2", DeviceTypeID: 230, MibModule: "SUNWAVE-DAS-MA2-MIB", Versions: []string{"2.2"}},
	{Schema: "corning", ProductTypeName: "SAU", DeviceTypeName: "Secondary A2", DeviceTypeID: 231, MibModule: "SUNWAVE-DAS-SA2-MIB", Versions: []string{"2.2"}},
	{Schema: "corning", ProductTypeName: "SAU", DeviceTypeName: "Slave A2", DeviceTypeID: 231, MibModule: "SUNWAVE-DAS-SA2-MIB", Versions: []string{"2.2"}},
	{Schema: "corning", ProductTypeName: "SAU", DeviceTypeName: "SA2", DeviceTypeID: 229, MibModule: "SUNWAVE-DAS-SA2-MIB", Versions: []string{"2.2"}},
	{Schema: "corning", ProductTypeName: "EU", DeviceTypeName: "EU-O", DeviceTypeID: 227, MibModule: "SUNWAVE-DAS-EU212-MIB", Versions: []string{"2.0"}},
	{Schema: "corning", ProductTypeName: "EU", DeviceTypeName: "EU", DeviceTypeID: 227, MibModule: "SUNWAVE-DAS-EU212-MIB", Versions: []string{"2.0"}},
	{Schema: "corning", ProductTypeName: "EU", DeviceTypeName: "E2-O", DeviceTypeID: 210, MibModule: "SUNWAVE-DAS-e2O-MIB", Versions: []string{"2.0"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "N2-RU", DeviceTypeID: 235, MibModule: "SUNWAVE-DAS-N2RU-MIB", Versions: []string{"2.1"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "M2-RU", DeviceTypeID: 219, MibModule: "SUNWAVE-DAS-M2RU-MIB", Versions: []string{"2.2"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "H2-RU", DeviceTypeID: 211, MibModule: "SUNWAVE-DAS-H2RU-MIB", Versions: []string{"2.1"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "X2-RU", DeviceTypeID: 212, MibModule: "SUNWAVE-DAS-X2RU-MIB", Versions: []string{"0"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "HP-RU", DeviceTypeID: 245, MibModule: "SUNWAVE-DAS-HP-RU-MIB", Versions: []string{"2.0"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "HP-F-RU", DeviceTypeID: 232, MibModule: "SUNWAVE-DAS-HP-F-RU-MIB", Versions: []string{"2.0"}},

	{Schema: "bti", ProductTypeName: "AU", DeviceTypeName: "Master A2", DeviceTypeID: 230, MibModule: "BTI-DAS-MA2-MIB"},
	{Schema: "bti", ProductTypeName: "SAU", DeviceTypeName: "Slave A2", DeviceTypeID: 231, MibModule: "BTI-DAS-SA2-MIB"},
	{Schema: "bti", ProductTypeName: "SAU", DeviceTypeName: "SA2", DeviceTypeID: 229, MibModule: "BTI-DAS-SA2-MIB"},
	{Schema: "bti", ProductTypeName: "EU", DeviceTypeName: "EU-O", DeviceTypeID: 227, MibModule: "BTI-DAS-EU212-MIB"},
	{Schema: "bti", ProductTypeName: "EU", DeviceTypeName: "E2-O", DeviceTypeID: 210, MibModule: "BTI-DAS-e2O-MIB"},
	{Schema: "bti", ProductTypeName: "RU", DeviceTypeName: "N2-RU", DeviceTypeID: 235, MibModule: "BTI-DAS-N2RU-MIB"},
	{Schema: "bti", ProductTypeName: "RU", DeviceTypeName: "M2-RU", DeviceTypeID: 219, MibModule: "BTI-DAS-M2RU-MIB"},
	{Schema: "bti", ProductTypeName: "RU", DeviceTypeName: "HP-RU", DeviceTypeID: 245, MibModule: "BTI-DAS-HP-RU-MIB"},
	{Schema: "bti", ProductTypeName: "RU", DeviceTypeName: "H2-RU", DeviceTypeID: 211, MibModule: "BTI-DAS-H2RU-MIB"},
	{Schema: "bti", ProductTypeName: "RU", DeviceTypeName: "X2-RU", DeviceTypeID: 212, MibModule: "BTI-DAS-X2RU-MIB"},
	{Schema: "bti", ProductTypeName: "RU", DeviceTypeName: "HP-F-RU", DeviceTypeID: 232, MibModule: "BTI-DAS-HP-F-RU-MIB"},

	{Schema: "mavenir", ProductTypeName: "AU", DeviceTypeName: "Master A2", DeviceTypeID: 230, MibModule: "MAVENIR-DAS-MA2-MIB"},
	{Schema: "mavenir", ProductTypeName: "SAU", DeviceTypeName: "Slave A2", DeviceTypeID: 231, MibModule: "MAVENIR-DAS-SA2-MIB"},
	{Schema: "mavenir", ProductTypeName: "SAU", DeviceTypeName: "SA2", DeviceTypeID: 229, MibModule: "MAVENIR-DAS-SA2-MIB"},
	{Schema: "mavenir", ProductTypeName: "EU", DeviceTypeName: "EU-O", DeviceTypeID: 227, MibModule: "MAVENIR-DAS-EU212-MIB"},
	{Schema: "mavenir", ProductTypeName: "EU", DeviceTypeName: "E2-O", DeviceTypeID: 210, MibModule: "MAVENIR-DAS-E2O-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "N2-RU", DeviceTypeID: 235, MibModule: "MAVENIR-DAS-N2RU-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "M2-RU", DeviceTypeID: 219, MibModule: "MAVENIR-DAS-M2RU-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "HP-RU", DeviceTypeID: 245, MibModule: "MAVENIR-DAS-HP-RU-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "H2-RU", DeviceTypeID: 211, MibModule: "MAVENIR-DAS-H2RU-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "X2-RU", DeviceTypeID: 212, MibModule: "MAVENIR-DAS-X2RU-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "HP-F-RU", DeviceTypeID: 232, MibModule: "MAVENIR-DAS-HP-F-RU-MIB"},

	{Schema: "default", ProductTypeName: "AU", DeviceTypeName: "Primary A3", DeviceTypeID: 260, MibModule: "SUNWAVE-DAS-PA3-MIB",
		SupportLayout: true, Versions: []string{"1.2.4"}},
	{Schema: "default", ProductTypeName: "SAU", DeviceTypeName: "Secondary A3", DeviceTypeID: 261, MibModule: "SUNWAVE-DAS-SA3-MIB",
		SupportLayout: true, Versions: []string{"1.2.4"}},
	{Schema: "default", ProductTypeName: "EU", DeviceTypeName: "E3-O", DeviceTypeID: 255, MibModule: "SUNWAVE-DAS-E3O-MIB",
		SupportLayout: true, Versions: []string{"1.2.1"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "N3-RU", DeviceTypeID: 236, MibModule: "SUNWAVE-DAS-N3RU-MIB",
		SupportLayout: true, Versions: []string{"1.2.1"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "X3-RU", DeviceTypeID: 257, MibModule: "SUNWAVE-DAS-X3RU-MIB",
		SupportLayout: true, Versions: []string{"0"}},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "M3-RU-L", DeviceTypeID: 237, MibModule: "SUNWAVE-DAS-M3RU-L-MIB"},
	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "M3-RU-H", DeviceTypeID: 238, MibModule: "SUNWAVE-DAS-M3RU-H-MIB"},

	{Schema: "corning", ProductTypeName: "AU", DeviceTypeName: "Primary A3", DeviceTypeID: 260, MibModule: "CORNING-DDAS-EVERON6200-PA3-MIB",
		SupportLayout: true, Versions: []string{"demo", "0.12"}},
	{Schema: "corning", ProductTypeName: "SAU", DeviceTypeName: "Secondary A3", DeviceTypeID: 261, MibModule: "CORNING-DDAS-EVERON6200-SA3-MIB",
		SupportLayout: true, Versions: []string{"demo", "0.12"}},
	{Schema: "corning", ProductTypeName: "EU", DeviceTypeName: "E3-O", DeviceTypeID: 255, MibModule: "CORNING-DDAS-EVERON6200-E3O-MIB",
		SupportLayout: true, Versions: []string{"0.10"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "N3-RU", DeviceTypeID: 236, MibModule: "CORNING-DDAS-EVERON6200-N3RU-MIB",
		SupportLayout: true, Versions: []string{"0.11"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "M3-RU-L", DeviceTypeID: 237, MibModule: "CORNING-DDAS-EVERON6200-M3RU-L-MIB",
		SupportLayout: true, Versions: []string{"0.7"}},
	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "M3-RU-H", DeviceTypeID: 238, MibModule: "CORNING-DDAS-EVERON6200-M3RU-H-MIB",
		SupportLayout: true, Versions: []string{"0.5"}},

	{Schema: "mavenir", ProductTypeName: "AU", DeviceTypeName: "Primary A3", DeviceTypeID: 260, MibModule: "MAVENIR-DAS-PA3-MIB"},
	{Schema: "mavenir", ProductTypeName: "SAU", DeviceTypeName: "Secondary A3", DeviceTypeID: 261, MibModule: "MAVENIR-DAS-SA3-MIB"},
	{Schema: "mavenir", ProductTypeName: "EU", DeviceTypeName: "E3-O", DeviceTypeID: 255, MibModule: "MAVENIR-DAS-E3O-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "N3-RU", DeviceTypeID: 236, MibModule: "MAVENIR-DAS-N3RU-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "M3-RU-L", DeviceTypeID: 237, MibModule: "MAVENIR-DAS-M3RU-L-MIB"},
	{Schema: "mavenir", ProductTypeName: "RU", DeviceTypeName: "M3-RU-H", DeviceTypeID: 238, MibModule: "MAVENIR-DAS-M3RU-H-MIB"},

	{Schema: "bti", ProductTypeName: "AU", DeviceTypeName: "Primary A3", DeviceTypeID: 260, MibModule: "BTI-DAS-PA3-MIB"},
	{Schema: "bti", ProductTypeName: "SAU", DeviceTypeName: "Secondary A3", DeviceTypeID: 261, MibModule: "BTI-DAS-SA3-MIB"},
	{Schema: "bti", ProductTypeName: "EU", DeviceTypeName: "E3-O", DeviceTypeID: 255, MibModule: "BTI-DAS-E3O-MIB"},
	{Schema: "bti", ProductTypeName: "RU", DeviceTypeName: "N3-RU", DeviceTypeID: 236, MibModule: "BTI-DAS-N3RU-MIB"},
	{Schema: "bti", ProductTypeName: "RU", DeviceTypeName: "M3-RU-L", DeviceTypeID: 237, MibModule: "BTI-DAS-M3RU-L-MIB"},

	{Schema: "default", ProductTypeName: "RU", DeviceTypeName: "MP-PS", DeviceTypeID: 222, MibModule: "SUNWAVE-DAS-MPPS-MIB"},
	{Schema: "default", ProductTypeName: "SAU", DeviceTypeName: "Slave PSAU", DeviceTypeID: 217, MibModule: "SUNWAVE-DAS-SPSAU-MIB"},
	{Schema: "default", ProductTypeName: "AU", DeviceTypeName: "Master PSAU", DeviceTypeID: 216, MibModule: "SUNWAVE-DAS-MPSAU-MIB"},
	{Schema: "default", ProductTypeName: "SAU", DeviceTypeName: "Slave AUAIR", DeviceTypeID: 218, MibModule: "SUNWAVE-DAS-AUAIR-MIB"},

	{Schema: "corning", ProductTypeName: "RU", DeviceTypeName: "MP-PS", DeviceTypeID: 222, MibModule: "CORNING-DAS-MPPS-MIB"},
	{Schema: "corning", ProductTypeName: "SAU", DeviceTypeName: "Slave PSAU", DeviceTypeID: 217, MibModule: "CORNING-DAS-SPSAU-MIB"},
	{Schema: "corning", ProductTypeName: "AU", DeviceTypeName: "Master PSAU", DeviceTypeID: 216, MibModule: "CORNING-DAS-MPSAU-MIB"},
	{Schema: "corning", ProductTypeName: "SAU", DeviceTypeName: "Slave AUAIR", DeviceTypeID: 218, MibModule: "CORNING-DAS-AUAIR-MIB"},
}

func GetProductDefine(schema string, deviceTypeName string) *ProductDefine {
	for _, def := range ProductDefines {
		if def.Schema == schema && def.DeviceTypeName == deviceTypeName {
			return def
		}
	}
	return nil
}

func GetProductDefineByDeviceTypeID(schema string, deviceTypeID int) *ProductDefine {
	for _, def := range ProductDefines {
		if def.Schema == schema && def.DeviceTypeID == deviceTypeID {
			return def
		}
	}
	return nil
}
