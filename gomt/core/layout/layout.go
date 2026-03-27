package layout

import (
	"encoding/json"
	"fmt"
	"gomt/core/model"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func keyFromName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, "___", "_")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")

	return name
}

type Element struct {
	Type    string                `json:"Type,omitempty"`
	Key     string                `json:"Key,omitempty"`
	Name    string                `json:"Name,omitempty"`
	Unique  string                `json:"Unique,omitempty"`
	Style   map[string]any        `json:"Style,omitempty"`
	Items   []*Element            `json:"Items,omitempty"`
	Data    []map[string]*Element `json:"Data,omitempty"`
	Actions map[string]*Element   `json:"Actions,omitempty"`

	OID    string  `json:"OID,omitempty"`
	Access string  `json:"Access,omitempty"`
	Value  *string `json:"Value,omitempty"`
}

func (m *Element) LoadFromData(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return errors.Wrap(err, "unmarshal json data")
	}
	return nil
}

func (m *Element) SetName(name string) {
	m.Name = name
	m.Key = keyFromName(name)
}

func (m *Element) SetValue(value string) {
	m.Value = &value
}

func (m *Element) SetStyle(name string, value any) {
	if m.Style == nil {
		m.Style = map[string]any{}
	}
	m.Style[name] = value
}

func (m *Element) RemoveStyle(name string) {
	if m.Style == nil {
		return
	}
	delete(m.Style, name)
}

func (m *Element) GetAction(name string) *Element {
	if m.Actions == nil {
		return nil
	}
	for k, v := range m.Actions {
		if k == name {
			return v
		}
	}
	return nil
}

func (m *Element) SetAction(name string, value *Element) {
	if m.Actions == nil {
		m.Actions = map[string]*Element{}
	}
	m.Actions[name] = value
}

func (m Element) Dump(prefix string) string {
	items := []string{}
	for _, item := range m.Items {
		items = append(items, item.Dump(prefix+"  "))
	}
	if len(items) > 0 {
		return fmt.Sprintf("%v%v: %v, %v\n%v", prefix, m.Type, m.Name, len(m.Items), strings.Join(items, "\n"))
	}
	return fmt.Sprintf("%v%v: %v, %v", prefix, m.Type, m.Name, m.OID)
}

func (m Element) getItemsByType(typ string, nextLevel bool) []*Element {
	items := []*Element{}
	for _, item := range m.Items {
		if item.Type == typ {
			items = append(items, item)
		}
		if nextLevel {
			if _items := item.getItemsByType(typ, nextLevel); len(_items) > 0 {
				items = append(items, _items...)
			}
		}
	}
	return items
}

func (m Element) GetItemsByType(typ string) []*Element {
	return m.getItemsByType(typ, true)
}

func (m Element) GetChildItemsByType(typ string) []*Element {
	return m.getItemsByType(typ, false)
}

func (m Element) GetItemByType(typ string) *Element {
	for _, item := range m.Items {
		if item.Type == typ {
			return item
		}
		if found := item.GetItemByType(typ); found != nil {
			return found
		}
	}
	return nil
}

func (m Element) ReplaceItem(a *Element, b *Element) bool {
	for i, item := range m.Items {
		if item == a {
			m.Items[i] = b
			return true
		}
		if found := item.ReplaceItem(a, b); found {
			return found
		}
	}
	return false
}

func (m Element) GetParams(accessList []model.Access, ignoreInputs []string) []*Element {
	params := []*Element{}
	switch m.Type {
	case "Param":
		match := false
		ignore := false

		if input, ok := m.Style["input"]; ok {
			for _, v := range ignoreInputs {
				if input == v {
					ignore = true
					break
				}
			}
		}
		if !ignore {
			if len(accessList) == 0 {
				match = true
			} else {
				for _, v := range accessList {
					if v == model.Access(m.Access) {
						match = true
						break
					}
				}
			}
		}

		if match {
			params = append(params, &m)
		}
	case "Table":
		for _, row := range m.Data {
			for _, col := range row {
				itemParams := col.GetParams(accessList, ignoreInputs)
				params = append(params, itemParams...)
			}
		}
	default:
		for _, item := range m.Items {
			itemParams := item.GetParams(accessList, ignoreInputs)
			params = append(params, itemParams...)
		}
	}

	paramsMap := map[string]*Element{}
	params2 := []*Element{}
	for _, param := range params {
		tmp := param
		if _, ok := paramsMap[tmp.OID]; !ok {
			paramsMap[tmp.OID] = tmp
			params2 = append(params2, tmp)
		}
	}
	return params2
}

