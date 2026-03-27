package xml

import "github.com/pkg/errors"

func LoadXmlFile(filename string) (*MachineRoot, error) {
	root := MachineRoot{}
	if err := root.UnmarshalXmlFile(filename, true); err != nil {
		return nil, errors.Wrap(err, "unmarshal device xml file")
	}
	root.CleanUp()
	return &root, nil
}
