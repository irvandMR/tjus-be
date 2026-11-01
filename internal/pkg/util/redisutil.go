package util

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

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
