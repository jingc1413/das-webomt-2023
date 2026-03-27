package cmd

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gomt/core/layout"
	"gomt/core/model"
	"gomt/core/parser"
	"gomt/core/story"
	"gomt/core/utils"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var GenerateCommand = &cobra.Command{
	Use:   fmt.Sprintf("generate [--import filepath] [-o output]"),
	Short: fmt.Sprintf("generate json files for webomt"),
	Long:  fmt.Sprintf("generate json files for webomt"),
	Run:   generate,
}

var (
	importPath string
	outputPath string
)

func init() {
	GenerateCommand.PersistentFlags().StringVarP(&importPath, "import", "", "asserts/device-files", "device files path")
	GenerateCommand.PersistentFlags().StringVarP(&outputPath, "out", "o", "out/", "output file path")
}

func generate(cmd *cobra.Command, args []string) {
	if err := doGenerate2(); err != nil {
		logrus.Error(err)
	}
}

func doGenerate2() error {
	all := map[string]map[string]model.ParameterDefines{}
	alarmDefs := parser.GetAlarmDefines()
	dirs, err := utils.GetDirList(importPath)
	if err != nil {
		return errors.Wrap(err, "read import dir")
	}
	for _, dir := range dirs {
		name := path.Base(dir)

		// parts := strings.Split(name, "-")
		// partsLength := len(parts)
		// version := parts[partsLength-1]
		// schema := "default"
		// if partsLength > 2 {
		// 	schema = strings.ToLower(parts[1])
		// }

		logrus.Infof("load files %v %v", dir, name)

		infoFiles, err := utils.GetFileList(dir, ".yaml")
		if err != nil {
			return errors.Wrap(err, "get info files")
		}
		if len(infoFiles) < 1 {
			return errors.New("cant find info file")
		}
		info, err := loadInfo(infoFiles[0])
		if err != nil {
			return errors.Wrap(err, "load info")
		}

		schemaParamsMap, ok := all[info.Schema]
		if !ok || schemaParamsMap == nil {
			schemaParamsMap = map[string]model.ParameterDefines{}
			all[info.Schema] = schemaParamsMap
		}

		xmlFiles, err := utils.GetFileList(dir, ".xml")
		if err != nil {
			return errors.Wrap(err, "get xml files")
		}
		iniFiles, err := utils.GetFileList(dir, ".ini")
		if err != nil {
			return errors.Wrap(err, "get ini files")
		}

		paramsMap, err := parseDeviceFiles(info.Schema, xmlFiles, iniFiles)
		if err != nil {
			return errors.Wrap(err, "parse device files")
		}

		for deviceTypeName, importParams := range paramsMap {
			version, revision, err := info.GetVersion(deviceTypeName)
			if err != nil {
				return errors.Wrap(err, "get version")
			}
			logrus.Infof("import %v: %v %v", deviceTypeName, version, revision)
			if params, ok := schemaParamsMap[deviceTypeName]; ok && params != nil {
				for _, v := range importParams {
					if v2 := params.GetParameterDefine(v.PrivOid); v2 != nil && v2.Equal(v) {
						v2.AddVersion(info.Schema, version)
					} else {
						newParam := v
						newParam.AddVersion(info.Schema, version)
						params.AddParameterDefine(newParam)
					}
				}
			} else {
				for _, v := range importParams {
					v.AddVersion(info.Schema, version)
				}
				schemaParamsMap[deviceTypeName] = importParams
			}
		}
	}

	os.MkdirAll(outputPath, os.ModePerm)

	for schema, schemaParamsMap := range all {
		for deviceTypeName, params := range schemaParamsMap {
			productDef := model.GetProductDefine(schema, deviceTypeName)
			if productDef == nil {
				return errors.Errorf("invalid device type, schema=%v, device=%v", schema, deviceTypeName)
			}
			versions := productDef.Versions
			if n := len(versions); n > 0 {
				lastVersion := versions[n-1]
				for _, param := range params {
					if param.Versions.MatchVersion(schema, lastVersion) {
						param.AddVersion(schema, "latest")
					}
				}
				versions = append(versions, "latest")
			}

			for _, version := range versions {
				params2 := params.GetParameterDefinesByVersion(schema, version)
				if productDef.SupportLayout {
					logrus.Infof("setup application, schema=%v, device=%v, version=%v, params=%v", schema, deviceTypeName, version, len(params2))
					app, err := layout.MakeApplication(deviceTypeName, layout.DeviceInfo{
						Device: productDef,
						Params: params2,
						Alarms: alarmDefs,
					})
					if err != nil {
						return errors.Wrapf(err, "make application layout, schema=%v, device=%v, version=%v", schema, deviceTypeName, version)
					}
					if err := saveJsonFile(path.Join(outputPath, "mock", schema, deviceTypeName, version, "layout.json"), app); err != nil {
						return errors.Wrap(err, "save layout.json")
					}
					if productDef.Schema == "corning" && version != "demo" && version != "latest" {
						storyList := story.MakeApplictionStoryList(app, productDef, params2)
						if err := story.UpdateStoryListToJira(productDef, version, storyList); err != nil {
							logrus.Error(errors.Wrap(err, "update story to jira"))
						}
					}
				}

				if err := saveJsonFile(path.Join(outputPath, "mock", schema, deviceTypeName, version, "parameters.json"), params2); err != nil {
					return errors.Wrap(err, "save parameters.json")
				}
			}

			ignoreParamOids := []string{"TB2.P0CCC", "TB2-PA.P0060"}
			for _, param := range params {
				ignore := false
				for _, oid := range ignoreParamOids {
					if string(param.PrivOid) == oid {
						ignore = true
						break
					}
				}
				if ignore {
					continue
				}
				if productDef.SupportLayout {
					if len(param.Groups) >= 1 &&
						(strings.HasPrefix(param.Groups[0], "Settings,Engineering CMD") ||
							strings.HasPrefix(param.Groups[0], "Settings,Engineering CMD1") ||
							strings.HasPrefix(param.Groups[0], "Settings,Engineering CMD2")) {
						//ignore
					} else if param.Paths == nil || len(param.Paths) == 0 {
						logrus.Warnf("not used parameter: %v, %v, %v", param.PrivOid, param.Name, param.Paths)
					}
				}
			}
			os.MkdirAll(filepath.Join(outputPath, "models", schema, deviceTypeName), os.ModePerm)

			if err := saveParametersFile(path.Join(outputPath, "models", schema, deviceTypeName, "parameters.csv"), params); err != nil {
				return errors.Wrap(err, "save parameters.csv")
			}
			if err := saveJsonFile(path.Join(outputPath, "models", schema, deviceTypeName, "parameters.json"), params); err != nil {
				return errors.Wrap(err, "save parameters.json")
			}
		}
	}
	return nil
}

