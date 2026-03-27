package server

import (
	"encoding/json"
	"fmt"
	"gomt/core/iam"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var SECRET = "88866999"
var DEAFULT_EXPIRE_SECONDS = 1800

func init() {
	const CHARSET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	secret := ""
	for i := 0; i < 16; i++ {
		x := r.Intn(len(CHARSET))
		c := CHARSET[x]
		secret = secret + fmt.Sprintf("%s", string(c))
	}
	SECRET = secret
}

func (s *OMTServer) makeLoginHandler(prefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "not login")
			}
			schema := cast.ToString(sess.Values["schema"])
			username := cast.ToString(sess.Values["login"])
			firstTimeLogin := cast.ToInt(sess.Values["firstTimeLogin"])
			expireSeconds := cast.ToInt64(sess.Values["expireSeconds"])
			expireTime := cast.ToInt64(sess.Values["expireTime"])

			c.Set("login", username)
			c.Set("schema", schema)
			c.Set("firstTimeLogin", firstTimeLogin)

			path := c.Path()
			method := c.Request().Method
			if prefix != "" {
				if !strings.HasPrefix(path, prefix) {
					return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
				}
				path = path[len(prefix):]
			}

			if username != "" {
				if expireSeconds <= 0 {
					expireSeconds = int64(DEAFULT_EXPIRE_SECONDS)
				}
				ts := time.Now().Unix()
				if ts < expireTime-expireSeconds || ts > expireTime {
					s.auditLogger.WriteApiLog(time.Now(), "info", "expired", username, c.RealIP(), 200, "login expired")
					return echo.NewHTTPError(http.StatusUnauthorized, "expired")
				}
				sess.Values["expireTime"] = time.Now().Unix() + expireSeconds
				if err := sess.Save(c.Request(), c.Response()); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "save session")
				}
				logrus.Warnf("%v, %v, %v", username, expireTime, expireSeconds)
				if ok := iam.EnforceObject(schema, "*", "api", path, method); ok {
					return next(c)
				}
				if ok := iam.EnforceObject(schema, username, "api", path, method); ok {
					return next(c)
				}
			}
			if ok := iam.EnforceObject(schema, "", "api", path, method); ok {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden, "forbidden")
		}
	}
}

func (s *OMTServer) mustGetSessionExpireSeconds() int64 {
	if agent := s.getDasDeviceAgent("0"); agent != nil {
		if v, err := agent.GetParameterValueOfSessionExpireSeconds(); err == nil && v > 0 {
			return v * 60
		}
	}
	return int64(DEAFULT_EXPIRE_SECONDS)
}

type LoginData struct {
	Username string
	Password string
}

