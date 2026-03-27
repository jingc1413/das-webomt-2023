package utils

import (
	"github.com/go-ini/ini"
	"github.com/pkg/errors"
)

func IniReadValues(filename string, names []string) (map[string]map[string]string, error) {
	values := map[string]map[string]string{}

	f, err := ini.Load(filename)
	if err != nil {
		return values, errors.Wrap(err, "load file")
	}

	for _, name := range names {
		sectionValues := map[string]string{}
		if section := f.Section(name); section != nil {
			for _, key := range section.Keys() {
				sectionValues[key.Name()] = key.Value()
			}
		}
		values[name] = sectionValues
	}
	return values, nil
}

func IniUpdateValues(filename string, values map[string]map[string]string) error {
	f, err := ini.Load(filename)
	if err != nil {
		return errors.Wrap(err, "load file")
	}
	for sectionName, sectionValues := range values {
		for k, v := range sectionValues {
			section := f.Section(sectionName)
			if section == nil {
				section, err = f.NewSection(sectionName)
				if err != nil {
					return errors.Wrap(err, "new section")
				}
			}
			if key := section.Key(k); key != nil {
				key.SetValue(v)
			} else {
				key, err = section.NewKey(k, v)
				if err != nil {
					return errors.Wrap(err, "new key")
				}
			}
		}
	}
	return f.SaveTo(filename)
}
