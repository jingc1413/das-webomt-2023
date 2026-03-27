package story

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"gomt/core/layout"
	"gomt/core/model"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:embed templates
var res embed.FS

func readTemplate(name string) []byte {
	data, err := res.ReadFile("templates/" + name + ".tmpl")
	if err != nil {
		logrus.Fatalf("read template file error, name=%v", name)
	}
	return data
}

type UserStory struct {
	Summary     string
	Description string
}

func MakeApplictionStoryList(app *layout.Element, productDef *model.ProductDefine, defs model.ParameterDefines) []UserStory {
	out := []UserStory{}
	for _, module := range app.Items {
		if out2 := exportModuleStoryList(app, module, productDef, defs); len(out2) > 0 {
			out = append(out, out2...)
		}
	}
	return out
}

func exportModuleStoryList(app *layout.Element, module *layout.Element, productDef *model.ProductDefine, defs model.ParameterDefines) []UserStory {
	out := []UserStory{}
	if module.Type != "Module" {
		return out
	}
	for _, page := range module.Items {
		if tabs := page.GetItemByType("Layout:Tabs"); tabs != nil {
			for _, tabPage := range tabs.Items {
				page2 := tabPage
				prefix := fmt.Sprintf("%v-%v-%v", module.Name, page.Name, page2.Name)
				prefix = strings.ReplaceAll(prefix, " ", "")
				if out2 := exportPageStoryList(prefix, app, module, page2, productDef, defs); len(out2) > 0 {
					out = append(out, out2...)
				}
			}
		} else {
			prefix := fmt.Sprintf("%v-%v", module.Name, page.Name)
			prefix = strings.ReplaceAll(prefix, " ", "")
			if out2 := exportPageStoryList(prefix, app, module, page, productDef, defs); len(out2) > 0 {
				out = append(out, out2...)
			}
		}
	}
	return out
}

func exportPageStoryList(breadcrumb string, app *layout.Element, module *layout.Element, page *layout.Element,
	productDef *model.ProductDefine, defs model.ParameterDefines) []UserStory {
	out := []UserStory{}
	if page.Type != "Page" {
		return out
	}

	prefix := fmt.Sprintf("【%v】", breadcrumb)
	fmap := template.FuncMap{
		"getParam": func(elem *layout.Element) *model.ParameterDefine {
			def := defs.GetParameterDefine(model.PrivObjectId(elem.OID))
			return def
		},
		"getAction": func(elem *layout.Element, name string) *layout.Element {
			if elem.Actions != nil {
				if action, ok := elem.Actions[name]; ok {
					return action
				}
			}
			return nil
		},
		"getItem": func(elem *layout.Element, typ string) *layout.Element {
			return elem.GetItemByType(typ)
		},
		"getItems": func(elem *layout.Element, typ string) []*layout.Element {
			return elem.GetItemsByType(typ)
		},
		"getTableSize": func(elem *layout.Element) int {
			return len(elem.Data)
		},
		"isWritable": func(elem *layout.Element) bool {
			out := elem.GetWritableParams([]string{"button", "buttonGroup"})
			return len(out) > 0
		},
		"isProductAU": func() bool {
			return productDef.ProductTypeName == "AU"
		},
		"supportExportTable": func(table *layout.Element) bool {
			if unique := table.Data[0][table.Unique]; unique != nil {
				return true
			}
			return false
		},
		"supportImportTable": func(table *layout.Element) bool {
			if actions := table.Data[0]["_actions"]; actions != nil {
				if click := actions.Actions["click"]; click != nil {
					if table.Name == "RF Module Mapping" || table.Name == "SNMP User List" {
						//ignore
					} else {
						return true
					}
				}
			}
			return false
		},
		"supportViewTableItem": func(table *layout.Element) bool {
			if table.Name == "RF Module Mapping" {
				return true
			}
			if actions := table.Data[0]["_actions"]; actions != nil {
				if click := actions.Actions["click"]; click != nil {
					return true
				}
			}
			return false
		},
		"supportEditTableItem": func(table *layout.Element) bool {
			if table.Name == "RF Module Mapping" {
				return true
			}
			if actions := table.Data[0]["_actions"]; actions != nil {
				if click := actions.Actions["click"]; click != nil {
					if tmp := click.GetWritableParams([]string{"button", "buttonGroup"}); len(tmp) > 0 {
						return true
					}
				}
			}
			return false
		},
	}
	if strings.HasPrefix(breadcrumb, "Overview-DASTopo") {
		// ignore
	} else if strings.HasPrefix(breadcrumb, "SystemSettings-Account-Users") {
		// ignore
	} else if strings.HasPrefix(breadcrumb, "SystemSettings-Configuration") {
		// ignore
	} else if strings.HasPrefix(breadcrumb, "SystemSettings-Upgrade") {
		// ignore
	} else if strings.HasPrefix(breadcrumb, "SystemSettings-Logs") {
		// ignore
	} else if strings.HasPrefix(breadcrumb, "Maintenance-FirmwareInformation") {
		// ignore
	} else if strings.HasPrefix(breadcrumb, "FactoryMaintenance-AddressInterface") {
		// ignore
	} else {
		if form := page.GetItemByType("Form"); form != nil {
			out = appendStory(out, "form_mgmt", fmap, prefix, page.Name, form)
		}

		tables := page.GetItemsByType("Table")
		for _, v := range tables {
			if table := v; table != nil && len(table.Data) > 0 {
				name := table.Name
				if strings.HasSuffix(name, "List") {
					name = strings.ReplaceAll(name, "List", "")
					name = strings.TrimSpace(name)
				}
				out = appendStory(out, "table_mgmt", fmap, prefix, name, table)
			}
		}
	}

	return out
}

