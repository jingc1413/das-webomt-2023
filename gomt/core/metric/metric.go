package metric

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"gomt/core/logger/lumberjack"
	"gomt/core/model"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type MetricItem struct {
	ID    string
	Type  string
	Fixed int
}

type MetricsData struct {
	sync.Mutex
	Logger     *lumberjack.Logger
	LastTime   time.Time
	Params     []*model.ParameterDefine
	Items      []MetricItem
	SummaryMap map[string][]any
	GaugeMap   map[string]any
}

func NewMetricsData(filename string, maxsize, maxbackups int) *MetricsData {
	return &MetricsData{
		Params:     []*model.ParameterDefine{},
		Items:      []MetricItem{},
		SummaryMap: map[string][]any{},
		GaugeMap:   map[string]any{},
		Logger: &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxsize, // kilobyte
			MaxBackups: maxbackups,
			// Compress:   true, // disabled by default
		},
	}
}

func (m *MetricsData) GetItem(id string) *MetricItem {
	m.Lock()
	defer m.Unlock()
	id = strings.ReplaceAll(id, "|", ",")
	for _, item := range m.Items {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

func (m *MetricsData) AddParamItem(def *model.ParameterDefine, fixed int) {
	m.Lock()
	defer m.Unlock()
	id := fmt.Sprintf("%v", def.PrivOid)
	id = strings.ReplaceAll(id, "|", ",")
	m.Params = append(m.Params, def)
	m.Items = append(m.Items, MetricItem{ID: id, Type: string(def.DataType), Fixed: fixed})
}

func (m *MetricsData) AddItem(id string, typ string, fixed int) {
	m.Lock()
	defer m.Unlock()

	id = strings.ReplaceAll(id, "|", ",")
	m.Items = append(m.Items, MetricItem{ID: id, Type: typ, Fixed: fixed})
}

func (m *MetricsData) Set(id string, value any) {
	m.Lock()
	defer m.Unlock()
	id = strings.ReplaceAll(id, "|", ",")
	for _, item := range m.Items {
		if item.ID == id {
			if item.Fixed > 0 {
				v := float64Fixed(cast.ToFloat64(value), item.Fixed)
				m.GaugeMap[id] = v
			} else {
				m.GaugeMap[id] = value
			}
			return
		}
	}
}

func (m *MetricsData) SetInt64(id string, value int64) {
	m.Set(id, value)
}

func (m *MetricsData) SetUint64(id string, value uint64) {
	m.Set(id, value)
}
func (m *MetricsData) SetFloat64(id string, value float64) {
	m.Set(id, value)
}

func (m *MetricsData) Observe(id string, value any) {
	m.Lock()
	defer m.Unlock()

	id = strings.ReplaceAll(id, "|", ",")

	for _, item := range m.Items {
		if item.ID == id {
			summary := m.SummaryMap[id]
			if summary == nil {
				summary = []any{}
			}

			if item.Fixed > 0 {
				v := float64Fixed(cast.ToFloat64(value), item.Fixed)
				summary = append(summary, v)
				m.GaugeMap[id] = v
			} else {
				summary = append(summary, value)
				m.GaugeMap[id] = value
			}

			m.SummaryMap[id] = summary
			return
		}
	}

}

func (m *MetricsData) ObserveFloat64(id string, value float64) {
	m.Observe(id, value)
}

func (m *MetricsData) ObserveInt64(id string, value int64) {
	m.Observe(id, value)
}

func (m *MetricsData) ObserveUint64(id string, value uint64) {
	m.Observe(id, value)
}

func (m *MetricsData) ObserveParamValue(id string, value any, param *model.ParameterDefine) {
	if param.Ratio != nil {
		radio := *param.Ratio
		m.ObserveFloat64(id, cast.ToFloat64(value)/float64(radio))
		return
	}
	switch param.DataType {
	case model.DataTypeInt, model.DataTypeInt16, model.DataTypeInt32, model.DataTypeInt64:
		m.ObserveInt64(id, cast.ToInt64(value))
	case model.DataTypeUInt, model.DataTypeUInt16, model.DataTypeUInt32, model.DataTypeUInt64:
		m.ObserveUint64(id, cast.ToUint64(value))
	case model.DataTypeFloat32, model.DataTypeFloat64:
		m.ObserveFloat64(id, cast.ToFloat64(value))
	default:
		m.ObserveFloat64(id, cast.ToFloat64(value))
	}
}

func float64Fixed(value float64, fixed int) float64 {
	if fixed > 0 {
		k := float64(1)
		for i := 0; i < fixed; i++ {
			k = k * 10
		}
		value = math.Round(value*k) / k
	}
	return value
}

func (m *MetricsData) WriteLog(t time.Time) {
	m.Lock()
	defer m.Unlock()

	parts := []string{
		fmt.Sprintf("t=%v", t.Unix()),
	}
	for _, item := range m.Items {
		if summary, ok := m.SummaryMap[item.ID]; ok && summary != nil {
			sum := float64(0)
			if len(summary) > 0 {
				for _, v := range summary {
					sum += cast.ToFloat64(v)
				}
				mean := sum / float64(len(summary))
				mean = float64Fixed(mean, item.Fixed)
				parts = append(parts, fmt.Sprintf("%v=%v", item.ID, mean))
			} else {
				parts = append(parts, fmt.Sprintf("%v=", item.ID))
			}
			m.SummaryMap[item.ID] = []any{}
		} else if v, ok := m.GaugeMap[item.ID]; ok {
			parts = append(parts, fmt.Sprintf("%v=%v", item.ID, v))
		}
	}
	content := strings.Join(parts, "|") + "\n"

	if _, err := m.Logger.Write([]byte(content)); err != nil {
		logrus.Error(errors.Wrap(err, "write log"))
	}
	m.Logger.Close()
}
