package DataBase

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/Config"
)

type DB struct {
	cfg   Config.Config
	sql   *gorm.DB
	redis *redis.Client
}

func CreateAndConnectToDb(cfg Config.Config) (*DB, error) {
	c := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Sql.Host,
		cfg.Database.Sql.Port,
		cfg.Database.Sql.Username,
		cfg.Database.Sql.Name,
		cfg.Database.Sql.Password,
	)

	// Create a new connection
	db, err := gorm.Open(postgres.Open(c), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Database.Redis.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &DB{
		cfg:   cfg,
		sql:   db,
		redis: client,
	}, nil
}

func (gdb *DB) CreateModel() error {
	err := gdb.sql.AutoMigrate(Order{})
	if err != nil {
		return err
	}
	err = gdb.sql.AutoMigrate(User{})
	if err != nil {
		return err
	}
	err = gdb.sql.AutoMigrate(RestaurantAdmin{})
	if err != nil {
		return err
	}

	err = gdb.sql.AutoMigrate(Restaurant{})
	if err != nil {
		return err
	}
	err = gdb.sql.AutoMigrate(Admin{})
	if err != nil {
		return err
	}
	err = gdb.sql.AutoMigrate(Food{})
	if err != nil {
		return err
	}

	return nil

}
