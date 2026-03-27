package iam

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Rule struct {
	ObjectType string
	ObjectID   string
	Method     string
	Public     bool
	Default    bool
}

func DumpRules(rules []Rule) string {
	lines := []string{}
	for _, rule := range rules {
		lines = append(lines, fmt.Sprintf("%v %v %v", rule.Method, rule.ObjectType, rule.ObjectID))
	}
	return strings.Join(lines, "\n")
}

func GetRulesForRole(domain string, role string, objectType string) []Rule {
	rules := []Rule{}
	e := GetDefaultEnforcer()
	polices, err := e.GetFilteredPolicy(0, fmt.Sprintf("r:%v", role), "", "", fmt.Sprintf("d:%v", domain))
	if err != nil {
		logrus.Error(errors.Wrap(err, "get filtered policy"))
		return rules
	}
	for _, v := range polices {
		args := strings.Split(v[1], ":")
		if len(args) != 2 {
			continue
		}
		if objectType == "" || v[1][0:1] == objectType {
			rules = append(rules, Rule{ObjectType: args[0], ObjectID: args[1], Method: v[2]})
		}
	}
	return rules
}

func GetRulesForUser(domain string, user string, objectType string) []Rule {
	rules := []Rule{}
	roles, _ := GetRolesForUser(domain, user)
	for _, role := range roles {
		_rules := GetRulesForRole(domain, role, objectType)
		rules = append(rules, _rules...)
	}
	return rules
}

func DeleteAllRulesForRole(domain string, role string, objectType string) (bool, error) {
	polices := [][]string{}
	e := GetDefaultEnforcer()

	_polices, err := e.GetFilteredPolicy(0, fmt.Sprintf("r:%v", role), "", "", fmt.Sprintf("d:%v", domain))
	if err != nil {
		return false, errors.Wrap(err, "get filtered policy")
	}

	for _, v := range _polices {
		if objectType == "" || v[1][0:1] == objectType {
			polices = append(polices, v)
		}
	}

	if len(polices) > 0 {
		return e.RemovePolicies(polices)
	}
	return true, nil
}

func SetRulesForRole(domain string, role string, rules []Rule) (bool, error) {
	if role == "" {
		return false, errors.Errorf("invalid role, role=%v", role)
	}
	e := GetDefaultEnforcer()

	DeleteAllRulesForRole(domain, role, "")

	if len(rules) > 0 {
		policies := [][]string{}
		for _, rule := range rules {
			if rule.ObjectType == "" && rule.ObjectID == "" {
				return false, errors.Errorf("invalid rule, objectType=%v objectID=%v", rule.ObjectType, rule.ObjectID)
			}
			policies = append(policies, []string{
				fmt.Sprintf("r:%v", role),
				fmt.Sprintf("%v:%v", rule.ObjectType, rule.ObjectID),
				strings.ToUpper(rule.Method),
				fmt.Sprintf("d:%v", domain),
			})
		}

		return e.AddPolicies(policies)
	}
	return true, nil
}

func AddRuleForRole(domain string, role string, rule Rule) (bool, error) {
	if role == "" {
		return false, errors.Errorf("invalid role, role=%v", role)
	}
	e := GetDefaultEnforcer()
	if rule.ObjectType == "" && rule.ObjectID == "" {
		return false, errors.Errorf("invalid rule, objectType=%v objectID=%v", rule.ObjectType, rule.ObjectID)
	}
	return e.AddPolicy(
		fmt.Sprintf("r:%v", role),
		fmt.Sprintf("%v:%v", rule.ObjectType, rule.ObjectID),
		strings.ToUpper(rule.Method),
		fmt.Sprintf("d:%v", domain),
	)
}
