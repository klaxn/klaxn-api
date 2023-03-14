package outbound

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/klaxn/klaxn-api/internal/config"
)

type Twilio struct {
	client      *twilio.RestClient
	logger      *logrus.Entry
	fromAddress string
}

func (t *Twilio) SendMessage(to, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(t.fromAddress)
	params.SetBody(message)

	resp, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return err
	} else {
		response, _ := json.Marshal(*resp)
		t.logger.Infof("Response from twilio %s", string(response))
	}

	return nil
}

func NewTwilio(logger logrus.FieldLogger, configMap map[string]interface{}) (*Twilio, error) {
	accountSid, err := config.GetStringFromConfigMap(configMap, "accountSid")
	if err != nil {
		return nil, err
	}

	authToken, err := config.GetStringFromConfigMap(configMap, "authToken")
	if err != nil {
		return nil, err
	}

	fromAddress, err := config.GetStringFromConfigMap(configMap, "fromAddress")
	if err != nil {
		return nil, err
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &Twilio{
		client:      client,
		logger:      logger.WithField("outbound", "twilio"),
		fromAddress: fromAddress,
	}, nil
}
