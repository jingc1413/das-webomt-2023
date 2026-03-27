package agent

import (
	"gomt/core/iam"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s *DasDeviceAgent) ServeAuthUser(name string, password string) (bool, error) {
	if s.isLocalDevice && s.userMgmt != nil {
		return s.userMgmt.AuthUser(name, password), nil
	} else if s.supportCGI && s.cgiHandler != nil {
		result, err := s.cgiHandler.ServeLogin(name, password)
		if err != nil {
			return false, errors.Wrap(err, "serve auth")
		}
		logrus.Tracef("login result: %v", result)
		if strings.Contains(result, "login success") || strings.Contains(result, "multi login pass") {
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("not supported")
}

func (s *DasDeviceAgent) ServeGetUsers() ([]*iam.User, error) {
	if s.isLocalDevice && s.userMgmt != nil {
		return s.userMgmt.GetUsers(), nil
	} else if s.supportCGI && s.cgiHandler != nil {
		return s.cgiHandler.ServeGetUsers()
	}
	return nil, errors.New("not supported")
}

func (s *DasDeviceAgent) ServeGetUser(name string) (*iam.User, error) {
	if s.isLocalDevice && s.userMgmt != nil {
		return s.userMgmt.GetUser(name), nil
	} else if s.supportCGI && s.cgiHandler != nil {
		users, err := s.cgiHandler.ServeGetUsers()
		if err != nil {
			return nil, errors.Wrap(err, "serve get users")
		}
		for _, v := range users {
			if v.Name == name {
				return v, nil
			}
		}
		return nil, nil
	}
	return nil, errors.New("not supported")
}
func (s *DasDeviceAgent) ServeCreateUser(name string, password string, roles []string) error {
	if s.isLocalDevice && s.userMgmt != nil {
		return s.userMgmt.CreateUser(name, password, roles)
	} else if s.supportCGI && s.cgiHandler != nil {
		return s.cgiHandler.ServeCreateUser(name, password)
	}
	return errors.New("not supported")
}

func (s *DasDeviceAgent) ServeSetUserPassword(name string, password string) error {
	if s.isLocalDevice && s.userMgmt != nil {
		return s.userMgmt.SetUserPassword(name, password)
	} else if s.supportCGI && s.cgiHandler != nil {
		return s.cgiHandler.ServeSetUserPassword(name, password)
	}
	return errors.New("not supported")
}

func (s *DasDeviceAgent) ServeSetUserRoles(name string, roles []string) error {
	if s.isLocalDevice && s.userMgmt != nil {
		return s.userMgmt.SetUserRoles(name, roles)
	}
	return errors.New("not supported")
}

func (s *DasDeviceAgent) ServeDeleteUser(name string) error {
	if s.isLocalDevice && s.userMgmt != nil {
		return s.userMgmt.DeleteUser(name)
	} else if s.supportCGI && s.cgiHandler != nil {
		return s.cgiHandler.ServeDeleteUser(name)
	}
	return errors.New("not supported")
}
