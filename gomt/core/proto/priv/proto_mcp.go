package priv

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gomt/core/model"

	"github.com/pkg/errors"
)

type McpData []byte

var EMPTY_MCP_DATA McpData = McpData{}

type ObjectID uint32
type ObjectIDSort []ObjectID

func (m ObjectIDSort) Len() int           { return len(m) }
func (m ObjectIDSort) Less(i, j int) bool { return uint32(m[i]) < uint32(m[j]) }
func (m ObjectIDSort) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

const (
	OID_NOTIFICATION         ObjectID = 0x00000141
	OID_QUERY_ALL_OBJECT_IDS ObjectID = 0x00000009
	OID_QUERY_ALL_DEVICES    ObjectID = 0x00000AEC
)

func (m ObjectID) String() string { return fmt.Sprintf("%04X", uint32(m)) }
func (m *ObjectID) UnmarshalString(in string) error {
	bytes, err := hex.DecodeString(in)
	if err != nil {
		return err
	}
	len := len(bytes)
	if len < 1 || len > 4 {
		return errors.New("invalid length")
	}
	buf := [4]byte{}
	copy(buf[4-len:], bytes)
	v := uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	*m = ObjectID(v)
	return nil
}

type Object struct {
	OID   ObjectID
	Value any
	Size  int64
	Fault uint8
}

func (m Object) String() string {
	return fmt.Sprintf("{%04X=%v}", uint32(m.OID), m.Value)
}

func (m ObjectID) MarshalBinaryWithType(typ McpType, fault uint8) ([]byte, error) {
	w := new(bytes.Buffer)
	if typ == MCP_A {
		if m > 0xFFFF {
			return w.Bytes(), errors.New("invaild object id of MCP_A")
		}
		oid := uint16(((uint32(fault) & 0xf) << 12) | uint32(m)&0xfff)
		if err := binary.Write(w, binary.LittleEndian, oid); err != nil {
			return w.Bytes(), err
		}
	} else if typ == MCP_B {
		return w.Bytes(), errors.New("invaild object id of MCP_B")
	} else if typ == MCP_C {
		oid := uint32(((uint32(fault) & 0xf) << 28) | uint32(m)&0xfffffff)
		if err := binary.Write(w, binary.LittleEndian, oid); err != nil {
			return w.Bytes(), err
		}
	}
	return w.Bytes(), nil
}

func (m *ObjectID) UnmarshalBinaryWithType(typ McpType, in []byte, fault *uint8) error {
	length := len(in)
	if typ == MCP_A {
		if length < 2 {
			return errors.New("invalid length of object id")
		}
		var oid uint16
		r := bytes.NewReader(in[0:2])
		if err := binary.Read(r, binary.LittleEndian, &oid); err != nil {
			return errors.Wrap(err, "read object id")
		}
		if fault != nil {
			*fault = uint8((oid >> 12) & 0xf)
		}
		oid = oid & 0xfff
		*m = ObjectID(oid)
	} else if typ == MCP_B {
		return errors.New("invaild object id of MCP_B")
	} else if typ == MCP_C {
		if length < 4 {
			return errors.New("invalid length of object id")
		}
		var oid uint32
		r := bytes.NewReader(in[0:4])
		if err := binary.Read(r, binary.LittleEndian, &oid); err != nil {
			return errors.Wrap(err, "read object id")
		}
		if fault != nil {
			*fault = uint8((oid >> 28) & 0xf)
		}
		oid = oid & 0xfffffff
		*m = ObjectID(oid)
	}
	return nil
}

func (m McpData) MarshalBinary() ([]byte, error) {
	out := make([]byte, len(m))
	copy(out[:], m[:])
	return out, nil
}

func (m *McpData) UnmarshalBinary(in []byte) error {
	*m = in[:]
	return nil
}

func (m *McpData) SetupQueryObjectIdsRequest(typ McpType, total uint8, index uint8) error {
	w := new(bytes.Buffer)
	var length uint8 = 0
	queryId := OID_QUERY_ALL_OBJECT_IDS
	b, err := queryId.MarshalBinaryWithType(typ, 0)
	if err != nil {
		return errors.Wrap(err, "marshal object id")
	}
	binary.Write(w, binary.LittleEndian, length)
	binary.Write(w, binary.LittleEndian, b)
	binary.Write(w, binary.LittleEndian, total)
	binary.Write(w, binary.LittleEndian, index)
	buf := w.Bytes()
	buf[0] = byte(len(buf))
	*m = buf
	return nil
}

