package Handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/Validation"
	"net/http"
)

type SignUprRequestBody struct {
	FirstName   string `json:"FirstName" binding:"required"`
	LastName    string `json:"LastName" binding:"required"`
	Email       string `json:"Email" binding:"required"`
	PhoneNumber string `json:"PhoneNumber" binding:"required"`
	Address     string `json:"Address" binding:"required"`
	Username    string `json:"Username" binding:"required"`
	Password    string `json:"Password" binding:"required"`
}

func (s *Server) UserSignUpHandler(g *gin.Context) {
	var reqData SignUprRequestBody

	// unmarshal json
	err := g.BindJSON(&reqData)
	if err != nil {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	//Validate user data
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

	// check user duplicate
	UserDuplicate, err := s.Db.CheckUserDuplicate(reqData.Username, reqData.Email, reqData.PhoneNumber)
	if err != nil {
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
		s.Logger.WithError(err).Error("Most likely, there is a problem with the database")
		return
	}
	if UserDuplicate != "user not exist" {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user with this %v exist", UserDuplicate)})
		return
	}
	if err := s.Db.CreateTempUserOnRedis("saleh"); err != nil {
		g.IndentedJSON(http.StatusForbidden, "New User Created")
		return
	}
	//// Create new user in database
	//user := &DataBase.User{
	//	FirstName:   reqData.FirstName,
	//	LastName:    reqData.LastName,
	//	Email:       reqData.Email,
	//	PhoneNumber: reqData.PhoneNumber,
	//	Address:     reqData.Address,
	//	Login:       false,
	//	Username:    reqData.Username,
	//	Password:    reqData.Password,
	//	Currency:    0,
	//}
	//
	//if err := s.Db.CreateNewUser(user); err != nil {
	//	s.Logger.WithError(err).Error("can not create a new user")
	//	g.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong ,please try later!"})
	//	return
	//}
	g.IndentedJSON(http.StatusCreated, "New User Created")
	return
}
