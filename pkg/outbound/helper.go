package outbound

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/klaxn/klaxn-api/internal/config"
)

func CreateOutbounds(conf []*config.OutboundConfig, logger *logrus.Logger) []Sender {
	var out []Sender

	for _, outboundConfig := range conf {
		l := logger.WithField("outbound", outboundConfig.Name)

		if outboundConfig.Enabled {
			outbound, err := createOutbound(outboundConfig.Name, outboundConfig.Config, logger)
			if err != nil {
				l.Error(err)
			} else {
				out = append(out, outbound)
				l.Info("finished configuration")
			}
		} else {
			l.Info("disabled outbound")
		}
	}

	return out
}

func createOutbound(name string, config map[string]interface{}, logger *logrus.Logger) (Sender, error) {
	if name == "twilio" {
		return NewTwilio(logger, config)
	}

	if name == "alert_logger" {
		return NewAlertLogger(logger), nil
	}

	return nil, fmt.Errorf("unknown outbound type: %s", name)
}
