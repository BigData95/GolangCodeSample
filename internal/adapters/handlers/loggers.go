package handlers

import "github.com/sirupsen/logrus"

func NewLogger(processId string) *logrus.Entry {
	return logrus.WithFields(
		logrus.Fields{
			"processId": processId,
		})
}
