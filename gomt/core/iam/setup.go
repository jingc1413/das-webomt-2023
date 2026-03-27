package iam

import (
	"github.com/pkg/errors"
)

func Setup(rules map[string]Rule, roles []*Role, users []*User) error {
	schema := ""

	GetDefaultEnforcer()
	SetupRules(rules)

	publicRules := []Rule{}
	defaultRules := []Rule{}

	for _, rule := range rules {
		if rule.Public {
			publicRules = append(publicRules, rule)
		} else if rule.Default {
			defaultRules = append(defaultRules, rule)
		}
	}
	DeleteAllRulesForRole(schema, "_public", "")
	SetRulesForRole(schema, "_public", publicRules)
	DeleteAllRulesForRole(schema, "_default", "")
	SetRulesForRole(schema, "_default", defaultRules)

	SetRolesForUser(schema, "", []string{"_public"})
	SetRolesForUser(schema, "*", []string{"_default"})

	for _, role := range roles {
		if apiRules := GetRules(role.Rules...); apiRules != nil {
			if ok, err := SetRulesForRole(role.Schema, role.Name, apiRules); !ok {
				return errors.Errorf("setup api rules for %v, not ok", role.Name)
			} else if err != nil {
				return errors.Wrapf(err, "setup api rules for %v", role.Name)
			}
			// if rules := GetRulesForRole(role.Schema, role.Name, ""); rules != nil {
			// 	logrus.Tracef("!!!! %v role rules\n%v", role.Name, DumpRules(rules))
			// }
		}
	}

	for _, user := range users {
		SetRolesForUser(user.Schema, user.Name, user.Roles)
	}

	GetDefaultEnforcer().SavePolicy()

	return nil
}

// group.POST("/auth", HandleAuthenticate)

// meGroup := group.Group("/current")
// meGroup.GET("", HandleWhoAmI)
// meGroup.PUT("/change-password", HandleChangePassword)

// group.GET("/rules", HandleListApiRules)

// group.GET("/roles", HandleListRoles)
// group.POST("/roles", HandleCreateRole)
// group.DELETE("/roles/:id", HandleDeleteRole, echohandler.ObjectHandler[Role](db))
// group.GET("/roles/:id", HandleGetRole, echohandler.ObjectHandler[Role](db))
// group.PUT("/roles/:id", HandleUpdateRole, echohandler.ObjectHandler[Role](db))

// group.GET("/users", HandleListUsers)
// group.POST("/users", HandleCreateUser)
// group.DELETE("/users/:id", HandleDeleteUser, echohandler.ObjectHandler[User](db))
// group.GET("/users/:id", HandleGetUser, echohandler.ObjectHandler[User](db))
// group.PUT("/users/:id", HandleUpdateUser, echohandler.ObjectHandler[User](db))
// group.PUT("/users/:id/reset-password", HandleResetUserPassword, echohandler.ObjectHandler[User](db))
// group.PUT("/users/:id/change-password", HandleChangeUserPassword, echohandler.ObjectHandler[User](db))
// group.GET("/users/:id/roles", HandleGetUserRoles, echohandler.ObjectHandler[User](db))
// group.POST("/users/:id/roles", HandleAddUserRoles, echohandler.ObjectHandler[User](db))
// group.PUT("/users/:id/roles", HandleSetUserRoles, echohandler.ObjectHandler[User](db))
