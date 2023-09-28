package Authentication

import (
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"main/DataBase"
	"time"
)

type Authentication struct {
	db *DataBase.DB
	// jwtSecretKey is the JWT secret key. Each time the server starts, new key is generated.
	jwtUserSecretKey      []byte
	jwtAdminSecretKey     []byte
	jwtResAdminSecretKey  []byte
	jwtExpirationDuration time.Duration
	logger                *logrus.Logger
}
type claims struct {
	jwt.MapClaims
	Username string `json:"username"`
}
type Credentials struct {
	Username string
	Password string
}

func CreateAuthentication(db *DataBase.DB, jwtExpirationInMinutes int64, logger *logrus.Logger) (*Authentication, error) {
	UserSecretKey, err := generateRandomKey()
	if err != nil {
		return nil, err
	}
	AdminSecretKey, err := generateRandomKey()
	if err != nil {
		return nil, err
	}
	ResAdminSecretKey, err := generateRandomKey()
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("the database is essential for authentication")
	}

	return &Authentication{
		db:                    db,
		jwtUserSecretKey:      UserSecretKey,
		jwtAdminSecretKey:     AdminSecretKey,
		jwtResAdminSecretKey:  ResAdminSecretKey,
		jwtExpirationDuration: time.Duration(int64(time.Minute) * jwtExpirationInMinutes),
		logger:                logger,
	}, nil
}

func (a *Authentication) AuthenticateUserWithCredentials(cred Credentials) error {

	// check user exist
	user, err := a.db.GetUserByUsername(cred.Username)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New("user not exist")
		}
		return err
	}

	//check user password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password)); err != nil {
		return errors.New("the password is not correct")
	}
	return nil
}

func (a *Authentication) GenerateJwtToken(username string, Role string) (token *string, err error) {
	expirationTime := time.Now().Add(a.jwtExpirationDuration)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Username: username,
		MapClaims: jwt.MapClaims{
			"expired_at": expirationTime.Unix(),
		},
	})

	// Calculate the signed account string format of JWT key
	tokenString := ""
	switch Role {
	case "Admin":
		tokenString, err = tokenJWT.SignedString(a.jwtAdminSecretKey)
		break
	case "User":
		tokenString, err = tokenJWT.SignedString(a.jwtUserSecretKey)
		break
	case "ResAdmin":
		tokenString, err = tokenJWT.SignedString(a.jwtResAdminSecretKey)
		break
	default:
		return nil, errors.New("invalid Role for GenerateJwtToken")
	}

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (a *Authentication) checkToken(token string, Role string) (*claims, error) {
	c := &claims{}
	tkn, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		switch Role {
		case "Admin":
			return a.jwtAdminSecretKey, nil

		case "User":
			return a.jwtUserSecretKey, nil

		case "ResAdmin":
			return a.jwtResAdminSecretKey, nil
		default:
			return nil, errors.New("invalid Role")
		}

	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token")
		}

		a.logger.WithError(err).Warn("can not validate the token of the user")
		return nil, errors.New("bad error in validating the token")
	}

	if !tkn.Valid {
		return nil, errors.New("unauthorized")
	}

	return c, nil
}

// generateRandomKey
// Each time that Auth is initialized, generateRandomKey is called to
// generate another key
func generateRandomKey() ([]byte, error) {
	jwtKey := make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		return nil, err
	}

	return jwtKey, nil
}
