package model

import (
	"fmt"
	"gomt/core/utils"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

var (
	reVersionRevision = regexp.MustCompile(`^\d+(\.\d+)*_\d+(\.\d+)*$`)
	reVersion         = regexp.MustCompile(`^\d+(\.\d+)*$`)
	reSpecial         = regexp.MustCompile(`^(latest|dev|demo|test)$`)
)

const (
	VERSION_TYPE_INVALID = 0
	VERSION_TYPE_VER     = 1
	VERSION_TYPE_VERREV  = 2
	VERSION_TYPE_SPEC    = 3
)

type ParameterVersionRange struct {
	Prefix string
	Start  string
	End    string
	Type   int
}

func (m *ParameterVersionRange) UnmarshalString(raw string) error {
	args := strings.Split(raw, ",")
	argsLength := len(args)
	if argsLength != 1 && argsLength != 2 {
		return errors.Errorf("invalid parameter version range")
	}
	if argsLength > 1 {
		m.Prefix = args[0]
	} else {
		m.Prefix = ""
	}
	args2 := strings.Split(args[argsLength-1], "~")
	for i, v := range args2 {
		switch i {
		case 0:
			m.Start = strings.TrimSpace(v)
		case 1:
			m.End = strings.TrimSpace(v)
		default:
			return errors.Errorf("invalid parameter version range")
		}
	}
	if m.End == "" {
		m.End = m.Start
	}
	if m.Start == "" && m.End == "" {
		return errors.Errorf("invalid parameter version range")
	}

	startType := checkVersionType(m.Start)
	endType := checkVersionType(m.End)
	if startType <= 0 || endType <= 0 || startType != endType {
		return errors.Errorf("invalid parameter version range")
	}
	m.Type = startType

	if m.Start > m.End {
		return errors.Errorf("invalid parameter version range")
	}
	return nil
}

type ParameterVersionRanges []*ParameterVersionRange

func (m ParameterVersionRanges) Len() int { return len(m) }

func (m ParameterVersionRanges) Less(i, j int) bool {
	if m[i].Prefix == m[j].Prefix {
		return m[i].Start < m[j].Start
	}
	return m[i].Prefix < m[j].Prefix
}

func (m ParameterVersionRanges) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

func (m *ParameterVersionRanges) UnmarshalString(raw string) error {
	out := ParameterVersionRanges{}

	parts := strings.Split(raw, ";")
	for _, part := range parts {
		item := &ParameterVersionRange{}
		if err := item.UnmarshalString(part); err != nil {
			return errors.Wrap(err, part)
		}
		out = append(out, item)
	}
	*m = out
	sort.Sort(m)
	return nil
}

func (m ParameterVersionRanges) String() string {
	out := []string{}
	for _, item := range m {
		if item.Start == item.End {
			out = append(out, fmt.Sprintf("%v,%v", item.Prefix, item.Start))
		} else {
			out = append(out, fmt.Sprintf("%v,%v~%v", item.Prefix, item.Start, item.End))
		}
	}
	return strings.Join(out, ";")
}

func (m ParameterVersionRanges) MatchVersion(prefix string, version string) bool {
	if len(m) == 0 {
		return true
	}

	typ := checkVersionType(version)
	if typ <= 0 {
		return false
	}

	for _, item := range m {
		if item.Prefix != prefix {
			continue
		}
		if typ == VERSION_TYPE_SPEC {
			if item.Start == version || item.End == version {
				return true
			}
		} else {
			ltStart := compareVersions(typ, item.Start, version)
			gtEnd := compareVersions(typ, version, item.End)
			if ltStart <= 0 && gtEnd <= 0 {
				return true
			}
		}
	}
	return false
}

func (m *ParameterVersionRanges) AddVersion(prefix string, version string) error {
	if len(*m) > 0 && m.MatchVersion(prefix, version) {
		return nil
	}
	typ := checkVersionType(version)
	if typ <= 0 {
		return errors.New("invalid version format")
	}

	if typ == VERSION_TYPE_VERREV || typ == VERSION_TYPE_VER {
		for _, item := range *m {
			if item.Prefix != prefix || item.Type != typ {
				continue
			}
			ltStart := compareVersions(typ, item.Start, version)
			gtEnd := compareVersions(typ, version, item.End)
			if ltStart <= 0 && gtEnd <= 0 {
				return nil
			}
			if ltStart > 0 && gtEnd <= 0 {
				item.Start = version
				return nil
			}
			if ltStart <= 0 && gtEnd > 0 {
				item.End = version
				return nil
			}
		}
	}

	item := ParameterVersionRange{
		Prefix: prefix,
		Start:  version,
		End:    version,
		Type:   typ,
	}
	*m = append(*m, &item)
	sort.Sort(m)
	return nil
}

func checkVersionType(version string) int {
	if reVersion.MatchString(version) {
		return VERSION_TYPE_VER
	} else if reVersionRevision.MatchString(version) {
		return VERSION_TYPE_VERREV
	} else if reSpecial.MatchString(version) {
		return VERSION_TYPE_SPEC
	}
	return VERSION_TYPE_INVALID
}

func compareVersions(typ int, version, version2 string) int {
	if typ == VERSION_TYPE_SPEC {
		return strings.Compare(version, version2)
	} else if typ == VERSION_TYPE_VER {
		return utils.CompareVersion(version, version2)
	} else if typ == VERSION_TYPE_VERREV {
		args := strings.Split(version, "_")
		args2 := strings.Split(version, "_")
		if len(args) != 2 || len(args2) != 2 {
			return strings.Compare(version, version2)
		}
		return utils.CompareVersion(args[0], args[1])
	}
	return 0
}

func FormatVersion(version string) string {
	args := strings.Split(version, "-")
	return args[0]
}
