package sms_provider

import (
	"fmt"
	"strconv"

	"github.com/kavenegar/kavenegar-go"
	"github.com/supabase/gotrue/internal/conf"
)

type KavenegarProvider struct {
	Config *conf.KavenegarProviderConfiguration
	API    *kavenegar.Kavenegar
	OTP    *kavenegar.VerifyService
}

type kavenegarError struct {
	Code      int    `json:"code"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	ErrorType string `json:"error_type"`
}

func NewKavenegarProvider(config conf.KavenegarProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &KavenegarProvider{
		Config: &config,
		API:    kavenegar.New(config.ApiKey),
		OTP:    kavenegar.NewVerifyService(kavenegar.NewClient(config.ApiKey)),
	}, nil
}

func (t *KavenegarProvider) SendMessage(phone string, message string, channel string) (string, error) {
	switch channel {
	case SMSProvider:
		return t.SendSms(phone, message)
	default:
		return "", fmt.Errorf("channel type %q is not supported for Kavenegar")
	}
}

func (t *KavenegarProvider) SendSms(phone string, message string) (string, error) {
	//sender := ""
	receptor := []string{phone[2:]}
	
	template := t.Config.OTPTemplate
	
	res, err := t.OTP.Lookup(receptor[0], template, message, &kavenegar.VerifyLookupParam{
		Type: kavenegar.Type_VerifyLookup_Sms,
		Tokens: map[string]string{"%token": message},
	})
	if err != nil {
		return "", err
	}

	//res, err := t.API.Message.Send(sender, receptor, message, nil)
	//if err != nil {
		//return "", err
	//}

	status, err := t.API.Message.Status([]string{strconv.Itoa(res.MessageID)})
	if err != nil {
		return "", err
	}

	return strconv.Itoa(status[0].MessageId), nil
}
