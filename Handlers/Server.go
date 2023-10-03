package Handlers

import (
	"github.com/sirupsen/logrus"
	"main/Authentication"
	"main/DataBase"
	"main/Logic"
	"main/Validation"
)

type Server struct {
	Logger *logrus.Logger
	Db     *DataBase.DB
	Vln    *Validation.Validation
	At     *Authentication.Authentication
	Lo     *Logic.Logic
}
