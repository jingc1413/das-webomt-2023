package priv

import (
	"encoding/hex"
	"fmt"
	"gomt/core/model"

	"github.com/pkg/errors"
)

type Parameter struct {
	ID    string `json:",omitempty"`
	Value any    `json:",omitempty"`
	Code  string `json:",omitempty"`
}

func (s *PrivMgmtServer) RegisterDevice(
	apType ApType, vpType VpType, mcpType McpType,
	deviceTypeName string, version string, params model.ParameterDefines,
	mainId []byte, subId uint8, addr string) (*DeviceSession, error) {
	if params == nil {
		return nil, errors.Errorf("invalid registry device type name")
	}
	repo := GetDefaultDeviceSessionRepository()
	sess := repo.RegiterDeviceSession("udp", apType, vpType, mcpType, deviceTypeName, mainId, subId)
	sess.Addr = addr
	sess.Params = params
	s.log.Tracef("register device: %v, %v, %v, %v", deviceTypeName, hex.EncodeToString(mainId), subId, addr)
	return sess, nil
}

func (s *PrivMgmtServer) QueryDevices(
	sess *DeviceSession,
	schema string,
	queryString string,
	cb func(total uint8, index uint8, info *model.DeviceInfo),
) ([]*model.DeviceInfo, error) {
	if sess == nil {
		return nil, errors.New("invalid device session")
	}

	output, err := s.queryDevices(sess, schema, queryString, cb)
	if err != nil {
		return nil, errors.Wrap(err, "query devices")
	}
	return output, nil
}

// func (s *PrivMgmtServer) getSession(mainId []byte, subId uint8) (*DeviceSession, error) {
// 	repo := GetDefaultDeviceSessionRepository()
// 	sess := repo.GetDeviceSession(mainId, subId)
// 	if sess == nil {
// 		return nil, errors.New("invalid device")
// 	}
// 	return sess, nil
// }

func (s *PrivMgmtServer) GetParameterValues(sess *DeviceSession, inputs []Parameter) ([]Parameter, error) {
	if sess == nil {
		return nil, errors.New("invalid device session")
	}
	return s.getParameterValues(sess, sess.Params, inputs)
}

func (s *PrivMgmtServer) SetParameterValues(sess *DeviceSession, inputs []Parameter) ([]Parameter, error) {
	if sess == nil {
		return nil, errors.New("invalid device session")
	}
	return s.setParameterValues(sess, sess.Params, inputs)
}

func (s *PrivMgmtServer) QueryObjectValues(sess *DeviceSession, ids []string) ([]model.Object, error) {
	if sess == nil {
		return nil, errors.New("invalid device session")
	}
	return s.queryObjectValues(sess, sess.Params, ids)
}

func (s *PrivMgmtServer) SetObjectValues(sess *DeviceSession, inputs []model.Object) ([]model.Object, error) {
	if sess == nil {
		return nil, errors.New("invalid device session")
	}
	return s.setObjectValues(sess, sess.Params, inputs)
}

func (s *PrivMgmtServer) getParameterValues(sess *DeviceSession, params model.ParameterDefines, inputs []Parameter) ([]Parameter, error) {
	objsMap := map[string][]Object{}

	for _, v := range inputs {
		id, err := model.PrivObjectIdFromString(v.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid parameter: id=%v", id)
		}
		param := params.GetParameterDefine(id)
		if param == nil {
			return nil, errors.Wrapf(err, "invalid parameter: id=%v", id)
		}
		cid, _, oid := param.PrivOid.Values()
		objs, ok := objsMap[cid]
		if objs == nil || !ok {
			objs = []Object{}
		}
		objs = append(objs, Object{OID: ObjectID(oid), Value: v.Value})
		objsMap[cid] = objs
	}
	outputs := []Parameter{}
	for cmdId, objs := range objsMap {
		ret, err := s.queryObjects(sess, cmdId, objs)
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			_id, _ := model.MakePrivObjectId(cmdId, v.OID.String())
			id := string(_id)
			value := v.Value
			outputs = append(outputs, Parameter{ID: id, Value: value, Code: fmt.Sprintf("%02X", v.Fault)})
		}
	}

	return outputs, nil
}

