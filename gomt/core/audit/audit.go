package audit

import (
	"encoding/json"
	"fmt"
	"gomt/core/logger/lumberjack"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AuditLogger struct {
	log      *lumberjack.Logger
	hostname string
	tag      string
}

func NewAuditLogger(tag string, filename string, maxSize, maxBackups int) *AuditLogger {
	return &AuditLogger{
		tag: tag,
		log: &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize, // kilobyte
			MaxBackups: maxBackups,
		},
	}
}

func (s *AuditLogger) WriteApiLog(t time.Time, severity string, event string, user string, remote string, status int, reqBody any) {
	text := fmt.Sprintf(`operation="%v" source="%v" user="%v" status=%v`,
		event,
		remote, user,
		status,
	)
	if reqBody != nil {
		switch body := reqBody.(type) {
		case string:
			text += fmt.Sprintf(` content="%v"`, strings.ReplaceAll(body, `"`, `\"`))
		default:
			content, _ := json.Marshal(body)
			text += fmt.Sprintf(` content="%v"`, strings.ReplaceAll(string(content), `"`, `\"`))
		}
	}
	s.WriteLog(t, severity, text)
}

func (s *AuditLogger) WriteLog(t time.Time, severity string, message string) {
	if t.IsZero() {
		t = time.Now()
	}
	nl := ""
	if !strings.HasSuffix(message, "\n") {
		nl = "\n"
	}
	text := fmt.Sprintf("%v %v: %v%v",
		t.Format("2006-01-02T15:04:05"), s.tag,
		message, nl,
	)
	if _, err := s.log.Write([]byte(text)); err != nil {
		logrus.Error(errors.Wrap(err, "write log"))
	}
	s.log.Close()
}