type TemplateData struct {
	Data *layout.Element
	Name string
}

func appendStory(out []UserStory,
	tmplName string,
	fmap template.FuncMap,
	prefix string,
	pageName string,
	data *layout.Element,
) []UserStory {
	if summary, content := executeStoryTemplate(tmplName, fmap, pageName, data); summary != "" {
		out = append(out, UserStory{
			Summary:     fmt.Sprintf("%v %v", prefix, summary),
			Description: content,
		})
	}
	return out
}

func executeStoryTemplate(tmplTye string, fmap template.FuncMap, name string, data *layout.Element) (string, string) {
	summary := ""
	content := ""
	input := TemplateData{
		Name: name,
		Data: data,
	}

	tmpl := template.New("name")
	tmpl = tmpl.Funcs(fmap)
	tmpl, err := tmpl.Parse(string(readTemplate("story_summary")))
	if err != nil {
		logrus.Error(errors.Wrapf(err, "parse template file, %v, %v", data.Name, tmplTye))
		return summary, content
	}
	tmpl, err = tmpl.Parse(string(readTemplate(tmplTye)))
	if err != nil {
		logrus.Error(errors.Wrapf(err, "parse template file, %v, %v", data.Name, tmplTye))
		return summary, content
	}
	buf := bytes.NewBufferString("")
	if err := tmpl.Execute(buf, input); err != nil {
		logrus.Error(errors.Wrapf(err, "execute template, %v, %v", data.Name, tmplTye))
		return summary, content
	}
	summary = buf.String()

	tmpl2 := template.New("story")
	tmpl2 = tmpl2.Funcs(fmap)
	tmpl2, err = tmpl2.Parse(string(readTemplate("story_desc")))
	if err != nil {
		logrus.Error(errors.Wrapf(err, "parse template file, %v, %v", data.Name, tmplTye))
		return summary, content
	}
	tmpl2, err = tmpl2.Parse(string(readTemplate(tmplTye)))
	if err != nil {
		logrus.Error(errors.Wrapf(err, "parse template file, %v, %v", data.Name, tmplTye))
		return summary, content
	}
	buf2 := bytes.NewBufferString("")
	if err := tmpl2.Execute(buf2, input); err != nil {
		logrus.Error(errors.Wrapf(err, "execute template, %v, %v", data.Name, tmplTye))
		return summary, content
	}
	content = buf2.String()
	content = strings.ReplaceAll(content, "\n\n", "\n")
	return summary, content
}

func UpdateStoryListToJira(productDef *model.ProductDefine, version string, list []UserStory) error {
	for _, v := range list {
		// logrus.Warnf("%v, %v", v.Summary, v.Description)
		// logrus.Warnf("%v", v.Summary)
		// logrus.Infof("%v, %v, %v, %v", productDef.Schema, productDef.DeviceTypeName, v.Summary, len(v.Description))
		componentName := fmt.Sprintf("%v:%v:%v", productDef.Schema, productDef.DeviceTypeName, version)
		_, err := UpdateOrCreateIssue("11111", componentName, "故事", v.Summary, v.Description)
		if err != nil {
			logrus.Warnf("%v, %v, %v", componentName, v.Summary, len(v.Description))
			// logrus.Warnf("%v", v.Description)
			return err
		}
	}

	return nil
}
