package Handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/Authentication"
	"main/DataBase"
	"main/Validation"
	"net/http"
)

type RegisterRequestBody struct {
	FirstName   string `json:"FirstName" binding:"required"`
	LastName    string `json:"LastName" binding:"required"`
	Email       string `json:"Email" binding:"required"`
	PhoneNumber string `json:"PhoneNumber" binding:"required"`
	Address     string `json:"Address" binding:"required"`
	Username    string `json:"Username" binding:"required"`
	Password    string `json:"Password" binding:"required"`
}

type VerifyRequestBody struct {
	PhoneNumber string `json:"PhoneNumber" binding:"required"`
}
type VerifyBody struct {
	PhoneNumber string `json:"PhoneNumber" binding:"required"`
	Code        string `json:"Code" binding:"required"`
}

type LoginRequestBody struct {
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

func (s *Server) UserVerifyRequestHandler(g *gin.Context) {
	var reqData VerifyRequestBody

	// unmarshal json
	err := g.BindJSON(&reqData)
	if err != nil {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	////Validate user data
	if result := s.Vln.ValidatePhoneNumber(reqData.PhoneNumber); !result {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "PhoneNumber is invalid"})
		return
	}

	// check user duplicate
	UserDuplicate, err := s.Db.CheckUserDuplicateByPhoneNumber(reqData.PhoneNumber)
	if err != nil {
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		s.Logger.WithError(err).Error("Most likely, there is a problem with the database")
		return
	}
	if UserDuplicate {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user with this PhoneNumber exist"})
		return
	}

	//check user duplicate on redis
	_, err = s.Db.GetTempUserCode(reqData.PhoneNumber)
	if err != nil {
		if err.Error() != "redis: nil" {
			s.Logger.WithError(err).Error("redis database have problem for get data of temp user")
			g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
			return
		}

	} else {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The user has registered with this email and we are waiting for her validation confirmation !"})
		return
	}
	Code := s.Lo.GenerateOtpCode()
	if err := s.Db.CreateTempUserOnRedis(reqData.PhoneNumber, Code); err != nil {
		s.Logger.WithError(err).Error("redis database have problem for create temp user")
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		return
	}

	err = s.Lo.SendSMSToUser(Code, reqData.PhoneNumber)
	if err != nil {
		s.Logger.WithError(err).Error("can not send SMS ")
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		return
	}
	g.IndentedJSON(http.StatusOK, "New Request for Verify Created")
	return
}

func (s *Server) UserVerifyHandler(g *gin.Context) {
	var reqData VerifyBody

	// unmarshal json
	err := g.BindJSON(&reqData)
	if err != nil {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	//check user verify or not on redis
	val, err := s.Db.GetTempUserCode(reqData.PhoneNumber)
	if err != nil {
		if err.Error() == "redis: nil" {
			g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user must verify"})
			return
		}
		s.Logger.WithError(err).Error("redis database have problem for get data of temp user")
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		return
	}

	if val == "Verify-successfully" {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user already been verifyed"})
		return
	}

	if val != reqData.Code {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "code is wrong"})
		return
	}
	err = s.Db.CreateTempUserOnRedis(reqData.PhoneNumber, "Verify-successfully")
	if err != nil {
		s.Logger.WithError(err).Error("redis database have problem for set data of temp user")
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		return
	}
	g.IndentedJSON(http.StatusOK, gin.H{"message": "successful"})
	return
}

func (s *Server) UserRegisterHandler(g *gin.Context) {

	var reqData RegisterRequestBody

	// unmarshal json
	err := g.BindJSON(&reqData)
	if err != nil {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	//// Validate user data
	if message, err := s.Vln.ValidateData(Validation.ValidateData{
		FirstName:   reqData.FirstName,
		LastName:    reqData.LastName,
		Email:       reqData.Email,
		Address:     reqData.Address,
		PhoneNumber: reqData.PhoneNumber,
		Username:    reqData.Username,
		Password:    reqData.Password,
	}); err != nil {
		g.IndentedJSON(http.StatusBadRequest, message)
		return
	}

	//check user verify or not on redis
	val, err := s.Db.GetTempUserCode(reqData.PhoneNumber)
	if err != nil {
		if err.Error() == "redis: nil" {
			g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user must signup"})
			return
		}
		s.Logger.WithError(err).Error("redis database have problem for get data of temp user")
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		return
	}
	if val != "Verify-successfully" {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user must be verify"})
		return
	}

	// check user duplicate on sql
	UserDuplicate, err := s.Db.CheckUserDuplicate(reqData.Username, reqData.PhoneNumber)
	if err != nil {
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		s.Logger.WithError(err).Error("Most likely, there is a problem with the database")
		return
	}
	if UserDuplicate != "user not exist" {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user with this %v exist", UserDuplicate)})
		return
	}

	//// Create new user in database
	user := &DataBase.User{
		FirstName:   reqData.FirstName,
		LastName:    reqData.LastName,
		Email:       reqData.Email,
		PhoneNumber: reqData.PhoneNumber,
		Address:     reqData.Address,
		Login:       false,
		Username:    reqData.Username,
		Password:    reqData.Password,
		Currency:    0,
	}

	if err := s.Db.CreateNewUser(user); err != nil {
		s.Logger.WithError(err).Error("can not create a new user")
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		return
	}
	g.IndentedJSON(http.StatusCreated, "New User Created")
	return

}

func (s *Server) UserLoginHandler(g *gin.Context) {
	var reqData LoginRequestBody

	// unmarshal json
	err := g.BindJSON(&reqData)
	if err != nil {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	err = s.At.AuthenticateUserWithCredentials(Authentication.Credentials{
		Username: reqData.Username,
		Password: reqData.Password,
	})
	if err != nil {
		if err.Error() == "user not exist" {
			g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user not found pelese signup first"})
			return
		}
		if err.Error() == "the password is not correct" {
			g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user password is not correct"})
			return
		}
		s.Logger.WithError(err).Error("can not find user in sql database")
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		return
	}
	Token, err := s.At.GenerateJwtToken(reqData.Username, "User")
	if err != nil {
		s.Logger.WithError(err).Error("can not Generate JwtToken ")
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		return
	}
	g.IndentedJSON(http.StatusOK, gin.H{"Token": *Token})
	return
}
