package iam

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const modelString string = `
[request_definition]
r = sub, obj, act, dom

[policy_definition]
p = sub, obj, act, dom

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`

var defaultEnforcer *casbin.Enforcer
var defaultEnforcerSetupOnce sync.Once

func GetDefaultEnforcer() *casbin.Enforcer {
	defaultEnforcerSetupOnce.Do(func() {
		m, err := model.NewModelFromString(modelString)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "setup model for enforcer"))
		}
		f, err := os.OpenFile("policy.csv", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "setup policy file for enforcer"))
		}
		f.Close()
		a := fileadapter.NewAdapter("policy.csv")
		e, err := casbin.NewEnforcer(m, a)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "setup enforcer"))
		}
		defaultEnforcer = e
	})
	return defaultEnforcer
}

func Enforce(dom string, sub string, obj string, act string) bool {
	e := GetDefaultEnforcer()
	ok, err := e.Enforce(sub, obj, act, dom)
	if err != nil {
		return false
	}
	return ok
}

func EnforceObject(domain string, user string, objectType string, objectID string, method string) bool {
	ok := Enforce(
		fmt.Sprintf("d:%v", domain),
		fmt.Sprintf("u:%v", user),
		fmt.Sprintf("%v:%v", objectType, objectID),
		strings.ToUpper(method),
	)
	// logrus.Warnf("enforce %v %v %v %v", user, objectType, objectID, ok)
	return ok
}

func GetRolesForUser(domain string, user string) ([]string, error) {
	roles := []string{}
	e := GetDefaultEnforcer()
	_roles, err := e.GetRolesForUser(fmt.Sprintf("u:%v", user), fmt.Sprintf("d:%v", domain))
	if err != nil {
		return roles, err
	}
	for _, v := range _roles {
		if strings.HasPrefix(v, "r:") {
			roles = append(roles, v[2:])
		}
	}
	return roles, nil
}

func DeleteAllRolesForUser(domain string, user string) (bool, error) {
	e := GetDefaultEnforcer()
	return e.DeleteRolesForUser(
		fmt.Sprintf("u:%v", user),
		fmt.Sprintf("d:%v", domain),
	)
}

func AddRoleForUser(domain string, user string, role string) (bool, error) {
	e := GetDefaultEnforcer()

	return e.AddRoleForUser(
		fmt.Sprintf("u:%v", user),
		fmt.Sprintf("r:%v", role),
		fmt.Sprintf("d:%v", domain),
	)
}

func AddRolesForUser(domain string, user string, roles []string) (bool, error) {
	_roles := []string{}
	for _, role := range roles {
		_roles = append(_roles, fmt.Sprintf("r:%v", role))
	}
	e := GetDefaultEnforcer()
	return e.AddRolesForUser(
		fmt.Sprintf("u:%v", user),
		_roles,
		fmt.Sprintf("d:%v", domain),
	)
}

func SetRolesForUser(domain string, user string, roles []string) (bool, error) {
	DeleteAllRolesForUser(domain, user)
	return AddRolesForUser(domain, user, roles)
}
