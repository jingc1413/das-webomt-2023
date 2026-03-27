package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Options map[string]string

func NewOptionsFromMap(in map[string]string) Options {
	out := Options{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func NewOptionsFromOptionsOfHexKey(in Options) Options {
	if in == nil {
		return nil
	}
	out := Options{}
	for k, v := range in {
		i, _ := strconv.ParseInt(k, 16, 32)
		out[fmt.Sprintf("%v", i)] = v
	}
	return out
}

func (m Options) Equal(n Options) bool {
	for k, v := range n {
		if v2, ok := m[k]; !ok {
			return false
		} else if v2 != v {
			return false
		}
	}
	for k, v := range m {
		if v2, ok := n[k]; !ok {
			return false
		} else if v2 != v {
			return false
		}
	}
	return true
}
func (m Options) IDs() []string {
	pairs := []*SortOptionPair{}
	for k := range m {
		pairs = append(pairs, &SortOptionPair{k, k})
	}
	sort.Sort(BySortOptionPair{pairs})
	ids := []string{}
	for _, pair := range pairs {
		ids = append(ids, pair.Key())
	}
	return ids
}

func (m Options) String() string {
	v, _ := m.MarshalJSON()
	return strings.ReplaceAll(string(v), "\n", "")
}

func (m Options) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(true)
	ids := m.IDs()
	for i, id := range ids {
		v := m[id]
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := encoder.Encode(id); err != nil {
			return nil, err
		}
		buf.WriteByte(':')
		if err := encoder.Encode(v); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

type SortOptionPair struct {
	key   string
	value string
}

func (kv *SortOptionPair) Key() string   { return kv.key }
func (kv *SortOptionPair) Value() string { return kv.value }

type BySortOptionPair struct {
	Pairs []*SortOptionPair
}

func (a BySortOptionPair) Len() int           { return len(a.Pairs) }
func (a BySortOptionPair) Swap(i, j int)      { a.Pairs[i], a.Pairs[j] = a.Pairs[j], a.Pairs[i] }
func (a BySortOptionPair) Less(i, j int) bool { return a.Pairs[i].Value() < a.Pairs[j].Value() }
