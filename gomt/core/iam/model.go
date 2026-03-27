package iam

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Role struct {
	Schema  string `json:"Schema,omitempty"`
	Name    string
	Default bool     `json:"Default,omitempty"`
	Rules   []string `json:"Rules,omitempty"`
}

func (m *Role) GetRules() {
	rules := GetRulesForRole(m.Schema, m.Name, "")
	m.Rules = GetRuleIdsByRules(rules)
}

type User struct {
	Schema           string `json:"Schema,omitempty"`
	Name             string
	Password         string `json:"-"`
	PasswordExpireAt int64  `json:"PasswordExpireAt,omitempty"`
	Default          bool   `json:"Default,omitempty"`
	Hidden           bool   `json:"Hidden,omitempty"`
	Roles            []string
	Rules            []string `json:"Rules,omitempty"`
	FirstTimeLogin   bool     `json:"FirstTimeLogin,omitempty"`
}

func (m *User) GetRules() {
	ids := []string{}
	for _, role := range m.Roles {
		if rules := GetRulesForRole(m.Schema, role, ""); rules != nil {
			ids = append(ids, GetRuleIdsByRules(rules)...)
		}
	}
	m.Rules = ids
}

func (s *User) SetPassword(password string) {
	s.Password = EncryptPassword(password)
}

func (s User) CheckPassword(password string) bool {
	if password != "" && s.Password == EncryptPassword(password) {
		return true
	}
	return false
}

func (s *User) ChangePassword(old string, password string) bool {
	if !s.CheckPassword(old) {
		return false
	}
	s.SetPassword(password)
	return true
}

func EncryptPassword(password string) string {
	if true {
		return password
	} else {
		h := sha256.New()
		io.WriteString(h, password)
		sum := h.Sum(nil)
		return base64.StdEncoding.EncodeToString(sum)
	}
}

func GeneratePassword() string {
	const CHARSET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := ""
	for i := 0; i < 8; i++ {
		x := r.Intn(len(CHARSET))
		c := CHARSET[x]
		password = password + fmt.Sprintf("%s", string(c))
	}
	return password
}

func MakeUser(name string, password string, roles []string) *User {
	user := &User{
		Name:     name,
		Password: password,
		Roles:    roles,
	}
	if user.Name == "admin" {
		user.Default = true
		user.Hidden = false
		user.Roles = []string{"admin"}
	}
	if user.Name == "root0" {
		user.Default = true
		user.Hidden = true
		user.Roles = []string{"super"}
	}
	return user
}

func EncodeUserName(name string) string {
	if name == "root" {
		return "root0"
	}
	return name
}

func DecodeUserName(name string) string {
	if name == "root0" {
		return "root"
	}
	return name
}
