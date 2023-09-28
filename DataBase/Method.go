package DataBase

import "golang.org/x/crypto/bcrypt"

func (gdb *DB) GetUserByUsername(username string) (*User, error) {
	u := User{}
	err := gdb.sql.Where(User{Username: username}).First(u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (gdb *DB) CheckUserDuplicate(username string, email string, phoneNumber string) (string, error) {
	UsernameDuplicate, err := gdb.CheckUserDuplicateByUsername(username)
	if err != nil {
		return "", err
	}
	if UsernameDuplicate {
		return "username", nil
	}
	EmailDuplicate, err := gdb.CheckUserDuplicateByEmail(email)
	if err != nil {
		return "", err
	}
	if EmailDuplicate {
		return "email", nil
	}
	PhoneNumberDuplicate, err := gdb.CheckUserDuplicateByPhoneNumber(phoneNumber)
	if err != nil {
		return "", err
	}
	if PhoneNumberDuplicate {
		return "phoneNumber", nil
	}
	return "user not exist", nil
}

func (gdb *DB) CheckUserDuplicateByUsername(username string) (bool, error) {
	_, err := gdb.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (gdb *DB) CheckUserDuplicateByEmail(email string) (bool, error) {
	u := User{}
	err := gdb.sql.Where(User{Email: email}).First(u).Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (gdb *DB) CheckUserDuplicateByPhoneNumber(phoneNumber string) (bool, error) {
	u := User{}
	err := gdb.sql.Where(User{PhoneNumber: phoneNumber}).First(u).Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (gdb *DB) CreateNewUser(u *User) error {
	if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0); err != nil {
		return err
	} else {
		u.Password = string(pw)
	}

	err := gdb.sql.Create(u).Error
	return err
}
