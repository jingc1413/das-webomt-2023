package model

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"inet.af/netaddr"
)

var defaultTimeLocation *time.Location = time.UTC

func SetTimeLocation(name string) {
	if strings.HasPrefix(name, "UTC") {
		v := cast.ToInt(name[3:]) * 3600
		logrus.Tracef("set time locaion, %v, %v", name, v)
		defaultTimeLocation = time.FixedZone(name, v)
	}
}

func (m ParameterDefine) UnmarshalBinaryValue(data []byte) (any, error) {
	r := bytes.NewReader(data)
	switch m.DataType {
	case DataTypeBinary:
		return hex.EncodeToString(data), nil
	case DataTypeString:
		end := len(data)
		for i, v := range data {
			if v == 0 {
				end = i
				break
			}
		}
		value := cast.ToString(data[:end])
		return value, nil
	case DataTypeBool:
		var v uint8
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return cast.ToUint8(v), nil
	case DataTypeUInt8:
		var v uint8
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeUInt16:
		var v uint16
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeUInt32:
		var v uint32
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeUInt64:
		var v uint64
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeInt8:
		var v int8
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeInt16:
		var v int16
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeInt32:
		var v int32
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeInt64:
		var v int64
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeFloat32:
		var v float32
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeFloat64:
		var v float64
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeTimestamp:
		var v int64
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, errors.Wrap(err, "read binary")
		}
		return v, nil
	case DataTypeIPV4:
		if len(data) != 4 {
			return "", errors.New("invalid value for ipv4")
		}
		v := [4]byte{}
		copy(v[0:3], data[0:3])
		ip := netaddr.IPFrom4(v)
		if !ip.Is4() {
			return "", errors.New("invalid value for ipv4")
		}
		return ip.String(), nil
	case DataTypeIPV6:
		if len(data) != 16 {
			return "", errors.New("invalid value for ipv6")
		}
		v := [16]byte{}
		copy(v[0:16], data[0:16])
		ip := netaddr.IPFrom16(v)
		if !ip.Is4() {
			return "", errors.New("invalid value for ipv6")
		}
		return ip.String(), nil
	case DataTypeDateTime:
		v := fmt.Sprintf("%02d%02d-%02d-%02d %02d:%02d:%02d",
			data[0], data[1], data[2], data[3], data[4], data[5], data[6])
		t, err := time.ParseInLocation(time.DateTime, v, defaultTimeLocation)
		if err != nil {
			return nil, errors.New("invalid value for datetime")
		}
		return t.Unix(), nil
	case DataTypeObject:
		v := []any{}
		offset := 0
		for _, child := range m.Child {
			v2, err := child.UnmarshalBinaryValue(data[offset : offset+int(child.ByteSize)])
			if err != nil {
				return nil, errors.Wrapf(err, "parse child")
			}
			v = append(v, v2)
			offset = offset + int(child.ByteSize)
		}
		return v, nil
	}
	if strings.HasPrefix(string(m.DataType), string(DataTypeStringArray)) {
		v := strings.Split(cast.ToString(data), m.Split)
		if len(v) != len(m.Child) {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		return v, nil
	} else if strings.HasPrefix(string(m.DataType), string(DataTypeInt32Array)) {
		v := []int32{}
		for {
			var v2 int32
			if err := binary.Read(r, binary.LittleEndian, &v2); err != nil {
				break
			}
			v = append(v, v2)
		}
		if len(v) != len(m.Child) {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		return v, nil
	} else if strings.HasPrefix(string(m.DataType), string(DataTypeUint32Array)) {
		v := []uint32{}
		for {
			var v2 uint32
			if err := binary.Read(r, binary.LittleEndian, &v2); err != nil {
				break
			}
			v = append(v, v2)
		}
		if len(v) != len(m.Child) {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		return v, nil
	}

	return nil, errors.Errorf("invalid data type %v", m.DataType)
}

func (m ParameterDefine) MarshalBinaryValue(value any) ([]byte, error) {
	w := new(bytes.Buffer)
	switch m.DataType {
	case DataTypeBinary:
		if m.MultipleOption {
			values, ok := value.([]any)
			if !ok {
				return nil, errors.Errorf("invalid value for %v", m.DataType)
			}
			sum := uint64(0)
			for _, v := range values {
				n := cast.ToUint64(v)
				sum += n
			}
			format := fmt.Sprintf("%%0%dX", m.ByteSize*2)
			value = fmt.Sprintf(format, sum)
		}
		v, err := hex.DecodeString(cast.ToString(value))
		if err != nil {
			return nil, errors.Wrap(err, "decode hex string")
		}
		buf := make([]byte, m.ByteSize)
		size := len(v)
		if size > int(m.ByteSize) {
			size = int(m.ByteSize)
		}
		copy(buf[0:size], v[0:size])
		if err := binary.Write(w, binary.LittleEndian, buf); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeString:
		v := []byte(cast.ToString(value))
		buf := make([]byte, m.ByteSize)
		size := len(v)
		if size > int(m.ByteSize) {
			size = int(m.ByteSize)
		}
		copy(buf[0:size], v[0:size])
		if err := binary.Write(w, binary.LittleEndian, buf); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeBool:
		v := cast.ToUint8(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeUInt8:
		v := cast.ToUint8(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeUInt16:
		v := cast.ToUint16(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeUInt32:
		v := cast.ToUint32(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeUInt64:
		v := cast.ToUint64(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeInt8:
		v := cast.ToInt8(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeInt16:
		v := cast.ToInt16(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeInt32:
		v := cast.ToInt32(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeInt64:
		v := cast.ToInt64(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeFloat32:
		v := cast.ToFloat32(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeFloat64:
		v := cast.ToFloat64(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeTimestamp:
		v := cast.ToInt64(value)
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeIPV4:
		ip, err := netaddr.ParseIP(cast.ToString(value))
		if err != nil {
			return nil, errors.Wrap(err, "parse ip")
		}
		if !ip.Is4() {
			return nil, errors.New("invalid value for ipv4")
		}
		buf, err := ip.MarshalBinary()
		if err != nil {
			return nil, errors.Wrap(err, "marshal ip")
		}
		if err := binary.Write(w, binary.LittleEndian, buf); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeIPV6:
		ip, err := netaddr.ParseIP(cast.ToString(value))
		if err != nil {
			return nil, errors.Wrap(err, "parse ip")
		}
		if !ip.Is6() {
			return nil, errors.New("invalid value for ipv6")
		}
		buf, err := ip.MarshalBinary()
		if err != nil {
			return nil, errors.Wrap(err, "marshal ip")
		}
		if err := binary.Write(w, binary.LittleEndian, buf); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeDateTime:
		t := time.Unix(cast.ToInt64(value), 0).In(defaultTimeLocation)
		buf := make([]byte, m.ByteSize)
		buf[0] = byte(t.Year() / 100)
		buf[1] = byte(t.Year() % 100)
		buf[2] = byte(t.Month())
		buf[3] = byte(t.Day())
		buf[4] = byte(t.Hour())
		buf[5] = byte(t.Minute())
		buf[6] = byte(t.Second())
		if err := binary.Write(w, binary.LittleEndian, buf); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	case DataTypeObject:
		values, ok := value.([]any)
		if !ok {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		if len(values) != len(m.Child) {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		for i, child := range m.Child {
			v2, err := child.MarshalBinaryValue(values[i])
			if err != nil {
				return nil, errors.Wrapf(err, "marshal child")
			}
			if _, err := w.Write(v2); err != nil {
				return nil, errors.Wrap(err, "write binary")
			}
		}
	}

	if strings.HasPrefix(string(m.DataType), string(DataTypeStringArray)) {
		values, ok := value.([]string)
		if !ok {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		v := []byte(strings.Join(values, m.Split))
		size := len(v)
		if size > int(m.ByteSize) {
			size = int(m.ByteSize)
		}
		buf := make([]byte, m.ByteSize)
		copy(buf[0:size], v[0:size])
		if err := binary.Write(w, binary.LittleEndian, buf); err != nil {
			return nil, errors.Wrap(err, "write binary")
		}
	} else if strings.HasPrefix(string(m.DataType), string(DataTypeInt32Array)) {
		values, ok := value.([]int32)
		if !ok {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		for _, v := range values {
			if err := binary.Write(w, binary.LittleEndian, v); err != nil {
				return nil, errors.Wrap(err, "write binary")
			}
		}
	} else if strings.HasPrefix(string(m.DataType), string(DataTypeUint32Array)) {
		values, ok := value.([]uint32)
		if !ok {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		for _, v := range values {
			if err := binary.Write(w, binary.LittleEndian, v); err != nil {
				return nil, errors.Wrap(err, "write binary")
			}
		}
	}

	out := w.Bytes()
	if len(out) == 0 {
		return nil, errors.Errorf("invalid data type %v", m.DataType)
	}
	if len(out) != int(m.ByteSize) {
		logrus.Warnf("%v: %v", len(out), out)
		return nil, errors.New("incorrect byte size")
	}
	return out, nil
}

func (m ParameterDefine) UnmarshalCgiStringValue(data string) (any, error) {
	if data == "invalid" {
		return nil, nil
	}
	if data == "module not installed" {
		return nil, nil
	}

	var number int64
	if m.DataType.IsNumeric() {
		d, err := decimal.NewFromString(data)
		if err != nil {
			return nil, errors.Wrapf(err, "parse decimal value %v", data)
		}
		if m.Ratio != nil && *m.Ratio != 1 {
			d = d.Mul(decimal.NewFromInt(*m.Ratio))
		}
		number = d.IntPart()
	}
	switch m.DataType {
	case DataTypeBinary:
		if m.MultipleOption {
			n, err := strconv.ParseUint(data, 16, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid value %v", data)
			}
			out := []string{}
			format := fmt.Sprintf("%%0%dX", m.ByteSize*2)
			for k, _ := range m.Options {
				n2, err := strconv.ParseUint(k, 16, 64)
				if err != nil {
					return nil, errors.Wrapf(err, "invalid option %v", k)
				}
				if n&n2 != 0 {
					out = append(out, fmt.Sprintf(format, n2))
				}
			}
			return out, nil
		} else if len(m.Options) > 0 {
			for k, v := range m.Options {
				if v == data {
					return k, nil
				}
			}
		}
		return cast.ToString(data), nil
	case DataTypeString:
		return cast.ToString(data), nil
	case DataTypeBool:
		return cast.ToUint8(data), nil
	case DataTypeUInt8:
		return cast.ToUint8(number), nil
	case DataTypeUInt16:
		return cast.ToUint16(number), nil
	case DataTypeUInt32:
		return cast.ToUint32(number), nil
	case DataTypeUInt64:
		return cast.ToUint64(number), nil
	case DataTypeInt8:
		return cast.ToInt8(number), nil
	case DataTypeInt16:
		return cast.ToInt16(number), nil
	case DataTypeInt32:
		return cast.ToInt32(number), nil
	case DataTypeInt64:
		return cast.ToInt64(number), nil
	case DataTypeFloat32:
		return cast.ToFloat32(number), nil
	case DataTypeFloat64:
		return cast.ToFloat64(number), nil
	case DataTypeTimestamp:
		return cast.ToInt64(number), nil
	case DataTypeIPV4:
		v, err := netaddr.ParseIP(data)
		if err != nil {
			return nil, errors.Wrap(err, "invalid value for ipv4")
		}
		return v.String(), nil
	case DataTypeIPV6:
		v, err := netaddr.ParseIP(data)
		if err != nil {
			return nil, errors.Wrap(err, "invalid value for ipv4")
		}
		return v.String(), nil
	case DataTypeDateTime:
		t, err := time.ParseInLocation(time.DateTime, data, defaultTimeLocation)
		if err != nil {
			return nil, errors.Wrap(err, "invalid value for datetime")
		}
		return t.Unix(), nil
	case DataTypeObject:
		if m.Split == "" {
			buf, err := hex.DecodeString(data)
			if err != nil {
				return nil, errors.New("invalid value for object")
			}
			return m.UnmarshalBinaryValue(buf)
		} else {
			args := strings.Split(data, m.Split)
			if len(args) != len(m.Child) {
				return nil, errors.New("invalid value for object")
			}
			v := []any{}
			for i, child := range m.Child {
				v2, err := child.UnmarshalCgiStringValue(args[i])
				if err != nil {
					return nil, errors.Wrapf(err, "unmarshal child value")
				}
				v = append(v, v2)
			}
			return v, nil
		}

	}
	if strings.HasPrefix(string(m.DataType), string(DataTypeStringArray)) {
		v := strings.Split(data, m.Split)
		if len(v) != len(m.Child) {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		return v, nil
	} else if strings.HasPrefix(string(m.DataType), string(DataTypeInt32Array)) {
		args := strings.Split(data, m.Split)
		v := []int32{}
		for _, arg := range args {
			d, err := decimal.NewFromString(arg)
			if err != nil {
				return nil, errors.Wrapf(err, "parse decimal value %v", data)
			}
			if m.Ratio != nil && *m.Ratio != 1 {
				d = d.Mul(decimal.NewFromInt(*m.Ratio))
			}
			number := d.IntPart()
			v2 := cast.ToInt32(number)
			v = append(v, v2)
		}
		if len(v) != len(m.Child) {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		return v, nil
	} else if strings.HasPrefix(string(m.DataType), string(DataTypeUint32Array)) {
		args := strings.Split(data, m.Split)
		v := []uint32{}
		for _, arg := range args {
			d, err := decimal.NewFromString(arg)
			if err != nil {
				return nil, errors.Wrapf(err, "parse decimal value %v", data)
			}
			if m.Ratio != nil && *m.Ratio != 1 {
				d = d.Mul(decimal.NewFromInt(*m.Ratio))
			}
			number := d.IntPart()
			v2 := cast.ToUint32(number)
			v = append(v, v2)
		}
		if len(v) != len(m.Child) {
			return nil, errors.Errorf("invalid value for %v", m.DataType)
		}
		return v, nil
	}

	return nil, errors.Errorf("invalid data type %v", m.DataType)
}

func (m ParameterDefine) MarshalCgiStringValue(value any) (string, error) {
	if value == nil {
		return "", nil
	}

	var numberString string
	if m.DataType.IsNumeric() {
		d, err := decimal.NewFromString(fmt.Sprintf("%v", value))
		if err != nil {
			return "", errors.Wrapf(err, "parse decimal value %v", value)
		}
		if m.Ratio != nil && *m.Ratio != 1 {
			d = d.Div(decimal.NewFromInt(*m.Ratio))
		}
		if d.IsInteger() {
			numberString = fmt.Sprintf("%d", d.IntPart())
		} else {
			numberString = d.String()
		}
	}
	switch m.DataType {
	case DataTypeBinary:
		if m.MultipleOption {
			values, ok := value.([]any)
			if !ok {
				return "", errors.Errorf("invalid value for %v", m.DataType)
			}
			sum := uint64(0)
			for _, v := range values {
				n := cast.ToUint64(v)
				sum += n
			}
			format := fmt.Sprintf("%%0%dX", m.ByteSize*2)
			return fmt.Sprintf(format, sum), nil
		}
		v := cast.ToString(value)
		if v == "" {
			v = hex.EncodeToString(make([]byte, m.ByteSize))
		}
		return v, nil
	case DataTypeString:
		return cast.ToString(value), nil
	case DataTypeBool:
		return cast.ToString(value), nil
	case DataTypeUInt8:
		return cast.ToString(numberString), nil
	case DataTypeUInt16:
		return cast.ToString(numberString), nil
	case DataTypeUInt32:
		return cast.ToString(numberString), nil
	case DataTypeUInt64:
		return cast.ToString(numberString), nil
	case DataTypeInt8:
		return cast.ToString(numberString), nil
	case DataTypeInt16:
		return cast.ToString(numberString), nil
	case DataTypeInt32:
		return cast.ToString(numberString), nil
	case DataTypeInt64:
		return cast.ToString(numberString), nil
	case DataTypeFloat32:
		return cast.ToString(numberString), nil
	case DataTypeFloat64:
		return cast.ToString(numberString), nil
	case DataTypeTimestamp:
		return cast.ToString(value), nil
	case DataTypeIPV4:
		return cast.ToString(value), nil
	case DataTypeIPV6:
		return cast.ToString(value), nil
	case DataTypeDateTime:
		t := time.Unix(cast.ToInt64(value), 0).In(defaultTimeLocation)
		return t.Format(time.DateTime), nil
	case DataTypeObject:
		if m.Split == "" {
			buf, err := m.MarshalBinaryValue(value)
			if err != nil {
				return "", errors.Wrapf(err, "invalid value for %v", m.DataType)
			}
			return hex.EncodeToString(buf), nil
		} else {
			values, ok := value.([]any)
			if !ok {
				return "", errors.Errorf("invalid value for %v", m.DataType)
			}
			if len(values) != len(m.Child) {
				return "", errors.Errorf("invalid value for %v", m.DataType)
			}
			args := []string{}
			for i, v := range values {
				arg, err := m.Child[i].MarshalCgiStringValue(v)
				if err != nil {
					return "", errors.Wrapf(err, "marshal child value")
				}
				args = append(args, arg)
			}
			return strings.Join(args, m.Split), nil
		}

	}

	if strings.HasPrefix(string(m.DataType), string(DataTypeStringArray)) {
		values, ok := value.([]string)
		if !ok {
			return "", errors.Errorf("invalid value for %v", m.DataType)
		}
		return strings.Join(values, m.Split), nil
	} else if strings.HasPrefix(string(m.DataType), string(DataTypeInt32Array)) {
		values, ok := value.([]int32)
		if !ok {
			return "", errors.Errorf("invalid value for %v", m.DataType)
		}
		args := []string{}
		for _, v := range values {
			d, err := decimal.NewFromString(fmt.Sprintf("%v", v))
			if err != nil {
				return "", errors.Wrapf(err, "parse decimal value %v", v)
			}
			if m.Ratio != nil && *m.Ratio != 1 {
				d = d.Div(decimal.NewFromInt(*m.Ratio))
			}
			number := d.IntPart()
			args = append(args, cast.ToString(number))
		}
		return strings.Join(args, m.Split), nil
	} else if strings.HasPrefix(string(m.DataType), string(DataTypeUint32Array)) {
		values, ok := value.([]uint32)
		if !ok {
			return "", errors.Errorf("invalid value for %v", m.DataType)
		}
		args := []string{}
		for _, v := range values {
			d, err := decimal.NewFromString(fmt.Sprintf("%v", v))
			if err != nil {
				return "", errors.Wrapf(err, "parse decimal value %v", v)
			}
			if m.Ratio != nil && *m.Ratio != 1 {
				d = d.Div(decimal.NewFromInt(*m.Ratio))
			}
			number := d.IntPart()
			args = append(args, cast.ToString(number))
		}
		return strings.Join(args, m.Split), nil
	}
	return "", errors.Errorf("invalid data type %v", m.DataType)
}