func (s *OMTServer) handleLogin(c echo.Context) error {
	schema := ""
	var loginData LoginData
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &loginData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if loginData.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid username")
	}
	if loginData.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
	}
	loginData.Username = iam.EncodeUserName(loginData.Username)

	ok := false
	agent := s.getDasDeviceAgent("local")
	if user, err := agent.ServeGetUser(loginData.Username); err == nil && user != nil {
		ok, err = agent.ServeAuthUser(loginData.Username, loginData.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	if !ok {
		if user := getDefaultUser(loginData.Username); user != nil {
			if user.CheckPassword(loginData.Password) {
				ok = true
			}
		}
	}
	if !ok {
		s.auditLogger.WriteApiLog(time.Now(), "info", "login", loginData.Username, c.RealIP(), 200, "login failed, incorrect username or password")
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect username or password")
	}

	firstTimeLogin := 0
	for _, defaultUser := range defaultUsers {
		if defaultUser.Name == loginData.Username && defaultUser.Password == loginData.Password {
			firstTimeLogin = 1
			break
		}
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}
	expireSeconds := s.mustGetSessionExpireSeconds()
	sess.Values["login"] = loginData.Username
	sess.Values["schema"] = schema
	sess.Values["expireSeconds"] = expireSeconds
	sess.Values["expireTime"] = time.Now().Unix() + expireSeconds
	sess.Values["firstTimeLogin"] = firstTimeLogin
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	s.auditLogger.WriteApiLog(time.Now(), "info", "login", loginData.Username, c.RealIP(), 200, nil)
	return c.NoContent(http.StatusOK)
}

func (s *OMTServer) handleLogout(c echo.Context) error {
	name := cast.ToString(c.Get("login"))
	s.auditLogger.WriteApiLog(time.Now(), "info", "logout", name, c.RealIP(), 200, nil)
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (s *OMTServer) handleGetCurrentUser(c echo.Context) error {
	name := cast.ToString(c.Get("login"))
	firstTimeLogin := cast.ToInt(c.Get("firstTimeLogin"))

	name = iam.EncodeUserName(name)
	if name == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	agent := s.getDasDeviceAgent("local")
	user, err := agent.ServeGetUser(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		user = getDefaultUser(name)
	}
	if user == nil {
		return c.NoContent(http.StatusNotFound)
	}

	user.GetRules()
	user.Name = iam.DecodeUserName(user.Name)
	user.Password = ""
	if firstTimeLogin != 0 {
		user.FirstTimeLogin = true
	}
	return c.JSON(http.StatusOK, user)
}

type changePasswordBody struct {
	Password    string `json:"Password"`
	NewPassword string `json:"NewPassword"`
}

func (s *OMTServer) handleChangeCurrentUserPassword(c echo.Context) error {
	name := cast.ToString(c.Get("login"))
	name = iam.EncodeUserName(name)
	if name == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var data changePasswordBody
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "invalid input parameter").Error())
	}

	if data.NewPassword == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
	}

	agent := s.getDasDeviceAgent("local")
	user, err := agent.ServeGetUser(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "get user").Error())
	} else if user != nil {
		if user.CheckPassword(data.Password) == false {
			return echo.NewHTTPError(http.StatusBadRequest, "incorrect password")
		}
		if err := agent.ServeSetUserPassword(name, data.NewPassword); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, nil)
	}

	if user = getDefaultUser(name); user != nil {
		if user.CheckPassword(data.Password) == false {
			return echo.NewHTTPError(http.StatusBadRequest, "incorrect password")
		}
		if err := agent.ServeCreateUser(user.Name, data.NewPassword, user.Roles); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, nil)
	}

	return echo.NewHTTPError(http.StatusForbidden, "not support")
}

func (s *OMTServer) handleListRoles(c echo.Context) error {
	roles := iam.GetRoles()
	return c.JSON(http.StatusOK, roles)
}

func (s *OMTServer) handleListUsers(c echo.Context) error {
	agent := s.getDasDeviceAgent("local")
	users, err := agent.ServeGetUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	result := []*iam.User{}
	for _, user := range users {
		if user.Hidden {
			continue
		}
		user.Password = ""
		result = append(result, user)
	}
	return c.JSON(http.StatusOK, result)
}

type createUserData struct {
	Name     string
	Password string
	Roles    []string
}

func (s *OMTServer) handleCreateUser(c echo.Context) error {
	var data createUserData
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if data.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid name")
	}
	if data.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
	}
	data.Name = iam.EncodeUserName(data.Name)
	if data.Name == "root0" || data.Name == "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}

	agent := s.getDasDeviceAgent("local")
	if user, err := agent.ServeGetUser(data.Name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user != nil {
		return echo.NewHTTPError(http.StatusForbidden, "already exists")
	}

	if err := agent.ServeCreateUser(data.Name, data.Password, data.Roles); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (s *OMTServer) handleDeleteUser(c echo.Context) error {
	name := c.Param("name")
	name = iam.EncodeUserName(name)
	if name == "root0" || name == "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}

	agent := s.getDasDeviceAgent("local")
	if user, err := agent.ServeGetUser(name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "not exists")
	}

	if err := agent.ServeDeleteUser(name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

type updateUserInfoData struct {
	Roles []string
}

func (s *OMTServer) handleUpdateUser(c echo.Context) error {
	schema := cast.ToString(c.Get("schema"))
	name := c.Param("name")
	name = iam.EncodeUserName(name)
	if name == "root0" || name == "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}

	var data updateUserInfoData
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	agent := s.getDasDeviceAgent("local")
	if user, err := agent.ServeGetUser(name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "not exists")
	}
	if err := agent.ServeSetUserRoles(name, data.Roles); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if ok, err := iam.SetRolesForUser(schema, name, data.Roles); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	} else if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "already has the role")
	}
	return c.NoContent(http.StatusOK)
}

type updateUserPasswordData struct {
	Password string
}

func (s *OMTServer) handleUpdateUserPassword(c echo.Context) error {
	name := c.Param("name")
	name = iam.EncodeUserName(name)
	if name == "root0" || name == "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}

	var data updateUserPasswordData
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if data.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
	}

	agent := s.getDasDeviceAgent("local")
	if user, err := agent.ServeGetUser(name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "not exists")
	}
	if err := agent.ServeSetUserPassword(name, data.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
