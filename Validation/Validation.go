package Validation

import (
	"errors"
	"strings"
)

const (
	EnglishChar        = "qwertyuiopasdfghjklzxcvbnm"
	EnglishCharCapital = "QWERTYUIOPASDFGHJKLZXCVBNM"
	Numbers            = "0123456789"
	SpecialCharValid   = "!@#$%^&*_.-+=?"
	SpecialCharInValid = "`~':;][{}(),<> \\ | / \"  "
)

type Validation struct {
	ValidDomain       []string
	IllegalWords      []string
	ValidPhoneCountry []string
}

// return message for ValidateData func

type ValidationMessage struct {
	Message     string
	FirstName   bool
	LastName    bool
	Email       bool
	Address     bool
	PhoneNumber bool
	Username    bool
	Password    bool
}

type ValidateData struct {
	FirstName   string
	LastName    string
	Email       string
	Address     string
	PhoneNumber string
	Username    string
	Password    string
}

// Create a new instance of Validation

func CreateValidation(validDomain []string, IllegalWords []string, validPhoneCountry []string) (*Validation, error) {
	if len(IllegalWords) == 0 || len(validDomain) == 0 || len(validPhoneCountry) == 0 {
		return nil, errors.New("config is invalid")
	}
	return &Validation{
		IllegalWords:      IllegalWords,
		ValidDomain:       validDomain,
		ValidPhoneCountry: validPhoneCountry,
	}, nil
}

func (v *Validation) ValidateData(inputData ValidateData) (ValidationMessage, error) {
	var message string
	fn := v.validateFirstName(inputData.FirstName)
	ln := v.validateLastName(inputData.LastName)
	ad := v.validateAddress(inputData.Address)
	pn := v.ValidatePhoneNumber(inputData.PhoneNumber)
	em := v.ValidateEmail(inputData.Email)
	un := v.validateUserName(inputData.Username)
	ps := v.validatePassword(inputData.Password)
	var err error
	if result := fn && ln && ad && pn && em && un && ps; result == true {
		message = "Successful"
		err = nil
	} else {
		message = "invalid data"
		err = errors.New("invalid data")
	}
	return ValidationMessage{
		Message:     message,
		FirstName:   fn,
		LastName:    ln,
		Email:       em,
		Address:     ad,
		PhoneNumber: pn,
		Username:    un,
		Password:    ps,
	}, err

}

func (v *Validation) ValidateEmail(email string) bool {
	if email == "" || !(strings.Contains(email, "@")) {
		return false
	}
	emailPart := strings.Split(email, "@")
	if len(emailPart) != 2 {
		return false
	}
	if !v.checkValidDomain(emailPart[1]) {
		return false
	}
	for _, i2 := range emailPart[0] {
		if strings.Contains(SpecialCharInValid, string(i2)) {
			return false
		}

	}
	return true
}

func (v *Validation) validateFirstName(firstname string) bool {
	if firstname == "" {
		return false
	}
	for _, i2 := range firstname {
		if !strings.Contains(EnglishChar, strings.ToLower(string(i2))) {
			return false
		}
	}

	return true
}

func (v *Validation) validateLastName(lastname string) bool {
	if lastname == "" {
		return false
	}
	// check for invalid char
	for _, i2 := range lastname {
		if !strings.Contains(EnglishChar, strings.ToLower(string(i2))) {
			return false
		}
	}

	return true
}

func (v *Validation) validatePassword(password string) bool {
	if password == "" || len(password) < 8 {
		return false
	}
	EnChar := false
	EnCharCapital := false
	numbers := false
	SpCharValid := false

	for _, i2 := range password {
		// check for invalid char
		if strings.Contains(SpecialCharInValid, string(i2)) {
			return false
		}
		// check for EnglishChar
		if EnChar == false && strings.Contains(EnglishChar, string(i2)) {
			EnChar = true
			continue
		}
		// check for EnglishCharCapital
		if EnCharCapital == false && strings.Contains(EnglishCharCapital, string(i2)) {
			EnCharCapital = true
			continue
		}
		// check for Numbers
		if numbers == false && strings.Contains(Numbers, string(i2)) {
			numbers = true
			continue
		}
		// check for SpecialChar Valid
		if SpCharValid == false && strings.Contains(SpecialCharValid, string(i2)) {
			SpCharValid = true
			continue
		}
	}

	return EnChar && EnCharCapital && numbers && SpCharValid
}

func (v *Validation) validateUserName(username string) bool {
	if username == "" {
		return false
	}
	for _, i2 := range username {
		if strings.Contains(SpecialCharInValid, string(i2)) {
			return false
		}

	}

	return true
}

func (v *Validation) validateAddress(address string) bool {
	if address == "" {
		return false
	}
	for _, i2 := range address {
		if strings.Contains(SpecialCharInValid, string(i2)) {
			return false
		}
		if strings.Contains(SpecialCharValid, string(i2)) {
			return false
		}
	}

	return true
}

func (v *Validation) ValidatePhoneNumber(phoneNumber string) bool {
	if phoneNumber == "" {
		return false
	}
	PhoneNumberPart := strings.Split(phoneNumber, "-")
	if len(PhoneNumberPart) != 2 {
		return false
	}

	if !v.checkValidPhoneCountry(PhoneNumberPart[0]) {
		return false
	}

	if len(PhoneNumberPart[1]) > 10 || len(PhoneNumberPart[1]) < 6 {
		return false
	}
	for _, i2 := range PhoneNumberPart[1] {
		if !strings.Contains(Numbers, string(i2)) {
			return false
		}
	}
	return true
}

func (v *Validation) checkValidDomain(domain string) bool {
	for _, d := range v.ValidDomain {
		if d == domain {
			return true
		}
	}
	return false
}

func (v *Validation) checkValidPhoneCountry(input string) bool {
	for _, s := range v.ValidPhoneCountry {
		if s == input {
			return true
		}
	}
	return false
}
