package agent

import (
	"bufio"
	"gomt/core/metric"
	"gomt/core/model"
	"gomt/core/proto/priv"
	"gomt/core/utils"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func (s *DasDeviceAgent) GetMetricsData() *metric.MetricsData {
	return s.metrics
}

func (s *DasDeviceAgent) setupMetrics(filename string, maxsize, maxbackups int) {
	s.metrics = metric.NewMetricsData(filename, maxsize, maxbackups)

	s.metrics.AddItem("cpu_usage", "percent`", 2)
	s.metrics.AddItem("memory_usage", "percent", 2)
	s.metrics.AddItem("disk_usage", "percent", 2)
	s.metrics.AddItem("disk_total", "int64", 0)
	s.metrics.AddItem("disk_used", "int64", 0)
	s.metrics.AddItem("disk_free", "int64", 0)

	if def := s.defs.GetParameterDefine("T02.P0501"); def != nil {
		s.metrics.AddParamItem(def, 2)
	}

	if re := regexp.MustCompile(`(Radio|Input) Module \d+ Port \d+ Input Power`); re != nil {
		if defs := s.defs.GetParameterDefinesByRegexpName(re); defs != nil {
			for _, def := range defs {
				s.metrics.AddParamItem(def, 3)
			}
		}
	}

	if re := regexp.MustCompile(`(Radio|Amplifier) Module \d+ UL Input Power`); re != nil {
		if defs := s.defs.GetParameterDefinesByRegexpName(re); defs != nil {
			for _, def := range defs {
				s.metrics.AddParamItem(def, 3)
			}
		}
	}

	if re := regexp.MustCompile(`(Radio|Amplifier) Module \d+ DL Output Power`); re != nil {
		if defs := s.defs.GetParameterDefinesByRegexpName(re); defs != nil {
			for _, def := range defs {
				s.metrics.AddParamItem(def, 3)
			}
		}
	}
}

func (s *DasDeviceAgent) updateMetrics() {
	ts := time.Now().Unix()

	if ts%30 != 0 {
		return
	}
	t := time.Unix(ts, 0)
	if stats, err := utils.GetSystemStatus(1); err == nil {
		s.metrics.ObserveFloat64("cpu_usage", stats.CPU.Percent)
		s.metrics.ObserveFloat64("memory_usage", stats.Memory.UsedPercent)
	} else {
		s.log.Error(errors.Wrap(err, "get system stats"))
	}
	if stats, err := utils.GetDirectoryDiskStatus([]string{"/drgfly"}); err == nil {
		for _, point := range stats.MountPoints {
			if point.MountPoint == "/drgfly" {
				s.metrics.SetUint64("disk_used", point.Used)
				s.metrics.SetUint64("disk_total", point.Total)
				s.metrics.SetUint64("disk_free", point.Free)
				s.metrics.SetFloat64("disk_usage", point.UsedPercent)
			}
		}
	} else {
		s.log.Error(errors.Wrap(err, "get disk stats"))
	}

	params := map[string]*model.ParameterDefine{}
	inputValues := []priv.Parameter{}
	for _, def := range s.metrics.Params {
		inputValues = append(inputValues, priv.Parameter{ID: string(def.PrivOid)})
		params[string(def.PrivOid)] = def
	}

	if output, err := s.GetParameterValues(inputValues); err == nil {
		for _, v := range output {
			// if v.Code == "00" {
			if param := params[v.ID]; param != nil {
				s.metrics.ObserveParamValue(v.ID, v.Value, param)
			}
			// }
		}
	} else {
		s.log.Error(errors.Wrap(err, "get parameter values"))
	}

	if ts%300 == 0 {
		s.metrics.WriteLog(t)
		s.metrics.LastTime = t
	}
}

func (s *DasDeviceAgent) QueryStats(beginTime int64, endTime int64, keys []string) ([]map[string]any, error) {
	data := []map[string]any{}
	if s.metrics == nil {
		return data, errors.New("not supported")
	}

	f, err := os.Open(s.metrics.Logger.Filename)
	if err != nil {
		return data, errors.Wrap(err, "open stats file")
	}
	defer f.Close()

	scaner := bufio.NewScanner(f)
	scaner.Split(bufio.ScanLines)

	for scaner.Scan() {
		line := scaner.Text()

		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		ts := int64(0)
		record := map[string]any{}
		for _, v := range parts {
			if subs := strings.Split(v, "="); len(subs) == 2 {
				if subs[0] == "t" {
					ts = cast.ToInt64(subs[1])
					record[subs[0]] = ts
				} else {
					if subs[1] == "" {
						record[subs[0]] = nil
					} else {
						if v, err := cast.ToUint64E(subs[1]); err == nil {
							record[subs[0]] = v
						} else if v, err := cast.ToInt64E(subs[1]); err == nil {
							record[subs[0]] = v
						} else if v, err := cast.ToFloat64E(subs[1]); err == nil {
							record[subs[0]] = v
						} else {
							record[subs[0]] = cast.ToFloat64(subs[1])
						}
					}
				}
			}
		}
		if ts == 0 {
			continue
		}
		if ts < beginTime || ts > endTime {
			continue
		}
		data = append(data, record)
	}

	return data, nil
}