func (m Element) GetWritableParams(ignoreInputs []string) []*Element {
	return m.GetParams([]model.Access{model.AccessWriteOnly, model.AccessReadWrite}, ignoreInputs)
}

func (m Element) GetReadableParams(ignoreInputs []string) []*Element {
	return m.GetParams([]model.Access{model.AccessReadOnly, model.AccessReadWrite}, ignoreInputs)
}

func (m Element) GetItemByTypeWithStyle(typ string, style map[string]any) *Element {
	for _, item := range m.Items {
		if item.Type == typ {
			match := true
			for k, v := range style {
				if v2, ok := item.Style[k]; !ok || v2 != v {
					match = false
					break
				}
			}
			if match {
				return item
			}
		}

		if found := item.GetItemByTypeWithStyle(typ, style); found != nil {
			return found
		}
	}
	return nil
}

func (m *Element) FilterItemsByKeys(keys []string) {
	items := []*Element{}
	for _, v := range m.Items {
		found := false
		for _, key := range keys {
			if key == v.Key {
				found = true
				break
			}
		}
		if !found {
			tmp := v
			items = append(items, tmp)
		}
	}
	m.Items = items
}

func NewRowLayout(name string) *Element {
	return &Element{
		Type:  "Layout:Row",
		Name:  name,
		Key:   keyFromName(name),
		Items: []*Element{},
	}
}

func NewColLayout(name string, span int) *Element {
	out := &Element{
		Type:  "Layout:Col",
		Name:  name,
		Key:   keyFromName(name),
		Items: []*Element{},
	}
	if span > 0 {
		out.Style = map[string]any{
			"span": span,
		}
	}
	return out
}

func NewTabsLayout(name string, pos string) *Element {
	if pos == "" {
		pos = "left"
	}
	return &Element{
		Type:  "Layout:Tabs",
		Name:  name,
		Key:   keyFromName(name),
		Items: []*Element{},
		Style: map[string]any{
			"position": pos,
		},
	}
}
func (m Element) IsTabsLayout() bool {
	return m.Type == "Layout:Tabs"
}

func NewRowLayoutWithItems(items ...*Element) *Element {
	layout := NewRowLayout("")
	layout.Items = append(layout.Items, items...)
	return layout
}

func NewColLayoutWithItems(span int, items ...*Element) *Element {
	layout := NewColLayout("", span)
	layout.Items = append(layout.Items, items...)
	return layout
}

func NewSingleColLayoutWithItems(items ...*Element) *Element {
	layoutSpan := 24
	col := NewColLayoutWithItems(layoutSpan, items...)
	return col
}

func NewApplication(name string) *Element {
	return &Element{
		Type:  "Application",
		Name:  name,
		Key:   keyFromName(name),
		Items: []*Element{},
	}
}

func NewModule(name string) *Element {
	return &Element{
		Type:  "Module",
		Name:  name,
		Key:   keyFromName(name),
		Items: []*Element{},
	}
}

func NewPage(name string) *Element {
	return &Element{
		Type:  "Page",
		Name:  name,
		Key:   keyFromName(name),
		Items: []*Element{},
	}
}

func NewPageWithLayouts(name string, layouts ...*Element) *Element {
	page := NewPage(name)
	page.Items = append(page.Items, layouts...)
	return page
}

func NewPageWithSetParameterValuesFormItems(name string, items ...*Element) *Element {
	form := NewSetParameterValuesForm(name, items...)
	return NewPageWithLayouts(name, NewSingleColLayoutWithItems(form))
}

