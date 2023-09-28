package Config

type Config struct {
	Database struct {
		Sql struct {
			Host     string `json:"HOST" env:"HOST"  env-default:"localhost"`
			Port     int    `json:"PORT" env:"PORT" env-default:"5432"`
			Name     string `json:"NAME" env:"NAME" env-default:"food_order"`
			Username string `json:"SQL_USERNAME" env:"SQL_USERNAME" env-default:"postgres"`
			Password string `json:"PASSWORD" env:"PASSWORD" env-default:"saleh"`
		}
		Redis struct {
			Addr string `json:"Addr" env:"Addr" env-default:"saleh"`
		}
	}
	ValidDomain       []string `json:"ValidDomain" env:"ValidDomain"`
	IllegalWords      []string `json:"IllegalWords" env:"IllegalWords"`
	ValidPhoneCountry []string `json:"ValidPhoneCountry" env:"ValidPhoneCountry"`
}
