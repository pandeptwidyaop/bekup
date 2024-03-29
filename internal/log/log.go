package log

import "github.com/sirupsen/logrus"

var instance *logrus.Logger

func GetInstance() *logrus.Logger {
	if instance == nil {
		instance = logrus.New()
	}

	return instance
}
