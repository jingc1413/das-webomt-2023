package layout

import (
	"fmt"
	"gomt/core/model"
	"gomt/core/parser"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var defaultAllLayoutMutex sync.Mutex
var defaultAllLayoutSetup bool
var defaultAllLayoutMap map[string]*Element
var defaultAllPathsMap map[string][]string

func GetLayoutByDeviceTypeName(deviceTypeName string, version string) *Element {
	version = model.FormatVersion(version)

	all := GetAllLayoutMap()
	key := fmt.Sprintf("%v:%v", deviceTypeName, version)
	v := all[key]
	return v
}

func GetAllLayoutMap() map[string]*Element {
	defaultAllLayoutMutex.Lock()
	defer defaultAllLayoutMutex.Unlock()
	if !defaultAllLayoutSetup {
		logrus.Fatal("all parameter defines is not setup")
	}
	return defaultAllLayoutMap
}
func GetAllPathsMap() map[string][]string {
	defaultAllLayoutMutex.Lock()
	defer defaultAllLayoutMutex.Unlock()
	if !defaultAllLayoutSetup {
		logrus.Fatal("all parameter defines is not setup")
	}
	return defaultAllPathsMap
}

func SetupAllLayoutMap(schema string, allModels map[string]model.ParameterDefines) error {
	defaultAllLayoutMutex.Lock()
	defer defaultAllLayoutMutex.Unlock()

	all, err := loadAllLayoutMapFromModels(schema, allModels)
	if err != nil {
		return err
	}

	pagesMap := map[string][]string{}
	for key, layout := range all {
		_paths := map[string]bool{}
		for _, m := range layout.Items {
			if m.Type == "Module" {
				pgs := m.getItemsByType("Page", false)
				for _, pg := range pgs {

					if tabs := pg.getItemsByType("Layout:Tabs", false); len(tabs) > 0 {
						pgs2 := tabs[0].getItemsByType("Page", false)
						for _, pg2 := range pgs2 {
							path := fmt.Sprintf("%v.%v.%v", m.Key, pg.Key, pg2.Key)
							_paths[path] = true
						}
					} else {
						path := fmt.Sprintf("%v.%v", m.Key, pg.Key)
						_paths[path] = true
					}
				}
			}
		}
		paths := []string{}
		for path, _ := range _paths {
			paths = append(paths, path)
		}
		pagesMap[key] = paths
	}

	defaultAllLayoutSetup = true
	defaultAllLayoutMap = all
	defaultAllPathsMap = pagesMap
	return nil
}

func loadAllLayoutMapFromModels(schema string, allModels map[string]model.ParameterDefines) (map[string]*Element, error) {
	out := map[string]*Element{}
	alarmDefs := parser.GetAlarmDefines()

	for deviceTypeName, defs := range allModels {
		productDef := model.GetProductDefine(schema, deviceTypeName)
		if productDef == nil {
			return out, errors.Errorf("invalid device type, device=%v", deviceTypeName)
		}
		if !productDef.SupportLayout {
			continue
		}
		versions := productDef.Versions
		for _, version := range versions {
			params := defs.GetParameterDefinesByVersion(schema, version)
			app, err := MakeApplication(deviceTypeName, DeviceInfo{
				Device: productDef,
				Params: params,
				Alarms: alarmDefs,
			})
			if err != nil {
				return out, errors.Wrapf(err, "make application layout, device=%v, version=%v", deviceTypeName, version)
			}
			key := fmt.Sprintf("%v:%v", deviceTypeName, version)
			// key := deviceTypeName
			out[key] = app
		}
		if version := "latest"; version != "" {
			params := defs.GetParameterDefinesByVersion(schema, version)
			app, err := MakeApplication(deviceTypeName, DeviceInfo{
				Device: productDef,
				Params: params,
				Alarms: alarmDefs,
			})
			if err != nil {
				return out, errors.Wrapf(err, "make application layout, device=%v, version=%v", deviceTypeName, version)
			}
			key := fmt.Sprintf("%v:%v", deviceTypeName, version)
			// key := deviceTypeName
			out[key] = app
		}
	}
	return out, nil
}

// func SetupAllLayoutMap(fsys fs.FS, base string) error {
// 	def=ltAllLayoutMutex.Lock()
// 	defer defaultAllLayoutMutex.Unlock()

// 	all, err := loadAllLayoutMapFromFS(fsys, base)
// 	if err != nil {
// 		return err
// 	}
// 	list := []*Element{}
// 	for _, v := range all {
// 		tmp := v
// 		list = append(list, tmp)
// 	}
// 	defaultAllLayoutSetup = true
// 	defaultAllLayoutMap = all
// 	defaultAllLayoutList = list
// 	return nil
// }

// func loadAllLayoutMapFromFS(fsys fs.FS, base string) (map[string]*Element, error) {
// 	out := map[string]*Element{}

// 	entries, err := fs.ReadDir(fsys, base)
// 	if err != nil {
// 		return out, errors.Wrap(err, "read models dir")
// 	}
// 	for _, entry := range entries {
// 		if entry.IsDir() {
// 			if deviceTypeName := entry.Name(); deviceTypeName != "" {
// 				filename := filepath.Join(base, deviceTypeName, "layout.json")
// 				filename = filepath.ToSlash(filepath.Clean(filename))
// 				f, err := fsys.Open(filename)
// 				if err != nil {
// 					return out, errors.Wrap(err, "read file "+deviceTypeName)
// 				}
// 				data, err := io.ReadAll(f)
// 				if err != nil {
// 					return out, errors.Wrap(err, "decompress file "+deviceTypeName)
// 				}
// 				def := Element{}
// 				if err := def.LoadFromData(data); err != nil {
// 					return out, errors.Wrap(err, "load "+deviceTypeName)
// 				}
// 				out[deviceTypeName] = &def
// 			}
// 		}
// 	}
// 	for k, _ := range out {
// 		logrus.Tracef("load model layout: %v", k)
// 	}
// 	return out, nil
// }
