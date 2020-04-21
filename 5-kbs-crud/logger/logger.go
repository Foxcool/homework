package logger

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type (
	coreStatus   string
	coreResource string

	userStatus    string
	userEventType string
)

const (
	CORECONNECTED coreStatus = "CONNECTED"
	COREFAILED    coreStatus = "FAILED"
	CORESTARTED   coreStatus = "STARTED"

	COREDB coreResource = "DB"

	USERSUCCESS userStatus = "SUCCESS"
	USERFAIL    userStatus = "FAIL"

	USERCREATE userEventType = "USER_CREATE"
	USERUPDATE userEventType = "USER_UPDATE"
	USERDELETE userEventType = "USER_DELETE"
)

var logger Logger

// Logger Init

type Logger interface {
	Core() CoreEntrier
	Auth() AuthEntrier
}

type homeworkLog struct {
	lg *logrus.Logger
}

func NewLogger(out io.Writer, levStr string, format string) Logger {
	lg := logrus.New()
	level, err := logrus.ParseLevel(levStr)
	if err != nil {
		lg.WithError(err).Warnln("Can't set logging level")
	}

	lg.SetLevel(level)
	lg.SetOutput(out)

	if format == "JSON" {
		lg.SetFormatter(&logrus.JSONFormatter{})
	}

	logger = &homeworkLog{lg: lg}
	return logger
}

// Core Entrier

type CoreEntrier interface {
	Info(resource coreResource, host, port, version string, status coreStatus)
	Fatal(resource coreResource, host, port string, status coreStatus)
	Debug(...interface{})
}

func (l *homeworkLog) Core() CoreEntrier {
	entry := l.lg.WithField("context", "CORE")

	return &coreEntry{
		entry,
	}
}

type coreEntry struct {
	entry *logrus.Entry
}

func (e *coreEntry) Info(resource coreResource, host, port, version string, status coreStatus) {
	if resource != "" {
		resource = "resource=" + resource
	}

	addr := host + ":" + port
	if addr != ":" {
		addr = " addr=" + strings.ToUpper(addr)
	} else {
		addr = ""
	}

	if version != "" {
		version = " version=" + strings.ToUpper(version)
	}

	e.entry.Infof("%s%s%s status=%s", resource, addr, version, status)
}

func (e *coreEntry) Fatal(resource coreResource, host, port string, status coreStatus) {
	addr := host + ":" + port

	e.entry.Fatalf("resource=%s addr=%s status=%s", resource, addr, status)
}

func (e *coreEntry) Debug(val ...interface{}) {
	e.entry.Debugln(val...)
}

// Auth Entrier

type AuthEntrier interface {
	Info(eventType userEventType, entityLogin string, status userStatus)
}

func (l *homeworkLog) Auth() AuthEntrier {
	entry := l.lg.WithField("context", "USER")

	return &authEntry{
		entry,
	}
}

type authEntry struct {
	entry *logrus.Entry
}

func (e *authEntry) Info(eventType userEventType, entityLogin string, status userStatus) {
	e.entry.Infof("eventType=%s login=%s status=%s", eventType, entityLogin, status)
}