func (m McpData) ParseQueryObjectIdsRequest(typ McpType) (uint8, uint8, error) {
	queryId := ObjectID(0)
	var (
		total uint8 = 0
		index uint8 = 0
	)
	offset := 0
	paramIdSize := 2
	if typ == MCP_C {
		paramIdSize = 4
	}
	size := len(m)
	if size < 1 {
		return 0, 0, ErrInvalidLength
	}
	length := int(m[0])
	if length != size || length < 1+2+paramIdSize {
		return 0, 0, ErrInvalidLength
	}
	offset += 1
	if err := queryId.UnmarshalBinaryWithType(typ, m[offset:], nil); err != nil {
		return 0, 0, errors.Wrap(err, "unmarshal object id")
	}
	offset += paramIdSize
	total = uint8(m[offset])
	offset += 1
	index = uint8(m[offset])
	offset += 1
	return total, index, nil
}

func (m *McpData) SetupQueryObjectIdsResponse(typ McpType, total uint8, index uint8, oids []ObjectID) error {
	w := new(bytes.Buffer)
	queryId := OID_QUERY_ALL_OBJECT_IDS
	b, err := queryId.MarshalBinaryWithType(typ, 0)
	if err != nil {
		return errors.Wrap(err, "marshal object id")
	}
	var (
		length uint8 = 0
	)
	binary.Write(w, binary.LittleEndian, length)
	binary.Write(w, binary.LittleEndian, b)
	binary.Write(w, binary.LittleEndian, total)
	binary.Write(w, binary.LittleEndian, index)

	for _, oid := range oids {
		b, err := oid.MarshalBinaryWithType(typ, 0)
		if err != nil {
			return errors.Wrap(err, "marshal object id")
		}
		binary.Write(w, binary.LittleEndian, b)
	}
	buf := w.Bytes()
	buf[0] = byte(len(buf))
	*m = buf
	return nil
}

func (m McpData) ParseQueryObjectIdsResponse(typ McpType) (uint8, uint8, []ObjectID, error) {
	queryId := ObjectID(0)
	oids := []ObjectID{}
	var (
		total uint8 = 0
		index uint8 = 0
	)
	offset := 0
	paramIdSize := 2
	if typ == MCP_C {
		paramIdSize = 4
	}
	size := len(m)
	if size < 1 {
		return 0, 0, oids, ErrInvalidLength
	}
	length := int(m[0])
	if length != size || (length-1-paramIdSize-2)%paramIdSize != 0 {
		return 0, 0, oids, ErrInvalidLength
	}
	offset += 1
	if err := queryId.UnmarshalBinaryWithType(typ, m[offset:], nil); err != nil {
		return 0, 0, oids, errors.Wrap(err, "unmarshal object id")
	}
	offset += paramIdSize
	total = uint8(m[offset])
	offset += 1
	index = uint8(m[offset])
	offset += 1

	for {
		if length-offset < paramIdSize {
			break
		}
		oid := ObjectID(0)
		if err := oid.UnmarshalBinaryWithType(typ, m[offset:], nil); err != nil {
			return total, index, oids, errors.Wrap(err, "unmarshal object id")
		}
		offset += paramIdSize
		oids = append(oids, oid)
	}

	return total, index, oids, nil
}

func (m *McpData) SetupQueryDevicesRequest(typ McpType, total uint8, index uint8, inputs []uint8) error {
	w := new(bytes.Buffer)
	var (
		length       uint8 = 0
		inputsLength uint8 = uint8(len(inputs))
	)
	queryId := OID_QUERY_ALL_DEVICES
	b, err := queryId.MarshalBinaryWithType(typ, 0)
	if err != nil {
		return errors.Wrap(err, "marshal object id")
	}
	binary.Write(w, binary.LittleEndian, length)
	binary.Write(w, binary.LittleEndian, b)
	binary.Write(w, binary.LittleEndian, total)
	binary.Write(w, binary.LittleEndian, index)
	binary.Write(w, binary.LittleEndian, inputsLength)
	for _, input := range inputs {
		binary.Write(w, binary.LittleEndian, input)
	}
	buf := w.Bytes()
	buf[0] = byte(len(buf))
	*m = buf
	return nil
}

