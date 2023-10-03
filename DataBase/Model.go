package DataBase

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string
	//PhoneNumber string
	Login    bool
	Username string
	Password string
}
type RestaurantAdmin struct {
	gorm.Model
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Address     string
	Username    string
	Password    string
	Login       bool
	Restaurant  Restaurant
}

type User struct {
	gorm.Model
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Address     string
	Login       bool
	Username    string
	Password    string
	Currency    int64
	Orders      []Order
}
type Restaurant struct {
	gorm.Model
	Name              string
	Score             string
	Address           string
	Verified          bool
	Delivery          bool
	Category          string
	PosterLink        string
	RestaurantAdminID uint
	Menu              []Food
}
type Food struct {
	gorm.Model
	Name         string
	PosterLink   string
	Description  string
	Price        string
	RestaurantID uint
}
type Order struct {
	gorm.Model
	ID                int64
	Status            string
	Price             int64
	Name              string
	Restaurant        string
	PaymentMethod     string
	User              string
	UserAddress       string
	RestaurantAddress string
	UserID            uint
}
