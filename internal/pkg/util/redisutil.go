package util

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func TokenBlacklistKey(token string) string {
	return "auth:blacklist:" + token
}

// GetCache mengambil JSON dari Redis lalu otomatis mengubahnya ke struct tujuan (dest).
// Return true kalau data ketemu (Cache Hit), false kalau kosong/error (Cache Miss).
func GetCache(ctx context.Context, rdb *redis.Client, key string, dest any) bool {
	if rdb == nil {
		return false
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false
	}

	return true
}

// SetCache mengubah struct apapun menjadi JSON lalu menyimpannya ke Redis dengan batas umur (TTL).
func SetCache(ctx context.Context, rdb *redis.Client, key string, value any, ttl time.Duration) error {
	if rdb == nil {
		return nil
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, bytes, ttl).Err()
}

// DeleteCache menghapus satu atau beberapa key sekaligus dari Redis.
func DeleteCache(ctx context.Context, rdb *redis.Client, keys ...string) error {
	if rdb == nil || len(keys) == 0 {
		return nil
	}

	return rdb.Del(ctx, keys...).Err()
}

// ExistsCache mengecek apakah sebuah key ada di dalam RAM Redis.
func ExistsCache(ctx context.Context, rdb *redis.Client, key string) bool {
	if rdb == nil || key == "" {
		return false
	}

	count, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	return count > 0
}

// IncrementCounter menambah angka penghitung +1 di Redis secara atomic dan memasang TTL saat pertama kali dibuat.
func IncrementCounter(ctx context.Context, rdb *redis.Client, key string, ttl time.Duration) (int64, error) {
	if rdb == nil || key == "" {
		return 0, nil
	}

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 && ttl > 0 {
		_ = rdb.Expire(ctx, key, ttl).Err()
	}

	return count, nil
}
