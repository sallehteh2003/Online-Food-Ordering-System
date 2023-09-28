package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"main/Authentication"
	"main/Config"
	"main/DataBase"
	"main/Handlers"
	"main/Validation"
)

func main() {
	var cfg Config.Config
	logger := logrus.New()
	r := gin.Default()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetReportCaller(true)
	err := cleanenv.ReadConfig("Config/config.json", &cfg)
	if err != nil {
		logger.WithError(err).Panicln("failed to load the configs")
	} else {
		logger.Infof("successful to read the configs: %+v", cfg)
	}

	db, err := DataBase.CreateAndConnectToDb(cfg)
	if err != nil {
		logger.WithError(err).Fatalln("can not connect to database")
	}

	err = db.CreateModel()
	if err != nil {
		logger.WithError(err).Fatalln("can not create schema in database")
	}
	vln, err := Validation.CreateValidation(cfg.ValidDomain, cfg.IllegalWords, cfg.ValidPhoneCountry)
	if err != nil {
		logger.WithError(err).Fatalln("can not create instance of validation")
	}
	Auth, err := Authentication.CreateAuthentication(db, 15, logger)
	if err != nil {
		logger.WithError(err).Fatalln("can not create instance of Authentication")
	}
	server := Handlers.Server{
		Logger: logger,
		Db:     db,
		Vln:    vln,
		At:     Auth,
	}

	r.POST("/Api/User/Signup", server.UserSignUpHandler)
	if err := r.Run("localhost:8080"); err != nil {
		logrus.WithError(err).Fatalln("can not run server")
	}
}