func (s *PrivMgmtServer) setParameterValues(sess *DeviceSession, params model.ParameterDefines, inputs []Parameter) ([]Parameter, error) {
	objsMap := map[string][]Object{}

	for _, v := range inputs {
		id, err := model.PrivObjectIdFromString(v.ID)
		if err != nil {
			return nil, errors.Errorf("invalid parameter: id=%v", v.ID)
		}
		param := params.GetParameterDefine(id)
		if param == nil {
			return nil, errors.Errorf("invalid parameter: id=%v", v.ID)
		}
		cid, _, oid := param.PrivOid.Values()
		objs, ok := objsMap[cid]
		if objs == nil || !ok {
			objs = []Object{}
		}
		objs = append(objs, Object{OID: ObjectID(oid), Value: v.Value})
		objsMap[cid] = objs
	}

	outputs := []Parameter{}
	for cmdId, objs := range objsMap {
		ret, err := s.setObjects(sess, cmdId, objs)
		if err != nil {
			return outputs, err
		}
		for _, v := range ret {
			_id, _ := model.MakePrivObjectId(cmdId, v.OID.String())
			id := string(_id)
			value := v.Value
			outputs = append(outputs, Parameter{ID: id, Value: value, Code: fmt.Sprintf("%02X", v.Fault)})
		}
	}
	return outputs, nil
}

func (s *PrivMgmtServer) queryObjectValues(sess *DeviceSession, params model.ParameterDefines, ids []string) ([]model.Object, error) {
	objsMap := map[string][]Object{}

	for _, id := range ids {
		_id, err := model.PrivObjectIdFromString(id)
		if err != nil {
			return nil, errors.Errorf("invalid parameter: id=%v", id)
		}
		param := params.GetParameterDefine(_id)
		if param == nil {
			return nil, errors.Errorf("invalid parameter: id=%v", id)
		}
		cid, _, oid := param.PrivOid.Values()
		objs, ok := objsMap[cid]
		if objs == nil || !ok {
			objs = []Object{}
		}
		objs = append(objs, Object{OID: ObjectID(oid)})
		objsMap[cid] = objs
	}
	outputs := []model.Object{}
	for cmdId, objs := range objsMap {
		ret, err := s.queryObjects(sess, cmdId, objs)
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			_id, _ := model.MakePrivObjectId(cmdId, v.OID.String())
			outputs = append(outputs, model.Object{
				ID:    _id,
				Value: v.Value,
				Code:  int(v.Fault),
			})
		}
	}
	return outputs, nil
}

func (s *PrivMgmtServer) setObjectValues(sess *DeviceSession, params model.ParameterDefines, inputs []model.Object) ([]model.Object, error) {
	objsMap := map[string][]Object{}

	for _, v := range inputs {
		param := params.GetParameterDefine(v.ID)
		if param == nil {
			return nil, errors.Errorf("invalid parameter: id=%v", v.ID)
		}
		cid, _, oid := param.PrivOid.Values()
		objs, ok := objsMap[cid]
		if objs == nil || !ok {
			objs = []Object{}
		}
		objs = append(objs, Object{OID: ObjectID(oid), Value: v.Value})
		objsMap[cid] = objs
	}
	outputs := []model.Object{}
	for cmdId, objs := range objsMap {
		ret, err := s.setObjects(sess, cmdId, objs)
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			_id, _ := model.MakePrivObjectId(cmdId, v.OID.String())
			outputs = append(outputs, model.Object{
				ID:    _id,
				Value: v.Value,
				Code:  int(v.Fault),
			})
		}
	}

	return outputs, nil
}
