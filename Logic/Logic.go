package Logic

import "math/rand"

const (
	EnglishChar        = "qwertyuiopasdfghjklzxcvbnm"
	EnglishCharCapital = "QWERTYUIOPASDFGHJKLZXCVBNM"
	Numbers            = "0123456789"
	SpecialCharValid   = "!@#$%^&*_.-+="
)

type Logic struct {
	EmailAddr string
}

func (l *Logic) GenerateOtpLink() string {
	code := ""
	for i := 0; i < 8; i++ {
		ran := rand.Int31()
		switch ran % 4 {
		case 0:
			code += string(EnglishCharCapital[ran%26])
			break
		case 1:
			code += string(EnglishChar[ran%26])
			break
		case 2:
			code += string(Numbers[ran%10])
			break
		case 3:
			code += string(SpecialCharValid[ran%12])
			break
		}

	}
	return code
}