func (m McpData) ParseQueryDevicesRequest(typ McpType) (uint8, uint8, []uint8, error) {
	queryId := ObjectID(0)
	var (
		total        uint8 = 0
		index        uint8 = 0
		inputsLength uint8 = 0
	)
	inputs := []uint8{}
	offset := 0
	paramIdSize := 2
	if typ == MCP_C {
		paramIdSize = 4
	}
	size := len(m)
	if size < 1 {
		return 0, 0, inputs, ErrInvalidLength
	}
	length := int(m[0])
	if length > size || length < 1+2+paramIdSize {
		return 0, 0, inputs, ErrInvalidLength
	}
	offset += 1
	if err := queryId.UnmarshalBinaryWithType(typ, m[offset:], nil); err != nil {
		return 0, 0, inputs, errors.Wrap(err, "unmarshal object id")
	}
	offset += paramIdSize
	total = uint8(m[offset])
	offset += 1
	index = uint8(m[offset])
	offset += 1
	inputsLength = uint8(m[offset])
	offset += 1
	for i := 0; i < int(inputsLength); i++ {
		inputs = append(inputs, m[offset+i])
	}
	offset += int(inputsLength)
	return total, index, inputs, nil
}

func GetQueyrDeviceSize(inputs []uint8) int {
	total := 0
	for _, input := range inputs {
		if def, ok := QueryDeviceInputDefineMap[input]; ok {
			total += int(def.Size)
		}
	}
	return total
}

func (m *McpData) SetupQueryDevicesResponse(typ McpType, total uint8, index uint8, inputs []uint8, outputs []QueryDevcieOutput) error {
	w := new(bytes.Buffer)
	queryId := OID_QUERY_ALL_DEVICES
	b, err := queryId.MarshalBinaryWithType(typ, 0)
	if err != nil {
		return errors.Wrap(err, "marshal object id")
	}
	var (
		length       uint8 = 0
		inputsLength uint8 = uint8(len(inputs))
	)

	binary.Write(w, binary.LittleEndian, length)
	binary.Write(w, binary.LittleEndian, b)
	binary.Write(w, binary.LittleEndian, total)
	binary.Write(w, binary.LittleEndian, index)
	binary.Write(w, binary.LittleEndian, inputsLength)
	for _, p := range inputs {
		binary.Write(w, binary.LittleEndian, p)
	}
	for _, output := range outputs {
		b, err := output.MarshalBinaryWithInputs(inputs)
		if err != nil {
			return errors.Wrap(err, "marshal output binary")
		}
		if err := binary.Write(w, binary.LittleEndian, b); err != nil {
			return errors.Wrap(err, "write output binary")
		}
	}
	buf := w.Bytes()
	buf[0] = byte(len(buf))
	*m = buf
	return nil
}

func (m McpData) ParseQueryDevicesResponse(typ McpType) (uint8, uint8, []uint8, []QueryDevcieOutput, error) {
	queryId := ObjectID(0)
	inputs := []uint8{}
	outputs := []QueryDevcieOutput{}
	var (
		total        uint8 = 0
		index        uint8 = 0
		inputsLength uint8 = 0
	)
	offset := 0
	paramIdSize := 2
	if typ == MCP_C {
		paramIdSize = 4
	}
	size := len(m)
	if size < 1 {
		return 0, 0, inputs, outputs, ErrInvalidLength
	}
	length := int(m[0])
	if length != size || length < 1+paramIdSize+2 {
		return 0, 0, inputs, outputs, ErrInvalidLength
	}
	offset += 1
	if err := queryId.UnmarshalBinaryWithType(typ, m[offset:], nil); err != nil {
		return 0, 0, inputs, outputs, errors.Wrap(err, "unmarshal object id")
	}
	offset += paramIdSize
	total = uint8(m[offset])
	offset += 1
	index = uint8(m[offset])
	offset += 1
	inputsLength = uint8(m[offset])
	offset += 1

	for i := 0; i < int(inputsLength); i++ {
		inputs = append(inputs, m[offset+i])
	}
	offset += int(inputsLength)

	deviceValueSize := 0
	for _, input := range inputs {
		def, ok := QueryDeviceInputDefineMap[input]
		if !ok {
			return 0, 0, inputs, outputs, errors.New("invalid input of query device")
		}
		deviceValueSize += int(def.Size)
	}

	for {
		if length-offset < deviceValueSize {
			break
		}
		output := QueryDevcieOutput{}
		if err := output.UnmarshalBinaryWithInputs(inputs, m[offset:offset+deviceValueSize]); err != nil {
			return 0, 0, inputs, outputs, errors.Wrap(err, "unmarshal output")
		}
		offset += deviceValueSize
		outputs = append(outputs, output)
	}

	return total, index, inputs, outputs, nil
}

