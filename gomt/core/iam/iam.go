package iam

import "sync"

var allRoles = []*Role{
	// {Schema: "", Name: "super", Default: true},
	{Schema: "", Name: "admin", Default: true},
	{Schema: "", Name: "guest", Default: true},
}
var setupAllRolesOnce sync.Once

func GetRoles() []*Role {
	setupAllRolesOnce.Do(func() {
		for _, role := range allRoles {
			role.GetRules()
		}
	})
	return allRoles
}
