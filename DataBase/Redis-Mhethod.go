package DataBase

import (
	"context"
	"time"
)

func (gdb *DB) CreateTempUserOnRedis(email string, code string) error {
	ctx := context.Background()
	err := gdb.redis.SetEx(ctx, email, code, time.Duration(int64(time.Minute)*10)).Err()
	if err != nil {
		return err
	}

	return nil
}
func (gdb *DB) GetTempUserCode(email string) (string, error) {
	ctx := context.Background()
	val, err := gdb.redis.Get(ctx, email).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
func (gdb *DB) DeleteTempUserOnRedis(email string) error {
	ctx := context.Background()
	_, err := gdb.redis.Del(ctx, email).Result()
	if err != nil {
		return err
	}
	return nil
}
