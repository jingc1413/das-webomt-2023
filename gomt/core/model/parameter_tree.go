package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type ParameterTreeNode struct {
	Name   string
	Child  []*ParameterTreeNode
	Params []*ParameterDefine
}

var (
	reFrequencyRange = regexp.MustCompile(`(.*)(UL|Uplink|DL|Downlink)\s+(Freq|Frequency)\s+(Start|Begin|Low|End|High)`)
	reWorkingHours   = regexp.MustCompile(`Working\s*Hours\s*(Start|End)`)
)

func (m ParameterTreeNode) GetChild(name string) *ParameterTreeNode {
	for _, v := range m.Child {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func GetParameterTree(params ParameterDefines) *ParameterTreeNode {
	root := &ParameterTreeNode{
		Name:  "root",
		Child: []*ParameterTreeNode{},
	}
	re := regexp.MustCompile(`\s+`)

	for _, param := range params {
		for _, v := range param.Groups {
			v = re.ReplaceAllString(v, " ")
			v = strings.TrimSpace(v)

			args := strings.Split(v, ",")
			if len(args) != 3 {
				logrus.Warnf("invalid group %v", v)
				continue
			}

			module := root.GetChild(args[0])
			if module == nil {
				module = &ParameterTreeNode{
					Name:  args[0],
					Child: []*ParameterTreeNode{},
				}
				root.Child = append(root.Child, module)
			}

			page := module.GetChild(args[1])
			if page == nil {
				page = &ParameterTreeNode{
					Name:  args[1],
					Child: []*ParameterTreeNode{},
				}
				module.Child = append(module.Child, page)
			}

			group := page.GetChild(args[2])
			if group == nil {
				group = &ParameterTreeNode{
					Name:  args[2],
					Child: []*ParameterTreeNode{},
				}
				page.Child = append(page.Child, group)
			}

			var obj *ParameterTreeNode
			if subs := reFrequencyRange.FindStringSubmatch(param.Name); len(subs) == 5 {
				name := fmt.Sprintf("%v%v Frequency", subs[1], subs[2])
				obj = group.GetChild(name)
				if obj == nil {
					obj = &ParameterTreeNode{
						Name:   name,
						Params: []*ParameterDefine{},
					}
					group.Child = append(group.Child, obj)
				}
			} else if subs := reWorkingHours.FindStringSubmatch(param.Name); len(subs) == 2 {
				name := "Working Hours (24h)"
				obj = group.GetChild(name)
				if obj == nil {
					obj = &ParameterTreeNode{
						Name:   name,
						Params: []*ParameterDefine{},
					}
					group.Child = append(group.Child, obj)
				}
			} else {
				obj = &ParameterTreeNode{
					Name:   param.Name,
					Params: []*ParameterDefine{},
				}
				group.Child = append(group.Child, obj)
			}
			obj.Params = append(obj.Params, param)
		}

	}

	return root
}
