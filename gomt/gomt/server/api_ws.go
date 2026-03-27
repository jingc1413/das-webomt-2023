package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type DasNotifyMessage struct {
	Type string `json:"Type"`
	Data any    `json:"Data"`
}

func (s *OMTServer) sendDasNotifyMessage(typ string, data any) {
	msg := DasNotifyMessage{
		Type: typ,
		Data: data,
	}
	content, err := json.Marshal(msg)
	if err != nil {
		s.log.Error(errors.Wrap(err, "marshal websocket message content"))
		return
	}

	if s.wsConns != nil {
		for _, conn := range s.wsConns {
			if err := conn.WriteMessage(websocket.TextMessage, content); err != nil {
				s.log.Error(errors.Wrap(err, "write websocket message"))
			}
		}
	}
	return
}

func (s *OMTServer) handleDasWebSocket(c echo.Context) error {
	conn, err := s.wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		s.log.Error(err)
		return err
	}
	defer conn.Close()

	s.addWebSocketConn(conn)
	defer s.removeWebSocketConn(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			s.log.Error(err)
			break
		}
	}
	return nil
}

func (s *OMTServer) addWebSocketConn(conn *websocket.Conn) {
	s.wsLock.Lock()
	defer s.wsLock.Unlock()

	if s.wsConns == nil {
		s.wsConns = []*websocket.Conn{}
	}

	s.wsConns = append(s.wsConns, conn)
}

func (s *OMTServer) removeWebSocketConn(conn *websocket.Conn) {
	s.wsLock.Lock()
	defer s.wsLock.Unlock()

	conns := []*websocket.Conn{}
	if s.wsConns != nil {
		for _, v := range s.wsConns {
			if v == conn {
				continue
			}
			tmp := v
			conns = append(conns, tmp)
		}
	}

	s.wsConns = conns
}
