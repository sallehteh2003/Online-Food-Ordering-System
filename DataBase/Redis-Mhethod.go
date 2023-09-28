package DataBase

import (
	"context"
	"fmt"
	"time"
)

func (gdb *DB) CreateTempUserOnRedis(email string) error {
	ctx := context.Background()

	err := gdb.redis.SetEx(ctx, email, "saleh", time.Duration(int64(time.Minute)*10)).Err()
	if err != nil {
		return err
	}
	val, err := gdb.redis.Get(ctx, "email").Result()
	if err != nil {
		return err
	}
	fmt.Println("foo", val)
	return nil
}
