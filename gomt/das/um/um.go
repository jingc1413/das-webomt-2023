package um

import (
	"gomt/core/iam"

	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type UserManager struct {
	filename string
	f        *ini.File
}

func NewUserManager(configPath string) *UserManager {
	s := &UserManager{
		filename: path.Join(configPath, "Key.ini"),
	}
	f, err := ini.Load(s.filename)
	if err != nil {
		logrus.Error(errors.Wrap(err, "load ini"))
		return nil
	}
	s.f = f
	return s
}

func (s *UserManager) AuthUser(name string, password string) bool {
	userSection := s.f.Section("user")

	if userSection != nil {
		for _, userKey := range userSection.Keys() {
			if userKey.Name() == name && userKey.Value() == password {
				return true
			}
		}
	}
	return false
}

func (s *UserManager) GetUsers() []*iam.User {
	out := []*iam.User{}

	userSection := s.f.Section("user")
	rolesSection := s.f.Section("roles")

	if userSection != nil {
		for _, userKey := range userSection.Keys() {
			name := userKey.Name()
			password := userKey.Value()
			roles := []string{}
			if rolesSection != nil {
				if rolesKey := rolesSection.Key(name); rolesKey != nil {
					if rolesText := rolesKey.Value(); rolesText != "" {
						roles = strings.Split(rolesText, "|")
					}
				}
			}
			user := iam.MakeUser(name, password, roles)
			out = append(out, user)
		}
	}
	return out
}

func (s *UserManager) GetUser(name string) *iam.User {
	userSection := s.f.Section("user")
	rolesSection := s.f.Section("roles")

	if userSection != nil {
		for _, userKey := range userSection.Keys() {
			if userKey.Name() == name {
				password := userKey.Value()
				roles := []string{}
				if rolesSection != nil {
					if rolesKey := rolesSection.Key(name); rolesKey != nil {
						if rolesText := rolesKey.Value(); rolesText != "" {
							roles = strings.Split(rolesText, "|")
						}
					}
				}
				return iam.MakeUser(name, password, roles)
			}
		}
	}
	return nil
}

func (s *UserManager) CreateUser(name string, password string, roles []string) error {
	userSection := s.f.Section("user")
	rolesSection := s.f.Section("roles")
	if userSection == nil {
		section, err := s.f.NewSection("user")
		if err != nil {
			return errors.Wrap(err, "new user section")
		}
		userSection = section
	}
	if rolesSection == nil {
		section, err := s.f.NewSection("roles")
		if err != nil {
			return errors.Wrap(err, "new roles section")
		}
		rolesSection = section
	}

	if userKey := userSection.Key(name); userKey != nil {
		return errors.New("already exists")
	}
	if _, err := userSection.NewKey(name, password); err != nil {
		return errors.Wrap(err, "new user key")
	}

	rolesKey := rolesSection.Key(name)
	if rolesKey == nil {
		key, err := rolesSection.NewKey(name, "")
		if err != nil {
			return errors.Wrap(err, "new roles key")
		}
		rolesKey = key
	}
	rolesText := ""
	if len(roles) > 0 {
		rolesText = strings.Join(roles, "|")
	}
	rolesKey.SetValue(rolesText)

	if err := s.f.SaveTo(s.filename); err != nil {
		return errors.Wrap(err, "save to file")
	}

	return nil
}

func (s *UserManager) DeleteUser(name string) error {
	userSection := s.f.Section("user")
	rolesSection := s.f.Section("roles")
	if userSection != nil {
		userSection.DeleteKey(name)
	}
	if rolesSection != nil {
		rolesSection.DeleteKey(name)
	}

	if err := s.f.SaveTo(s.filename); err != nil {
		return errors.Wrap(err, "save to file")
	}
	return nil
}

func (s *UserManager) SetUserPassword(name string, password string) error {
	userSection := s.f.Section("user")
	if userSection == nil {
		section, err := s.f.NewSection("user")
		if err != nil {
			return errors.Wrap(err, "new user section")
		}
		userSection = section
	}

	userKey := userSection.Key(name)
	if userKey == nil {
		return errors.New("not found")
	}
	userKey.SetValue(password)

	if err := s.f.SaveTo(s.filename); err != nil {
		return errors.Wrap(err, "save to file")
	}
	return nil
}

func (s *UserManager) SetUserRoles(name string, roles []string) error {
	userSection := s.f.Section("user")
	rolesSection := s.f.Section("roles")
	if userSection == nil {
		section, err := s.f.NewSection("user")
		if err != nil {
			return errors.Wrap(err, "new user section")
		}
		userSection = section
	}
	if rolesSection == nil {
		section, err := s.f.NewSection("roles")
		if err != nil {
			return errors.Wrap(err, "new roles section")
		}
		rolesSection = section
	}

	userKey := userSection.Key(name)
	if userKey == nil {
		return errors.New("not found")
	}

	rolesKey := rolesSection.Key(name)
	if rolesKey == nil {
		key, err := rolesSection.NewKey(name, "")
		if err != nil {
			return errors.Wrap(err, "new roles key")
		}
		rolesKey = key
	}
	rolesText := ""
	if len(roles) > 0 {
		rolesText = strings.Join(roles, "|")
	}
	rolesKey.SetValue(rolesText)

	if err := s.f.SaveTo(s.filename); err != nil {
		return errors.Wrap(err, "save to file")
	}

	return nil
}