func loadInfo(infoFile string) (*parser.ReleaseInfo, error) {
	buf, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, errors.Wrap(err, "read info file")
	}
	info := &parser.ReleaseInfo{}
	if err := yaml.Unmarshal(buf, info); err != nil {
		return nil, errors.Wrap(err, "unmarshal info file")
	}
	return info, nil
}

func parseDeviceFiles(schema string, xmlFiles []string, iniFiles []string) (map[string]model.ParameterDefines, error) {
	all := map[string]model.ParameterDefines{}

	deviceParsers, err := parser.NewDeviceFileParsers(xmlFiles, iniFiles)
	if err != nil {
		return nil, errors.Wrap(err, "load device files")
	}

	for _, deviceParser := range deviceParsers {
		deviceTypeName := deviceParser.GetDeviceTypename()

		productDef := model.GetProductDefine(schema, deviceTypeName)
		if productDef == nil {
			return nil, errors.Errorf("invalid device type for %v:%v", schema, deviceTypeName)
		}
		productModel := model.GetProductModelByDeviceTypeName(deviceTypeName)
		if productModel == nil {
			return nil, errors.Errorf("invalid product model for %v:%v", schema, deviceTypeName)
		}

		logrus.Infof("loading %v:%v, %v",
			schema, deviceTypeName,
			productModel.String())
		params := deviceParser.GetParameters()
		params.FixParameters()
		all[deviceTypeName] = params
	}

	return all, nil
}

func saveParametersFile(filename string, params model.ParameterDefines) error {
	// Create a new CSV file
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "create file")
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	header := []string{"Groups", "Paths", "OID", "Name", "Unit", "Tips", "AccessMode", "DataType", "ByteSize", "Ratio", "Range", "Options"}
	if err := csvWriter.Write(header); err != nil {
		return errors.Wrap(err, "write header")
	}

	data := [][]string{}
	for _, param := range params {
		ratioString := ""
		if param.Ratio != nil {
			ratioString = fmt.Sprintf("%v", *param.Ratio)
		}
		rangeString := ""
		if param.Max != nil && param.Min != nil {
			rangeString = fmt.Sprintf("%v ~ %v", *param.Min, *param.Max)
		} else if param.Max != nil {
			rangeString = fmt.Sprintf(" ~ %v", *param.Max)
		} else if param.Min != nil {
			rangeString = fmt.Sprintf("%v ~ ", *param.Min)
		}
		optionsStrings := []string{}
		for k, v := range param.Options {
			optionsStrings = append(optionsStrings, fmt.Sprintf("%v=%v", k, v))
		}
		sort.Slice(optionsStrings, func(i, j int) bool {
			return optionsStrings[i] < optionsStrings[j]
		})
		item := []string{
			strings.Join(param.Groups, "\r\n"),
			strings.Join(param.Paths, "\r\n"),
			string(param.PrivOid),
			param.Name,
			param.UnitName,
			param.Tips,
			string(param.Access),
			string(param.DataType),
			fmt.Sprintf("%v", param.ByteSize),
			ratioString,
			rangeString,
			strings.Join(optionsStrings, ";"),
		}
		data = append(data, item)
	}
	for _, row := range data {
		if err := csvWriter.Write(row); err != nil {
			return errors.Wrap(err, "write data row")
		}
	}

	// Flush the CSV writer to ensure all data is written to the file
	csvWriter.Flush()

	return nil
}

func saveJsonFile(filename string, v any) error {
	content, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err := os.WriteFile(filename, content, 0666); err != nil {
		return errors.Wrap(err, "write file")
	}

	gzipFile, _ := os.Create(filename + ".gz")
	gzipWriter := gzip.NewWriter(gzipFile)
	defer gzipWriter.Close()
	if _, err := gzipWriter.Write(content); err != nil {
		return errors.Wrap(err, "write gzip file")
	}
	return nil
}
