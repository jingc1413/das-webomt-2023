package server

import (
	"fmt"
	"gomt/core/iam"
	"gomt/core/layout"

	"github.com/pkg/errors"
)

var rulesMap = map[string]iam.Rule{
	// api
	"api.app.info.get": {Method: "GET", ObjectType: "api", ObjectID: "/app/info", Public: true},

	"api.auth.login":  {Method: "POST", ObjectType: "api", ObjectID: "/auth/login", Public: true},
	"api.auth.logout": {Method: "GET", ObjectType: "api", ObjectID: "/auth/logout", Public: true},

	"api.current.get":             {Method: "GET", ObjectType: "api", ObjectID: "/current", Default: true},
	"api.current.change-password": {Method: "POST", ObjectType: "api", ObjectID: "/current/change-password", Default: true},

	"api.iam.roles.list":         {Method: "GET", ObjectType: "api", ObjectID: "/iam/roles", Public: true},
	"api.iam.users.list":         {Method: "GET", ObjectType: "api", ObjectID: "/iam/users"},
	"api.iam.users.create":       {Method: "POST", ObjectType: "api", ObjectID: "/iam/users"},
	"api.iam.users.delete":       {Method: "DELETE", ObjectType: "api", ObjectID: "/iam/users/:name"},
	"api.iam.users.update":       {Method: "POST", ObjectType: "api", ObjectID: "/iam/users/:name"},
	"api.iam.users.set-password": {Method: "POST", ObjectType: "api", ObjectID: "/iam/users/:name/set-password"},

	"api.diag.ping.jobs.create": {Method: "POST", ObjectType: "api", ObjectID: "/diag/ping/jobs"},
	"api.diag.ping.jobs.ws":     {Method: "GET", ObjectType: "api", ObjectID: "/diag/ping/jobs/:token/ws"},
	"api.diag.ping.jobs.run":    {Method: "POST", ObjectType: "api", ObjectID: "/diag/ping/jobs/:token/run"},
	"api.diag.ping.jobs.cancel": {Method: "POST", ObjectType: "api", ObjectID: "/diag/ping/jobs/:token/cancel"},

	"api.das.websocket":                         {Method: "GET", ObjectType: "api", ObjectID: "/das/ws"},
	"api.das.products.list":                     {Method: "GET", ObjectType: "api", ObjectID: "/das/products"},
	"api.das.products.get":                      {Method: "GET", ObjectType: "api", ObjectID: "/das/products/:name"},
	"api.das.device-types.list":                 {Method: "GET", ObjectType: "api", ObjectID: "/das/device-types"},
	"api.das.device-types.model.layout.get":     {Method: "GET", ObjectType: "api", ObjectID: "/das/device-types/:name/:version/model/layout"},
	"api.das.device-types.model.parameters.get": {Method: "GET", ObjectType: "api", ObjectID: "/das/device-types/:name/:version/model/parameters"},
	"api.das.query-devices.progress.get":        {Method: "GET", ObjectType: "api", ObjectID: "/das//query-devices/progress"},

	"api.das.devices.list":           {Method: "GET", ObjectType: "api", ObjectID: "/das/devices"},
	"api.das.devices.get":            {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub"},
	"api.das.devices.type.get":       {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/type"},
	"api.das.devices.available.get":  {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/available"},
	"api.das.devices.parameters.get": {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/parameters/get"},
	"api.das.devices.parameters.set": {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/parameters/set"},
	"api.das.devices.register.read":  {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/register/read"},
	"api.das.devices.register.write": {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/register/write"},

	"api.das.devices.metrics.data.query":  {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/metrics/data/query"},
	"api.das.devices.metrics.current.get": {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/metrics/current"},

	"api.das.devices.files.list":              {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/files/:ftype"},
	"api.das.devices.files.create":            {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/files/:ftype"},
	"api.das.devices.files.get":               {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/files/:ftype/:fname"},
	"api.das.devices.files.delete":            {Method: "DELETE", ObjectType: "api", ObjectID: "/das/devices/:device_sub/files/:ftype/:fname"},
	"api.das.devices.files.packet-info.get":   {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/files/:ftype/:fname/packet-info"},
	"api.das.devices.version.packet-info.get": {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/version/packet-info"},

	"api.das.devices.alarm.logs.list": {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/alarm/logs"},

	"api.das.devices.firmwares.list":      {Method: "GET", ObjectType: "api", ObjectID: "/das/devices/:device_sub/firmwares"},
	"api.das.devices.firmwares.delete":    {Method: "DELETE", ObjectType: "api", ObjectID: "/das/devices/:device_sub/firmwares/:name"},
	"api.das.devices.upgrade.start":       {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/upgrade/start"},
	"api.das.devices.upgrade.reboot":      {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/upgrade/reboot"},
	"api.das.devices.delete-key-and-logs": {Method: "POST", ObjectType: "api", ObjectID: "/das/devices/:device_sub/delete-key-and-logs"},
}

func getApiRule(method string, uri string) (string, *iam.Rule) {
	for id, v := range rulesMap {
		if v.Method == method && v.ObjectType == "api" && v.ObjectID == uri {
			tmp := v
			return id, &tmp
		}
	}
	return "", nil
}

var defaultRoles = []*iam.Role{
	{Name: "super", Default: true, Rules: []string{"api.iam.*"}},
	{Name: "admin", Default: true, Rules: []string{"api.*", "page.*"}},
	{Name: "guest", Default: true, Rules: []string{"api.*.list", "api.*.get", "page.*.get"}},
}

var defaultUsers = []*iam.User{
	iam.MakeUser("root0", "root@SW24", []string{"super"}),
	iam.MakeUser("admin", "admin", []string{"admin"}),
	iam.MakeUser("guest", "guest", []string{"guest"}),
}

func getDefaultUser(name string) *iam.User {
	name = iam.EncodeUserName(name)
	for _, user := range defaultUsers {
		if user.Name == name {
			tmp := *user
			return &tmp
		}
	}
	return nil
}

func (s *OMTServer) setupIAM() error {
	agent := s.dassys.GetAgent("local")
	if agent == nil {
		return errors.New("get primary agent")
	}

	pathsMap := layout.GetAllPathsMap()
	for _, paths := range pathsMap {
		for _, path := range paths {
			key := fmt.Sprintf("page.%v.get", path)
			key2 := fmt.Sprintf("page.%v.set", path)

			rulesMap[key] = iam.Rule{Method: "GET", ObjectType: "page", ObjectID: path}
			rulesMap[key2] = iam.Rule{Method: "SET", ObjectType: "page", ObjectID: path}
		}
	}
	if err := iam.Setup(rulesMap, defaultRoles, defaultUsers); err != nil {
		return err
	}
	return nil
}
