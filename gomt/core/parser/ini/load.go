package ini

import (
	ini "github.com/ochinchina/go-ini"
)

func Load(inifiles []string) map[string]map[string]string {
	out := map[string]map[string]string{}
	for _, inifile := range inifiles {
		data := ini.Load(inifile)
		sections := data.Sections()
		for _, section := range sections {
			m := map[string]string{}
			keys := section.Keys()
			for _, key := range keys {
				m[key.ValueWithDefault("")] = key.Name()
			}
			out[section.Name] = m
		}
	}
	return out
}
