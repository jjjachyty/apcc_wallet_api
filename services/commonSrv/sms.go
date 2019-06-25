package commonSrv

import (
	"errors"
	"time"

	"github.com/mojocn/base64Captcha/store"
)

type SMSService struct{}

var globalStore = store.NewMemoryStore(10240, time.Second*60)

func (SMSService) SendSMS(phone string) error {
	globalStore.Set(phone, "1234")
	return nil
}

func (SMSService) VerificationSMS(phone, sms string) error {
	storeSms := globalStore.Get(phone, false)
	if storeSms == sms {
		globalStore.Get(phone, true)
		return nil
	}
	return errors.New("短信校验失败")
}