func (m *McpData) SetupObjects(typ McpType, cmdId CmdID, objects []Object, defs model.ParameterDefines) error {
	w := new(bytes.Buffer)
	cid := cmdId & 0xFE
	for _, object := range objects {
		var length uint8 = 0
		oid, _ := model.MakePrivObjectId(fmt.Sprintf("%02X", cid), object.OID.String())
		def := defs.GetParameterDefine(oid)
		if def == nil {
			return errors.Errorf("unknow object id %v", oid)
		}
		b, err := object.OID.MarshalBinaryWithType(typ, object.Fault)
		if err != nil {
			return errors.Wrapf(err, "marshal object id %v", oid)
		}
		size := object.Size
		if size == 0 {
			size = def.ByteSize
		}
		length = uint8(len(b)) + uint8(size) + 1
		binary.Write(w, binary.LittleEndian, length)
		binary.Write(w, binary.LittleEndian, b)
		if object.Value == "" {
			empty := make([]byte, size)
			binary.Write(w, binary.LittleEndian, empty)
		} else {
			b2, err := def.MarshalBinaryValue(object.Value)
			if err != nil {
				return errors.Wrap(err, "marshal object value")
			}
			if len(b2) > int(size) {
				return errors.New(fmt.Sprintf("invalid object value size, %d>%d, %v, %s", len(b2), size, object.OID, object.Value))
			}
			if len(b2) < int(size) {
				add := make([]byte, int(size)-len(b2))
				for i := 0; i < len(add); i++ {
					add[i] = 0
				}
				b2 = append(b2, add...)
			}

			binary.Write(w, binary.LittleEndian, b2)
		}
	}
	*m = w.Bytes()
	return nil
}

func (m McpData) ParseObjects(typ McpType, cmdId CmdID, defs model.ParameterDefines) ([]Object, error) {
	cid := cmdId & 0xFE

	objects := []Object{}
	size := len(m)
	offset := 0

	paramIdSize := 2
	if typ == MCP_C {
		paramIdSize = 4
	}

	for {
		if size < offset+1+paramIdSize {
			break
		}
		length := int(m[offset])
		if length < 1+paramIdSize {
			return objects, errors.New("invalid object length")
		}
		object := Object{}
		if err := object.OID.UnmarshalBinaryWithType(typ, m[offset+1:], &object.Fault); err != nil {
			return objects, errors.Wrap(err, "unmarshal object id")
		}

		size := length - 1 - paramIdSize
		object.Size = int64(size)
		if object.OID == OID_QUERY_ALL_OBJECT_IDS {
			v := make([]byte, size)
			copy(v[:], m[offset+1+paramIdSize:offset+length])
			object.Value = hex.EncodeToString(v)
		} else if object.OID == OID_QUERY_ALL_DEVICES {
			v := make([]byte, size)
			copy(v[:], m[offset+1+paramIdSize:offset+length])
			object.Value = hex.EncodeToString(v)
		} else {
			oid, _ := model.MakePrivObjectId(fmt.Sprintf("%02X", cid), object.OID.String())
			def := defs.GetParameterDefine(oid)
			if def == nil {
				return objects, errors.New(fmt.Sprintf("unknow object id 0x%04X", object.OID))
			}
			v, err := def.UnmarshalBinaryValue([]byte(m[offset+1+paramIdSize : offset+length]))
			if err != nil {
				return objects, errors.Wrap(err, "unmarshal object value")
			}
			object.Value = v
		}
		offset += length
		objects = append(objects, object)
	}
	return objects, nil
}
