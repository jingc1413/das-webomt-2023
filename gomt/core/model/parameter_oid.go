package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type PrivObjectId string

func (m PrivObjectId) Values() (string, string, uint32) {
	re := regexp.MustCompile(`^(T([0-9A-F]{2})((-([A-Z0-9]+))|)\.|)P([0-9A-F]{4}|[0-9A-F]{8})$`)
	subs := re.FindSubmatch([]byte(m))
	if len(subs) != 7 {
		return "", "", 0
	}
	cid := string(subs[2])
	oid, _ := strconv.ParseUint(string(subs[6]), 16, 32)
	obj := string(subs[5])
	return cid, obj, uint32(oid)
}

func MakePrivObjectId(cid string, oid string) (PrivObjectId, error) {
	if strings.HasPrefix(oid, "P") || strings.HasPrefix(oid, "p") {
		oid = oid[1:]
	}
	return PrivObjectIdFromString(fmt.Sprintf("T%v.P%v", cid, oid))
}

func PrivObjectIdFromString(in string) (PrivObjectId, error) {
	re := regexp.MustCompile(`^(T([0-9A-F]{2})((-([A-Z0-9]+))|)\.|)P([0-9A-F]{4}|[0-9A-F]{8})$`)
	in = strings.ToUpper(strings.ReplaceAll(in, " ", ""))
	subs := re.FindSubmatch([]byte(in))
	if subs == nil {
		return "", errors.Errorf("invalid private object id %v", in)
	}
	var cid uint64 = 0x02
	var oid uint64 = 0
	var obj string = ""
	if len(subs[1]) != 0 {
		_cid, err := strconv.ParseUint(string(subs[2]), 16, 32)
		if err != nil {
			return "", errors.Wrap(err, "parse command id")
		}
		cid = _cid
	}
	_oid, err := strconv.ParseUint(string(subs[6]), 16, 32)
	if err != nil {
		return "", errors.Wrap(err, "parse object id")
	}
	oid = _oid

	obj = string(subs[5])
	objString := ""
	if obj != "" {
		objString = "-" + obj
	}
	if oid > 0xFFFF {
		return PrivObjectId(fmt.Sprintf("T%02X%v.P%08X", cid, objString, oid)), nil
	}
	return PrivObjectId(fmt.Sprintf("T%02X%v.P%04X", cid, objString, oid)), nil
}

type SnmpObjectId string

func (m SnmpObjectId) Short() string {
	tmp := string(m)
	tmp = strings.TrimSuffix(tmp, ".0")
	if parts := strings.Split(tmp, "."); len(parts) > 3 {
		return strings.Join(parts[len(parts)-3:], ".")
	}
	return string(m)
}
