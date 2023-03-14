package outbound

import "github.com/sirupsen/logrus"

type AlertLogger struct {
	logger *logrus.Entry
}

func NewAlertLogger(logger logrus.FieldLogger) *AlertLogger {
	return &AlertLogger{logger: logger.WithField("outbound", "alert_logger")}

}

func (al *AlertLogger) SendMessage(to, message string) error {
	al.logger.WithField("to", to).WithField("message", message).Info("alert logger sending message to logger")
	return nil
}
