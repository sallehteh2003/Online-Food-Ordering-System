package Logic

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"math/rand"
	"strings"
	"time"
)

const (
	EnglishChar        = "qwertyuiopasdfghjklzxcvbnm"
	EnglishCharCapital = "QWERTYUIOPASDFGHJKLZXCVBNM"
	Numbers            = "0123456789"
	SpecialCharValid   = "!@#$%^&*_.-+="
)

type Logic struct {
	EmailAddr  string
	AccountSid string
	AuthToken  string
}

func (l *Logic) GenerateOtpCode() string {
	code := ""
	for i := 0; i < 8; i++ {
		h, m, s := time.Now().Clock()
		ran := rand.Int63n(int64(h*3600 + m*60 + s))
		switch ran % 4 {
		case 0:
			code += string(EnglishChar[ran%26])
			continue
		case 1:
			code += string(EnglishCharCapital[ran%26])

			continue
		case 2:
			code += string(Numbers[ran%10])

			continue
		case 3:
			code += string(SpecialCharValid[ran%12])
			continue

		}
	}

	return code
}

func (l *Logic) SendSMSToUser(code string, to string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: l.AccountSid,
		Password: l.AuthToken,
	})
	to = strings.ReplaceAll(to, "-", "")
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom("+12182202927")
	params.SetBody(fmt.Sprintf("your otp code from online shop food order is %v plese verify youe acconut.", code))

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending SMS message: %v", err.Error()))

	}
	return nil
}
