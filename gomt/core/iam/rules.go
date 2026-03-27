package iam

import (
	"regexp"
	"strings"
)

var defaultRulesMap = map[string]Rule{}

func SetupRules(rules map[string]Rule) {
	if rules != nil {
		for key, rule := range rules {
			defaultRulesMap[key] = rule
		}
	}
}

func GetRuleIds(patterns ...string) []string {
	ids := []string{}
	ruleMap := make(map[string]bool)
	for _, pattern := range patterns {
		pattern = strings.ReplaceAll(pattern, ".", "\\.")
		pattern = strings.ReplaceAll(pattern, "*", ".*")
		re := regexp.MustCompile(pattern)
		for id, _ := range defaultRulesMap {
			if re.Match([]byte(id)) {
				if _, ok := ruleMap[id]; !ok {
					ruleMap[id] = true
					ids = append(ids, id)
				}
			}
		}
	}
	return ids
}

func GetRules(patterns ...string) []Rule {
	rules := []Rule{}
	ruleMap := make(map[string]bool)
	for _, pattern := range patterns {
		pattern = strings.ReplaceAll(pattern, ".", "\\.")
		pattern = strings.ReplaceAll(pattern, "*", ".*")
		re := regexp.MustCompile(pattern)
		for id, rule := range defaultRulesMap {
			if re.Match([]byte(id)) {
				if _, ok := ruleMap[id]; !ok {
					ruleMap[id] = true
					rules = append(rules, rule)
				}
			}
		}
	}
	return rules
}

func GetRuleIdsByRules(rules []Rule) []string {
	ids := []string{}
	ruleMap := make(map[string]bool)
	for _, rule := range rules {
		for id, rule2 := range defaultRulesMap {
			if rule.Method == rule2.Method && rule.ObjectType == rule2.ObjectType && rule.ObjectID == rule2.ObjectID {
				if _, ok := ruleMap[id]; !ok {
					ruleMap[id] = true
					ids = append(ids, id)
				}
			}
		}
	}
	return ids
}

func GetRuleIdByRule(rule Rule) string {
	for id, rule2 := range defaultRulesMap {
		if rule.Method == rule2.Method && rule.ObjectType == rule2.ObjectType && rule.ObjectID == rule2.ObjectID {
			return id
		}
	}
	return ""
}