func NewToolbar(name string) *Element {
	return &Element{
		Name:  name,
		Type:  "Toolbar",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
}

func NewToolbarWithItems(name string, items ...*Element) *Element {
	toolbar := NewToolbar(name)
	toolbar.Items = append(toolbar.Items, items...)
	return toolbar
}

func NewForm(name string) *Element {
	return &Element{
		Name:  name,
		Type:  "Form",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
}

func NewFormWithLayouts(name string, layouts ...*Element) *Element {
	form := NewForm(name)
	form.Items = append(form.Items, layouts...)
	return form
}

func NewFormWithItems(name string, items ...*Element) *Element {
	// layout := NewSingleColLayoutWithItems(items...)
	return NewFormWithLayouts(name, items...)
}

func NewSetParameterValuesForm(name string, items ...*Element) *Element {
	form := NewFormWithItems(name, items...)
	wparams := form.GetWritableParams([]string{"button", "buttonGroup"})
	readonly := true
	if len(wparams) > 0 {
		readonly = false
	}
	if !readonly {
		rparams := form.GetReadableParams([]string{"button", "buttonGroup"})

		gaction := NewGetParameterValuesAction(rparams...)
		saction := NewSetParameterValuesAction(wparams...)
		saction.SetStyle("showMessage", true)
		submitAction := NewMultipleActionsAction(saction, gaction)
		form.SetAction("submit", submitAction)
		// submit := NewButtonWithActions("Submit", "Submit", "primary", action, action2)
		// items = append(items, submit)
		// form = NewFormWithItems(name, items...)
	}
	return form
}

func NewDialog(name string) *Element {
	return &Element{
		Name:  name,
		Type:  "Dialog",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
}

func NewDialogWithLayouts(name string, layouts ...*Element) *Element {
	dialog := NewDialog(name)
	dialog.Items = append(dialog.Items, layouts...)
	return dialog
}

func NewDialogWithItems(name string, items ...*Element) *Element {
	layout := NewSingleColLayoutWithItems(items...)
	return NewDialogWithLayouts(name, layout)
}

func NewConfirmDialogWithMessage(name string, msg string, target *Element) *Element {
	dialog := NewDialogWithItems(name, NewText("Message", msg))
	dialog.SetAction("confirm", target)
	return dialog
}

func NewFilesTable(name string) *Element {
	table := &Element{
		Name:  name,
		Type:  "Table:Files",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
	return table
}

func NewUsersTable(name string) *Element {
	table := &Element{
		Name:  name,
		Type:  "Table:Users",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
	return table
}

func NewAlarmLogsTable(name string) *Element {
	table := &Element{
		Name:  name,
		Type:  "Table:AlarmLogs",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
	return table
}

func NewInventoryTable(name string) *Element {
	table := &Element{
		Name:  name,
		Type:  "Table:Inventory",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
	return table
}

func NewFirmwaresTable(name string) *Element {
	table := &Element{
		Name:  name,
		Type:  "Table:Firmwares",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
	return table
}

func NewTableColumn(name string, fixed string, width int, columns []*Element) *Element {
	columns2 := []*Element{}
	for _, col := range columns {
		found := false
		for _, col2 := range columns2 {
			if col2.Name == col.Name {
				found = true
			}
		}
		if !found {
			columns2 = append(columns2, col)
		}
	}
	elem := &Element{
		Name:  name,
		Type:  "TableColumn",
		Key:   keyFromName(name),
		Items: columns2,
	}
	if fixed != "" {
		elem.SetStyle("fixed", fixed)
	}
	if width > 0 {
		elem.SetStyle("width", width)
	}
	return elem
}

func NewTable(name string) *Element {
	table := &Element{
		Name:  name,
		Type:  "Table",
		Key:   keyFromName(name),
		Items: []*Element{},
	}
	return table
}

func (m *Element) addTableColumn(name string, isUnique bool, fixed string, width int, columns []*Element) (*Element, error) {
	if m.Type != "Table" {
		return nil, errors.New("not table type")
	}

	name = NameFormat(name)

	if isUnique {
		m.Unique = keyFromName(name)
	}
	if m.Items == nil {
		m.Items = []*Element{}
	}
	for _, v := range m.Items {
		if v.Name == name {
			return v, nil
		}
	}
	item := NewTableColumn(name, fixed, width, columns)
	m.Items = append(m.Items, item)
	return item, nil
}

func (m *Element) AddTableColumn(name string, width int) *Element {
	item, err := m.addTableColumn(name, false, "", width, nil)
	if err != nil {
		logrus.Fatal(errors.Wrapf(err, "add table column, table=%v, column=%v", m.Name, name))
	}
	return item
}
func (m *Element) AddTableColumnWithColumns(name string, width int, columns []*Element) *Element {
	item, err := m.addTableColumn(name, false, "", width, columns)
	if err != nil {
		logrus.Fatal(errors.Wrapf(err, "add table column, table=%v, column=%v", m.Name, name))
	}
	return item
}
func (m *Element) AddTableColumnUnique(name string, width int) *Element {
	item, err := m.addTableColumn(name, true, "left", width, nil)
	if err != nil {
		logrus.Fatal(errors.Wrapf(err, "add table column, table=%v, column=%v", m.Name, name))
	}
	return item
}

func (m *Element) MustSetTableRowData(index int64, name string, elem *Element) {
	if err := m.SetTableRowData(index, name, elem); err != nil {
		logrus.Warnf("table: %v", m.Dump(""))
		logrus.Warnf("elem: %v", elem.Dump(""))
		logrus.Fatal(errors.Wrapf(err, "set row data for %v", name))
	}
}

func (m *Element) SetTableRowClickActions(index int64, actions ...*Element) error {
	action := NewMultipleActionsAction(actions...)
	return m.SetTableRowAction(index, "click", action)
}

func (m *Element) SetTableRowClickAction(index int64, action *Element) error {
	return m.SetTableRowAction(index, "click", action)
}

func (m *Element) SetTableRowAction(index int64, name string, action *Element) error {
	if index < 0 {
		logrus.Warnf("invalid row: index=%v, action=%v", index, action)
		return errors.New("invalid row index")
	}
	if m.Data == nil {
		m.Data = []map[string]*Element{}
	}
	for i := len(m.Data); i <= int(index); i = len(m.Data) {
		m.Data = append(m.Data, map[string]*Element{})
	}
	row := m.Data[index]
	actions := row["_actions"]
	if actions == nil {
		actions = NewAction("Actions")
	}
	actions.SetAction(name, action)
	row["_actions"] = actions
	m.Data[index] = row
	return nil
}

func (m *Element) SetTableRowData(index int64, name string, elem *Element) error {
	colIndex := -1
	name = NameFormat(name)
	key := keyFromName(name)

	i := 0
	for _, col := range m.Items {
		if col.Items != nil && len(col.Items) > 0 {
			for _, col2 := range col.Items {
				if col2.Name == name {
					colIndex = i
					break
				}
				i += 1
			}
			if colIndex >= 0 {
				break
			}
			i += len(col.Items)
		} else {
			if col.Name == name {
				colIndex = i
				break
			}
			i += 1
		}
	}
	if colIndex < 0 {
		logrus.Warnf("invalid column: name=%v, elem=%v, %v", name, elem, m.Items)
		return errors.New("invalid column name")
	}
	if index < 0 {
		logrus.Warnf("invalid row: name=%v, elem=%v, %v", name, elem, m.Items)
		return errors.New("invalid row index")
	}
	if m.Data == nil {
		m.Data = []map[string]*Element{}
	}
	for i := len(m.Data); i <= int(index); i = len(m.Data) {
		m.Data = append(m.Data, map[string]*Element{})
	}
	row := m.Data[index]
	row[key] = elem
	m.Data[index] = row
	return nil
}

func NewComponent(name string, typ string, items ...*Element) *Element {
	if items == nil {
		items = []*Element{}
	}
	name = NameFormat(name)
	return &Element{
		Type:  "Component:" + typ,
		Name:  name,
		Key:   keyFromName(name),
		Items: items,
	}
}

func NewParamGroupComponent(name string, num int, params ...*Element) *Element {
	elem := NewComponent(name, fmt.Sprintf("ParamGroup:%v", num), params...)
	access := params[0].Access
	for _, v := range params {
		if v.Access != access {
			return elem
		}
	}
	elem.Access = access
	return elem
}

func NewStatisticComponent(name string, label *Element, value *Element) *Element {
	elem := NewComponent(name, "Statistic", label, value)
	elem.Access = "ro"
	return elem
}

func NewStatisticGroupComponent(name string, items ...*Element) *Element {
	elem := NewComponent(name, "StatisticGroup", items...)
	elem.Access = "ro"
	return elem
}

func NewParam(name string, oid string, access string) *Element {
	return newParamWithValue(name, oid, access, nil)
}
func NewParamWithValue(name string, oid string, access string, value string) *Element {
	v := value
	return newParamWithValue(name, oid, access, &v)
}

func newParamWithValue(name string, oid string, access string, value *string) *Element {
	name = NameFormat(name)
	return &Element{
		Type:   "Param",
		Name:   name,
		Key:    keyFromName(name),
		OID:    oid,
		Access: access,
		Value:  value,
	}
}

func CopyParamWithValue(p *Element, value string) *Element {
	return NewParamWithValue(p.Name, p.OID, p.Access, value)
}

func NewLabel(name string, value string) *Element {
	elem := &Element{
		Type: "Label",
		Name: name,
		Key:  keyFromName(name),
	}
	elem.SetValue(value)
	return elem
}

func NewText(name string, msg string) *Element {
	elem := &Element{
		Type: "Text",
		Name: name,
		Key:  keyFromName(name),
	}
	elem.SetValue(msg)
	return elem
}

func NewAlert(name string, value string, typ string) *Element {
	elem := &Element{
		Type: "Alert",
		Name: name,
		Key:  keyFromName(name),
	}
	elem.SetValue(value)
	elem.SetStyle("type", typ)
	return elem
}

func NewButton(name string, value string, typ string) *Element {
	btn := &Element{
		Type:  "Button",
		Name:  name,
		Key:   keyFromName(name),
		Value: &value,
		Style: map[string]any{
			"type": typ,
		},
	}
	return btn
}

func NewButtonWithAction(name string, value string, typ string, action *Element) *Element {
	btn := NewButton(name, value, typ)
	btn.SetAction("click", action)
	return btn
}

func NewButtonWithViewPage(name string, value string, typ string, page *Element) *Element {
	action := NewViewPageAction(page)
	return NewButtonWithAction(name, value, typ, action)
}

func NewButtonWithActions(name string, value string, typ string, actions ...*Element) *Element {
	action := NewMultipleActionsAction(actions...)
	return NewButtonWithAction(name, value, typ, action)
}

func NewAction(name string, items ...*Element) *Element {
	return &Element{
		Type:  "Action",
		Name:  name,
		Key:   keyFromName(name),
		Items: items,
	}
}

func NewMultipleActionsAction(items ...*Element) *Element {
	return NewAction("MultipleActions", items...)
}

func NewViewPageAction(page *Element) *Element {
	return NewAction("ViewPage", page)
}

func NewSetParameterValuesAction(params ...*Element) *Element {
	re := regexp.MustCompile(`^([A-Z0-9\.\-_]+)(\[\d+\]|)$`)
	check := map[string]bool{}

	elems := []*Element{}
	for _, param := range params {
		if subs := re.FindStringSubmatch(param.OID); len(subs) == 3 {
			if _, ok := check[subs[1]]; !ok {
				check[subs[1]] = true
				elem := &Element{
					OID:   subs[1],
					Value: param.Value,
				}
				elems = append(elems, elem)
			}
		} else {
			logrus.Fatalf("make action for SetParameterValues error, invalid oid '%v'", param.OID)
		}
	}
	return NewAction("SetParameterValues", elems...)
}

func NewGetParameterValuesAction(params ...*Element) *Element {
	re := regexp.MustCompile(`^([A-Z0-9\.\-_]+)(\[\d+\]|)$`)
	check := map[string]bool{}

	elems := []*Element{}
	for _, param := range params {
		if subs := re.FindStringSubmatch(param.OID); len(subs) == 3 {
			if _, ok := check[subs[1]]; !ok {
				check[subs[1]] = true
				elem := &Element{
					OID:   subs[1],
					Value: param.Value,
				}
				elems = append(elems, elem)
			}
		} else {
			logrus.Fatalf("make action for GetParameterValues error, invalid oid '%v'", param.OID)
		}
	}
	return NewAction("GetParameterValues", elems...)
}

type PageMap struct {
	Pages map[string]*Element
	Keys  []string
}

func GetPageMap(app *Element, typ string) PageMap {
	out := PageMap{
		Pages: map[string]*Element{},
		Keys:  []string{},
	}
	for _, module := range app.Items {
		if module.Type != "Module" {
			continue
		}
		for _, page := range module.Items {
			if page.Type != "Page" {
				continue
			}
			if len(page.Items) == 1 && page.Items[0].IsTabsLayout() {
				for _, tabPage := range page.Items[0].Items {
					if tabPage.Type != "Page" {
						continue
					}
					tmp := *tabPage
					key := fmt.Sprintf("%v.%v.%v", module.Key, page.Key, tabPage.Key)
					out.Pages[key] = &tmp
					out.Keys = append(out.Keys, key)
				}
			} else {
				tmp := *page
				key := fmt.Sprintf("%v.%v", module.Key, page.Key)
				out.Pages[key] = &tmp
				out.Keys = append(out.Keys, key)
			}
		}
	}
	return out
}

func (m PageMap) GetPage(key string) *Element {
	for k, v := range m.Pages {
		if k == key {
			tmp := v
			return tmp
		}
	}
	return nil
}

func (m PageMap) GetPageMapWithPrefix(prefix string) PageMap {
	out := PageMap{
		Pages: map[string]*Element{},
		Keys:  []string{},
	}
	for _, k := range m.Keys {
		v := m.Pages[k]
		if strings.HasPrefix(k, prefix) {
			out.Pages[k] = v
			out.Keys = append(out.Keys, k)
		}
	}
	return out
}
